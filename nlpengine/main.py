"""
Main function for training and predicting cloud resource state file assignments via Spacy.
"""
import os
import traceback

from ast import literal_eval
from copy import deepcopy
from random import randint, shuffle
from typing import List, Tuple, Union

import functions_framework
import requests
import spacy

import numpy as np
import pandas as pd

from sklearn.metrics import f1_score, precision_score, recall_score
from spacy.pipeline.textcat_multilabel import DEFAULT_MULTI_TEXTCAT_MODEL
from spacy.training import Example


@functions_framework.http
def train_and_predict(request):
    """
    Train and predict function. Input request is expected to have a json body with the following
    structure:
    {
        "workspace_docs": {...},
        "new_resource_docs": {...},
    }
    """
    try:
        # Loading data from request
        json_body = request.get_json()
        category_docs = literal_eval(json_body["workspace_docs"])
        new_resource_docs = literal_eval(json_body["new_resource_docs"])

        spacy.util.fix_random_seed(42)
        # Loading spacy model
        nlp = spacy.load("en_core_web_md")

        print("Loaded english model, beginning to preprocess training data.")
        all_training_examples = preprocess_training_data(nlp, category_docs=category_docs)

        print("Through preprocess_training_data().")

        train_eval_dict = _split_into_train_and_evaluation_data(
            training_data_examples=all_training_examples
        )

        training_data = train_eval_dict["train"]
        evaluation_data = train_eval_dict["evaluate"]

        print("Done preprocessing training data, beginning to initialize model.")
        textcat_multilabel_model = _initialize_classification_model(
            nlp=nlp, initial_examples=training_data[:3]
        )

        print("Done initializing model, beginning to train model.")
        final_loss_metrics = _train_model(
            nlp=nlp,
            textcat_multilabel_model=textcat_multilabel_model,
            training_data=training_data,
            evaluation_data=evaluation_data,
        )

        print("Done training model, beginning to make predictions.")
        resource_name_to_workspace_predictions = _predict(
            nlp=nlp,
            new_resource_docs=new_resource_docs,
            textcat_multilabel_model=textcat_multilabel_model,
        )

        print("Done making predictions, returning results.")

        return resource_name_to_workspace_predictions, 201
    except Exception as e:
        stack_trace = traceback.format_exc()
        print(f"{e}:\n{stack_trace}")
        return "Internal Server Error", 500


def preprocess_training_data(nlp: spacy.Language, category_docs: dict) -> List[Example]:
    """Creates all examples needed to train the text classification model."""
    list_of_examples = []

    # List of all categories
    categories = list(category_docs.keys())

    for category, doc in category_docs.items():
        current_category_examples = _create_examples_from_doc(
            nlp=nlp,
            current_cat=category,
            doc=doc,
            categories=categories,
        )
        list_of_examples.extend(current_category_examples)

    return list_of_examples


def _create_examples_from_doc(
    nlp: spacy.Language, current_cat: str, doc: str, categories: List[str]
) -> List[Example]:
    """
    Create a list of Spacy examples from the doc string.
    """
    doc_examples = []

    example_text_list = _doc_to_example_text_list(doc)

    gold_dict = _create_gold_dict(current_cat=current_cat, categories=categories)

    for example_text in example_text_list:
        doc = nlp(example_text)
        current_example = Example.from_dict(doc, gold_dict)
        doc_examples.append(current_example)

    if len(example_text_list) >= 5:
        for _ in range(len(example_text_list)):
            shuffle(example_text_list)
            end_index = randint(2, len(example_text_list))

            combined_example_text = _join_text_components(
                example_text_list=example_text_list, end_index=end_index
            )

            combined_doc = nlp(combined_example_text)
            current_example = Example.from_dict(combined_doc, gold_dict)
            doc_examples.append(current_example)

    return doc_examples


def _create_gold_dict(current_cat: str, categories: List[str]) -> dict:
    """
    Returns a "gold label" dictionary as required for training examples to be used
    within a Spacy Text classification model.

    For more details see: https://spacy.io/api/data-formats#dict-input
    """
    gold_dict = {
        "cats": {cat: 1.0 if cat == current_cat else 0.0 for cat in categories}
    }
    return gold_dict


def _join_text_components(example_text_list: List[str], end_index: int) -> str:
    """
    Joins the text components within example_text_list into a single string.
    """
    return ". ".join(example_text_list[:end_index]).replace("  ", " ") + "."


def _doc_to_example_text_list(doc: str) -> List[str]:
    """
    Converts a document into a series of individual pieces of text, separated by sentence.
    """
    return [text for text in doc.split(".") if text != " "]


def _split_into_train_and_evaluation_data(training_data_examples: List) -> dict:
    """
    Splits the training and evaluation data out of training_data_examples with an 80/20 train/eval split.
    """
    train_set_cutoff = int(len(training_data_examples) * 0.8)

    shuffle(training_data_examples)

    return {
        "train": training_data_examples[:train_set_cutoff],
        "evaluate": training_data_examples[train_set_cutoff:],
    }


def _initialize_classification_model(nlp: spacy.Language, initial_examples) -> Union:
    """Initialize the multi-label classification model."""
    config = {
        "threshold": 0.5,
        "model": DEFAULT_MULTI_TEXTCAT_MODEL,
    }

    textcat_multilabel = nlp.add_pipe("textcat_multilabel", config=config)

    textcat_multilabel.initialize(lambda: initial_examples, nlp=nlp)

    return textcat_multilabel


def _train_model(
    nlp: spacy.Language,
    textcat_multilabel_model,
    training_data: List[Example],
    evaluation_data: List[Example],
) -> dict:
    """
    Divides the training data into sets of 10 and updates the model using them.

    After each training batch, output the loss, as well as evaluate the accuracy on the evaluation
    data set set.

    Outputs a dict of the final loss and recall, precision and f1-score on the evaluation data set.
    """
    optimizer = nlp.create_optimizer()
    training_batches = np.array_split(training_data, 10)

    for i, batch in enumerate(training_batches):
        losses = textcat_multilabel_model.update(batch, sgd=optimizer)
        evaluation_data_scores = _performance_on_evaluation_data(
            textcat_multilabel_model=textcat_multilabel_model,
            evaluation_data=evaluation_data,
        )
        print(f"Batch {i} results:")
        print(f"Loss: {losses['textcat_multilabel']: .4f}")
        print(
            f"F1-Score: {evaluation_data_scores['f1']: .3f}, "
            f"Precision: {evaluation_data_scores['precision']: .3f}, "
            f"Recall: {evaluation_data_scores['recall']: .3f}\n"
        )

    return {"loss": losses, **evaluation_data_scores}


def _performance_on_evaluation_data(
    textcat_multilabel_model,
    evaluation_data: List[Example],
) -> dict:
    """
    Assessing performance on a held-out evaluation set of data. Works by making predictions
    with the current model on the evaluation data set and then calculating
    the f1, recall, and precision score.
    """
    actual_list_of_dicts = []

    prediction_list_of_dicts = []

    for eval_example in evaluation_data:
        eval_example_dict = deepcopy(eval_example.to_dict())

        actual_list_of_dicts.append(eval_example_dict["doc_annotation"]["cats"])

        scores = textcat_multilabel_model.predict([eval_example.reference])
        cats = textcat_multilabel_model.labels

        prediction_list_of_dicts.append(
            {cats[i]: score for i, score in enumerate(scores[0])}
        )

    return _score_evaluation_data_performance(
        actual_list_of_dicts=actual_list_of_dicts,
        prediction_list_of_dicts=prediction_list_of_dicts,
    )


def _score_evaluation_data_performance(
    actual_list_of_dicts: List[dict],
    prediction_list_of_dicts: List[dict],
) -> dict:
    """
    Helper function to score prediction performance on evaluation data.
    """
    actual_score_df = pd.DataFrame(actual_list_of_dicts)

    predicted_score_df = pd.DataFrame(prediction_list_of_dicts)

    predicted_score_df["max_score"] = predicted_score_df.max(axis=1)

    for column in predicted_score_df.columns[:-1]:
        predicted_score_df[column] = (
            predicted_score_df[column] == predicted_score_df["max_score"]
        ).astype(float)

    predicted_score_df = predicted_score_df.drop(columns=["max_score"])

    recall = recall_score(
        y_true=actual_score_df.to_numpy(),
        y_pred=predicted_score_df.to_numpy(),
        average="weighted",
        zero_division=0,
    )

    precision = precision_score(
        y_true=actual_score_df.to_numpy(),
        y_pred=predicted_score_df.to_numpy(),
        average="weighted",
        zero_division=0,
    )

    f1 = f1_score(
        y_true=actual_score_df.to_numpy(),
        y_pred=predicted_score_df.to_numpy(),
        average="weighted",
        zero_division=0,
    )

    return {"recall": recall, "precision": precision, "f1": f1}


def _predict(
    nlp: spacy.Language, new_resource_docs: dict, textcat_multilabel_model: Union
) -> dict:
    """
    Predict the workspace category for each resource within `new_resource_docs`
    using the trained `textcat_multilabel_model`.
    """
    resource_name_to_workspace = {}
    prediction_labels = textcat_multilabel_model.labels

    for resource_name, doc in new_resource_docs.items():
        current_document = nlp(doc)
        scores_array = textcat_multilabel_model.predict([current_document])
        resource_name_to_workspace[resource_name] = prediction_labels[
            scores_array.argmax()
        ]

    return resource_name_to_workspace
