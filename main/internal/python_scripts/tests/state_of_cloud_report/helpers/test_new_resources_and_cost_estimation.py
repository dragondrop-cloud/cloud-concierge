"""
Unit tests for helpers in estimating the cost of cloud resources
within the state of cloud report.
"""
from unittest import TestCase
import pandas as pd
from mdutils.mdutils import MdUtils
from main.internal.python_scripts.state_of_cloud_report.helpers.new_resources_and_cost_estimation import (
    create_markdown_table_new_resources,
    process_new_resources,
    _calculate_aggregate_costs_across_scan,
    _dataframe_from_cost_estimates_json,
    _uncontrolled_cost_by_div_by_type,
    _query_sort_and_clip_grouped_data,
)


def _create_baseline_expected_df(is_new_resource: bool = True) -> pd.DataFrame:
    return pd.DataFrame(
        [
            {
                "cost_component": "SQL instance (db-f1-micro, zonal)",
                "is_usage_based": False,
                "monthly_cost": 7.665,
                "monthly_quantity": 730.0,
                "price": "hours",
                "resource_name": "google_sql_database_instance.tfer--outside-of-terraform-control-db",
                "sub_resource_name": "",
                "unit": "hours",
                "provider": "google",
                "resource_type": "google_sql_database_instance",
                "is_new_resource": is_new_resource,
            },
            {
                "cost_component": "Storage (SSD, zonal)",
                "is_usage_based": False,
                "monthly_cost": 1.7,
                "monthly_quantity": 10.0,
                "price": "GB",
                "resource_name": "google_sql_database_instance.tfer--outside-of-terraform-control-db",
                "sub_resource_name": "",
                "unit": "GB",
                "provider": "google",
                "resource_type": "google_sql_database_instance",
                "is_new_resource": is_new_resource,
            },
        ]
    )


def test_dataframe_from_cost_estimates_json():
    """Unit test for _dataframe_from_divisions_to_cost_estimates_dict"""
    input_cost_estimates = [
            {
                "cost_component": "SQL instance (db-f1-micro, zonal)",
                "is_usage_based": False,
                "monthly_cost": "7.665",
                "monthly_quantity": "730",
                "price": "hours",
                "resource_name": "google_sql_database_instance.tfer--outside-of-terraform-control-db",
                "sub_resource_name": "",
                "unit": "hours",
                "provider": "google",
                "resource_type": "google_sql_database_instance",
            },
            {
                "cost_component": "Storage (SSD, zonal)",
                "is_usage_based": False,
                "monthly_cost": "1.7",
                "monthly_quantity": "10",
                "price": "GB",
                "resource_name": "google_sql_database_instance.tfer--outside-of-terraform-control-db",
                "sub_resource_name": "",
                "unit": "GB",
                "provider": "google",
                "resource_type": "google_sql_database_instance",
            },
        ]

    input_new_resources = {
        "google-dragondrop-dev.google_sql_database.tfer--outside-of-terraform-control-db-postgres": "terraform name of tfer  outs",
        "google-dragondrop-dev.google_sql_database_instance.tfer--outside-of-terraform-control-db": "ter. ",
        "google-dragondrop-dev.google_storage_bucket.tfer--testing-out-this-bucket": "terraform name ",
    }

    expected_output_df = _create_baseline_expected_df()

    output_df = _dataframe_from_cost_estimates_json(
        cost_estimates_json=input_cost_estimates,
        new_resources=input_new_resources,
    )
    pd.testing.assert_frame_equal(expected_output_df, output_df)


def test_calculate_aggregate_costs_across_scan():
    """Unit test for _calculate_aggregate_costs_across_scan"""
    # all uncontrolled monthly costs
    input_df = _create_baseline_expected_df()

    output_df = _calculate_aggregate_costs_across_scan(input_df)

    expected_output_df = pd.DataFrame(
        [
            {
                "provider": "google",
                "Uncontrolled Resources Monthly Cost": "$9.36",
                "Terraform Controlled Resources Monthly Cost": "$0.0",
                "Total Cost": "$9.36",
            }
        ]
    )

    pd.testing.assert_frame_equal(output_df, expected_output_df)

    # all controlled monthly costs
    input_df = _create_baseline_expected_df(is_new_resource=False)

    output_df = _calculate_aggregate_costs_across_scan(input_df)

    expected_output_df = pd.DataFrame(
        [
            {
                "Uncontrolled Resources Monthly Cost": "$0.0",
                "provider": "google",
                "Terraform Controlled Resources Monthly Cost": "$9.36",
                "Total Cost": "$9.36",
            }
        ]
    )

    pd.testing.assert_frame_equal(output_df, expected_output_df)


def test_uncontrolled_cost_by_div_by_type():
    """Unit test for _uncontrolled_cost_by_div_by_type"""
    input_df = _create_baseline_expected_df()

    expected_output_df = pd.DataFrame(
        [
            {
                "resource_type": "google_sql_database_instance",
                "num_cost_components": 2,
                "monthly_cost": "$9.36",
                "is_usage_based": False,
            }
        ]
    )

    output_df = _uncontrolled_cost_by_div_by_type(input_df)

    pd.testing.assert_frame_equal(expected_output_df, output_df)


def test_create_markdown_table_new_resources():
    """Unit test for create_new_markdown_table()"""
    case = TestCase()

    input_current_resource_count_df = pd.DataFrame(
        [{"type": "abc", "num_resources": 12}, {"type": "def", "num_resources": 8}]
    )

    _, markdown_string = create_markdown_table_new_resources(
        current_resource_count_df=input_current_resource_count_df,
        column="type",
        markdown_file=MdUtils("", ""),
    )

    expected_markdown_string = (
        "\n|Type|# Resources|\n| :---: | :---: |\n|abc|12|\n|def|8|\n"
    )

    case.assertEqual(expected_markdown_string, markdown_string)


def test_process_new_resource():
    """
    Unit test for the process_new_resources helper function.
    """
    input_new_resources = {
        "provider_resource_type.name": "asdasdasd",
        "provider_resource_type.name2": "asdasdasd",
        "provider2_resource_type.name": "asadasd",
        "provider2_resource_type.name2": "asdasdsd",
    }

    output_dict = process_new_resources(new_resources=input_new_resources)

    # Testing provider df
    expected_provider_df = pd.DataFrame(
        [
            {"provider": "provider", "num_resources": 2},
            {"provider": "provider2", "num_resources": 2},
        ]
    )
    pd.testing.assert_frame_equal(expected_provider_df, output_dict["provider_df"])

    # Testing provider_by_type_df
    expected_provider_by_type_df = pd.DataFrame(
        [
            {
                "provider": "provider",
                "type": "provider_resource_type",
                "num_resources": 2,
            },
            {
                "provider": "provider2",
                "type": "provider2_resource_type",
                "num_resources": 2,
            },
        ]
    )
    pd.testing.assert_frame_equal(
        expected_provider_by_type_df, output_dict["provider_by_type_df"]
    )

    # Testing provider_by_division_df
    expected_provider_by_division_df = pd.DataFrame(
        [
            {"provider": "provider", "num_resources": 2},
            {"provider": "provider2", "num_resources": 2},
        ]
    )
    pd.testing.assert_frame_equal(
        expected_provider_by_division_df, output_dict["provider_df"]
    )


def test_query_sort_and_clip_grouped_data():
    """
    _query_sort_and_clip_grouped_data()
    """
    input_grouped_df = pd.DataFrame(
        [
            {"provider": "google", "num_resources": 1},
            {"provider": "google", "num_resources": 2},
            {"provider": "google", "num_resources": 3},
            {"provider": "google", "num_resources": 4},
            {"provider": "google", "num_resources": 5},
            {"provider": "google", "num_resources": 6},
            {"provider": "google", "num_resources": 7},
            {"provider": "google", "num_resources": 8},
            {"provider": "google", "num_resources": 9},
            {"provider": "google", "num_resources": 10},
            {"provider": "google", "num_resources": 11},
            {"provider": "google", "num_resources": 12},
            {"provider": "aws", "num_resources": 1},
            {"provider": "aws", "num_resources": 2},
        ]
    )

    expected_output_df = pd.DataFrame(
        [
            {"provider": "google", "num_resources": 12},
            {"provider": "google", "num_resources": 11},
            {"provider": "google", "num_resources": 10},
            {"provider": "google", "num_resources": 9},
            {"provider": "google", "num_resources": 8},
            {"provider": "google", "num_resources": 7},
            {"provider": "google", "num_resources": 6},
            {"provider": "google", "num_resources": 5},
            {"provider": "google", "num_resources": 4},
            {"provider": "google", "num_resources": 3},
        ]
    )

    output_df, count = _query_sort_and_clip_grouped_data(
        grouped_df=input_grouped_df, current_provider="google"
    )

    pd.testing.assert_frame_equal(output_df, expected_output_df)
    assert count == 10

    _, count = _query_sort_and_clip_grouped_data(
        grouped_df=input_grouped_df, current_provider="aws"
    )
    assert count == 2
