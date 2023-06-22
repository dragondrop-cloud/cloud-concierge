"""
Unit tests for functions within the backend/nlpengine/main.py file.
"""
from unittest import TestCase
from random import seed

from driftmitigation.python_scripts.nlpengine.main import (
    _create_gold_dict,
    _doc_to_example_text_list,
    _join_text_components,
    _split_into_train_and_evaluation_data,
    _score_evaluation_data_performance,
)


def test_create_examples_from_doc():
    """
    Unit test for helper functions of _create_examples_from_doc()
    """
    case = TestCase()

    input_cat = "cat_1"
    input_categories = ["cat_1", "cat_2", "cat_3", "cat_4"]
    input_doc = "sent 1. sent 2. sent 3. "

    expected_output_text_list = ["sent 1", " sent 2", " sent 3"]
    output_text_list = _doc_to_example_text_list(doc=input_doc)

    case.assertListEqual(output_text_list, expected_output_text_list)

    expected_gold_dict = {
        "cats": {
            "cat_1": 1.0,
            "cat_2": 0.0,
            "cat_3": 0.0,
            "cat_4": 0.0,
        }
    }

    output_gold_dict = _create_gold_dict(
        current_cat=input_cat, categories=input_categories
    )

    case.assertDictEqual(expected_gold_dict, output_gold_dict)

    expected_output_combined_text = "sent 1. sent 2."
    output_combined_text = _join_text_components(
        example_text_list=expected_output_text_list,
        end_index=2,
    )

    case.assertEqual(expected_output_combined_text, output_combined_text)


def test_split_into_train_and_evaluation_data():
    """Unit test for _split_into_train_and_evaluation_data"""
    case = TestCase()

    seed(42)

    input = [1, 2, 3, 4, 5]

    expected_output = {
        "train": [4, 2, 3, 5],
        "evaluate": [1],
    }

    output = _split_into_train_and_evaluation_data(input)

    case.assertDictEqual(expected_output, output)


def test_score_evaluation_data_performance():
    """Unit test for _score_evaluation_data_performance"""
    case = TestCase()

    input_actual_list_of_dicts = [
        {"cat_1": 0, "cat_2": 0, "cat_3": 1},
        {"cat_1": 0, "cat_2": 1, "cat_3": 0},
        {"cat_1": 1, "cat_2": 0, "cat_3": 0},
    ]

    input_prediction_list_of_dicts = [
        {"cat_1": 0.02, "cat_2": 0.01, "cat_3": 0.73},
        {"cat_1": 0.32, "cat_2": 0.76, "cat_3": 0.12},
        {"cat_1": 0.98, "cat_2": 0.02, "cat_3": 0.45},
    ]

    output = _score_evaluation_data_performance(
        actual_list_of_dicts=input_actual_list_of_dicts,
        prediction_list_of_dicts=input_prediction_list_of_dicts,
    )

    expected_output = {
        "recall": 1.0,
        "precision": 1.0,
        "f1": 1.0,
    }

    case.assertDictEqual(output, expected_output)
