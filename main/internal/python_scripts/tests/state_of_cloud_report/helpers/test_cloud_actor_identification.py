"""
Unit tests for state_of_cloud_report/cloud_actor_identification.py
"""
from unittest import TestCase
import pandas as pd
from mdutils.mdutils import MdUtils

from main.internal.python_scripts.state_of_cloud_report.helpers.cloud_actor_identificaton import (
    create_markdown_table_cloud_actor_summary,
    process_cloud_actor_actions,
)


def _func_create_cloud_actor_dataframe() -> pd.DataFrame:
    return pd.DataFrame(
        [
            {
                "Cloud Provider": "google",
                "Cloud Division": "dragondrop-dev",
                "Cloud Actor": "g.b@dragondrop.cloud",
                "Action": "Create Resource",
                "Count": 2,
            },
            {
                "Cloud Provider": "google",
                "Cloud Division": "dragondrop-dev",
                "Cloud Actor": "g.b@dragondrop.cloud",
                "Action": "Modify Resource",
                "Count": 1,
            },
            {
                "Cloud Provider": "google",
                "Cloud Division": "dragondrop-extra",
                "Cloud Actor": "mcguill.jimmy@dragondrop.cloud",
                "Action": "Modify Resource",
                "Count": 1,
            },
        ]
    )


def test_create_markdown_table_cloud_actor_summary():
    case = TestCase()

    input_actor_action_count_df = _func_create_cloud_actor_dataframe()

    _, output_markdown_string = create_markdown_table_cloud_actor_summary(
        actor_action_count_df=input_actor_action_count_df,
        markdown_file=MdUtils("", ""),
    )

    expected_output_markdown_string = "\n|Cloud Provider|Division|Actor|Action|Count|\n| :---: | :---: | :---: | :---: | :---: |\n|google|dragondrop-dev|g.b@dragondrop.cloud|Create Resource|2|\n|google|dragondrop-dev|g.b@dragondrop.cloud|Modify Resource|1|\n|google|dragondrop-extra|mcguill.jimmy@dragondrop.cloud|Modify Resource|1|\n"

    case.assertEqual(output_markdown_string, expected_output_markdown_string)


def test_process_cloud_actor_actions():
    resources_to_cloud_actions = {
        "google": {
            "dragondrop-dev": {
                "google_storage_bucket.testing_out_this_bucket": {
                    "creation": {
                        "actor": "g.b@dragondrop.cloud",
                        "timestamp": "2023-02-25",
                    },
                    "modified": {
                        "actor": "g.b@dragondrop.cloud",
                        "timestamp": "2023-03-08",
                    },
                },
                "google_storage_bucket.new_unique_bucket": {
                    "creation": {
                        "actor": "g.b@dragondrop.cloud",
                        "timestamp": "2023-02-25",
                    },
                },
            },
            "dragondrop-extra": {
                "google_storage_bucket.testing_out_this_bucket": {
                    "modified": {
                        "actor": "mcguill.jimmy@dragondrop.cloud",
                        "timestamp": "2023-03-08",
                    }
                },
            },
        }
    }

    output_df = process_cloud_actor_actions(
        resources_to_cloud_actions=resources_to_cloud_actions
    )

    expected_output_df = _func_create_cloud_actor_dataframe()

    pd.testing.assert_frame_equal(output_df, expected_output_df)
