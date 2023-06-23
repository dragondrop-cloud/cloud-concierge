"""
Helper functions for formatting cloud actor identification results.
"""
from typing import Tuple
import pandas as pd
from mdutils.mdutils import MdUtils


def create_markdown_table_cloud_actor_summary(
    actor_action_count_df: pd.DataFrame, markdown_file: MdUtils
) -> Tuple[MdUtils, str]:
    """Create a new Markdown table out of actor_action_count_df"""
    list_of_strings = ["Cloud Provider", "Division", "Actor", "Action", "Count"]
    for record in actor_action_count_df.to_dict("records"):
        list_of_strings.extend(
            [
                record["Cloud Provider"],
                record["Cloud Division"],
                record["Cloud Actor"],
                record["Action"],
                record["Count"],
            ]
        )

    new_table_str = markdown_file.new_table(
        columns=5,
        rows=len(actor_action_count_df) + 1,
        text=list_of_strings,
        text_align="center",
    )
    return markdown_file, new_table_str


def process_cloud_actor_actions(
    resources_to_cloud_actions: dict,
) -> pd.DataFrame:
    """
    Function take in cloud actor actions in dictionary of the form:
    {
        'google': {
            'dragondrop-dev': {
                'google_storage_bucket.testing_out_this_bucket': {
                    'creation': {
                        'actor': 'goodman.benjamin@dragondrop.cloud', 'timestamp': '2023-02-25'
                    },
                    'modified': {
                        'actor': 'goodman.benjamin@dragondrop.cloud', 'timestamp': '2023-03-08'
                    }
                }
            }
        }
    }

    And returns a pandas dataframe with the following columns:
    | Cloud Actor | Cloud Division | Cloud Actor | Action | Count |
    """
    list_of_dicts = []
    for provider, division_dict in resources_to_cloud_actions.items():
        for division, resource_dict in division_dict.items():
            for resource, actions in resource_dict.items():
                if "creation" in actions:
                    creation_actor = actions["creation"]["actor"]
                    list_of_dicts.append(
                        {
                            "Cloud Provider": provider,
                            "Cloud Division": division,
                            "Cloud Actor": creation_actor,
                            "Action": "Create Resource",
                            "resource": resource,
                        }
                    )
                if "modified" in actions:
                    modified_actor = actions["modified"]["actor"]
                    list_of_dicts.append(
                        {
                            "Cloud Provider": provider,
                            "Cloud Division": division,
                            "Cloud Actor": modified_actor,
                            "Action": "Modify Resource",
                            "resource": resource,
                        }
                    )

    cloud_actors_df = pd.DataFrame(list_of_dicts)
    actor_action_count_df = (
        cloud_actors_df.groupby(
            by=["Cloud Provider", "Cloud Division", "Cloud Actor", "Action"]
        )
        .agg(Count=pd.NamedAgg(column="resource", aggfunc="nunique"))
        .reset_index()
    )
    actor_action_count_df = actor_action_count_df.sort_values(
        by=["Cloud Provider", "Cloud Division", "Count"], ascending=[True, True, False]
    ).reset_index(drop=True)
    return actor_action_count_df
