"""
Helper functions for estimating the cost of cloud resources
within the state of cloud report.
"""
from typing import Tuple
import pandas as pd
from mdutils.mdutils import MdUtils


def create_markdown_table_new_resources(
    current_resource_count_df: pd.DataFrame,
    column: str,
    markdown_file: MdUtils,
) -> Tuple[MdUtils, str]:
    """Create a new Markdown table out of the current_resource_count_df"""
    list_of_strings = [column.title(), "# Resources"]
    for record in current_resource_count_df.to_dict("records"):
        list_of_strings.extend([record[column], record["num_resources"]])

    new_table_str = markdown_file.new_table(
        columns=2,
        rows=len(current_resource_count_df) + 1,
        text=list_of_strings,
        text_align="center",
    )
    return markdown_file, new_table_str


def process_new_resources(new_resources: dict) -> dict:
    """
    Parses input resources file and returns relevant, processed versions of the data as dataframe.

    Specifically, produces three data frames:
    provider_df, provider_by_type_df, provider_by_division_df

    Each with different groupings with the critical variable being num_resources discovered for that category.
    """
    list_of_dicts = []
    for resource_key, _ in new_resources.items():
        division, resource_type, _ = resource_key.split(".")
        provider = resource_type.split("_")[0]
        current_resource_dict = {
            "division": division,
            "type": resource_type,
            "provider": provider,
        }
        list_of_dicts.append(current_resource_dict)

    new_resources_df = pd.DataFrame(list_of_dicts)

    count_by_provider_df = (
        new_resources_df.groupby(by=["provider"])
        .agg(
            num_resources=pd.NamedAgg(column="provider", aggfunc="count"),
        )
        .reset_index()
        .sort_values(by=["provider"], ascending=True)
        .reset_index(drop=True)
    )

    count_by_provider_by_type_df = (
        new_resources_df.groupby(by=["provider", "type"])
        .agg(
            num_resources=pd.NamedAgg(column="provider", aggfunc="count"),
        )
        .reset_index()
        .sort_values(by=["provider", "type"], ascending=True)
        .reset_index(drop=True)
    )

    count_by_provider_by_division_df = (
        new_resources_df.groupby(by=["provider", "division"])
        .agg(num_resources=pd.NamedAgg(column="provider", aggfunc="count"))
        .reset_index()
        .sort_values(by=["provider", "division"], ascending=True)
        .reset_index(drop=True)
    )

    return {
        "provider_df": count_by_provider_df,
        "provider_by_type_df": count_by_provider_by_type_df,
        "provider_by_division_df": count_by_provider_by_division_df,
    }


def single_provider_new_resources_by_division_tabular_output(
    markdown_file: MdUtils,
    current_provider: str,
    by_division_df: pd.DataFrame,
) -> MdUtils:
    """
    Create graphical output (pie chart and table) for resource counts by division
    for single provider.
    """
    (
        current_by_division_df,
        number_of_top_types_count,
    ) = _query_sort_and_clip_grouped_data(
        grouped_df=by_division_df, current_provider=current_provider
    )

    markdown_file.new_line()
    markdown_file, _ = create_markdown_table_new_resources(
        current_resource_count_df=current_by_division_df,
        column="division",
        markdown_file=markdown_file,
    )

    return markdown_file


def _query_sort_and_clip_grouped_data(
    grouped_df: pd.DataFrame, current_provider: str
) -> Tuple[pd.DataFrame, int]:
    current_grouped_df = grouped_df.query(f"provider == '{current_provider}'")

    current_grouped_df = current_grouped_df.sort_values(
        by=["num_resources"], ascending=False
    ).reset_index(drop=True)

    number_of_top_types_count = min(len(current_grouped_df), 10)
    current_grouped_df = current_grouped_df.loc[: number_of_top_types_count - 1, :]
    return current_grouped_df, number_of_top_types_count


# Cost estimation functions
def create_new_resource_tabular_breakdowns_with_cost(
    markdown_file: MdUtils,
    resource_count_dict_of_dfs: dict,
    cost_by_provider_by_type_df: pd.DataFrame,
) -> MdUtils:
    """
    Function that coordinates the creation of all tabular breakdown
    of cloud resources identified by dragondrop.
    """
    provider_breakdown_df = resource_count_dict_of_dfs["provider_df"]
    provider_to_resource_totals = provider_breakdown_df.to_dict("records")

    markdown_file, _ = create_markdown_table_new_resources(
        markdown_file=markdown_file,
        column="provider",
        current_resource_count_df=provider_breakdown_df,
    )

    by_division_df = resource_count_dict_of_dfs["provider_by_division_df"]
    by_type_df = resource_count_dict_of_dfs["provider_by_type_df"]

    # Creating outputs by provider
    for provider in provider_to_resource_totals:
        current_provider = provider["provider"]
        markdown_file.new_header(
            level=2,
            title=current_provider.title()
            if current_provider.lower() != "aws"
            else "AWS",
            add_table_of_contents="n",
        )

        markdown_file = single_provider_new_resources_by_division_tabular_output(
            markdown_file=markdown_file,
            current_provider=current_provider,
            by_division_df=by_division_df,
        )

        markdown_file = single_provider_costs_by_type_tabular_output(
            markdown_file=markdown_file,
            current_provider=current_provider,
            by_type_df=by_type_df,
            cost_by_provider_by_type_df=cost_by_provider_by_type_df,
        )

        return markdown_file


def create_markdown_table_cost_by_resource_type(
    current_by_type_df: pd.DataFrame,
    resource_cost_by_provider_by_type_df: pd.DataFrame,
    markdown_file: MdUtils,
) -> Tuple[MdUtils, str]:
    """Create a new Markdown table out of cost_summary_df"""
    if not resource_cost_by_provider_by_type_df.empty:
        list_of_strings = [
            "Type",
            "# Resources",
            "Cost Components",
            "Monthly Cost",
            "Usage Based*",
        ]

        output_df = current_by_type_df.merge(
            resource_cost_by_provider_by_type_df.rename(
                columns={"resource_type": "type"}
            ),
            how="left",
        ).fillna("No Charge")

        for record in output_df.to_dict("records"):
            list_of_strings.extend(
                [
                    record["type"],
                    record["num_resources"],
                    str(record["num_cost_components"]).split(".")[0],
                    record["monthly_cost"],
                    record["is_usage_based"],
                ]
            )

        new_table_str = markdown_file.new_table(
            columns=5,
            rows=len(output_df) + 1,
            text=list_of_strings,
            text_align="center",
        )
    else:
        list_of_strings = ["Type", "# Resources"]

        for record in current_by_type_df.to_dict("records"):
            list_of_strings.extend(
                [
                    record["type"],
                    record["num_resources"],
                ]
            )

        new_table_str = markdown_file.new_table(
            columns=2,
            rows=len(current_by_type_df) + 1,
            text=list_of_strings,
            text_align="center",
        )
    return markdown_file, new_table_str


def single_provider_costs_by_type_tabular_output(
    markdown_file: MdUtils,
    current_provider: str,
    by_type_df: pd.DataFrame,
    cost_by_provider_by_type_df: pd.DataFrame,
) -> MdUtils:
    """
    Create tabular output for New resource counts by
    resource type for a single provider.
    """
    current_by_type_df, number_of_top_types_count = _query_sort_and_clip_grouped_data(
        grouped_df=by_type_df, current_provider=current_provider
    )

    markdown_file.new_line()
    markdown_file, _ = create_markdown_table_cost_by_resource_type(
        current_by_type_df=current_by_type_df,
        resource_cost_by_provider_by_type_df=cost_by_provider_by_type_df,
        markdown_file=markdown_file,
    )
    return markdown_file


def process_pricing_data(
    divisions_to_cost_estimates: dict,
    new_resources: dict,
) -> dict:
    """
    Process pricing data in the following format:
    {
        'google-dragondrop-dev': [
            {
                'cost_component': 'SQL instance (db-f1-micro, zonal)',
                'is_usage_based': False,
                'monthly_cost': '7.665',
                'monthly_quantity': '730',
                'price': 'hours',
                'resource_name': 'google_sql_database_instance.tfer--outside-of-terraform-control-db',
                'sub_resource_name': '',
                'unit': 'hours',
                'provider': 'google',
                'division': 'google-dragondrop-dev',
                'resource_type': 'google_sql_database_instance'
            },
            ....
            {
                'cost_component': 'SQL instance (db-f1-micro, zonal)',
                ...
                'resource_type': 'google_sql_database_instance'
            }
        ]
    }
    into two dataframes that look as follows:
    1)
    provider | Uncontrolled Resources Monthly Cost | Terraform Controlled Resources Monthly Cost |
    google   |                          $16.665    |                          $16.665            |
    aws      |                          $16.665    |                          $16.665            |

    2)
    provider | division              |   resource_type              | num_cost_components | monthly_cost | is_usage_based |
    google   | google-dragondrop-dev | google_sql_database_instance | 4                   |   $16.665    |   False        |
    google   | google-dragondrop-dev | google_storage_bucket        | 8                   |   $0.0*      |      True      |
    """
    df = _dataframe_from_divisions_to_cost_estimates_dict(
        divisions_to_cost_estimates=divisions_to_cost_estimates,
        new_resources=new_resources,
    )

    combined_cost_summary_df = _calculate_aggregate_costs_across_scan(df)

    uncontrolled_cost_by_div_by_type_df = _uncontrolled_cost_by_div_by_type(df=df)

    return {
        "cost_summary": combined_cost_summary_df,
        "uncontrolled_cost_by_div_by_type_df": uncontrolled_cost_by_div_by_type_df,
    }


def _dataframe_from_divisions_to_cost_estimates_dict(
    divisions_to_cost_estimates: dict, new_resources: dict
) -> pd.DataFrame:
    """
    Convert the divisions_to_cost_estimates dictionary into a dataframe with some basic feature engineering
    which is leveraged by all downstream operations.
    """
    complete_data_list_of_dicts = []

    for div, list_of_dicts in divisions_to_cost_estimates.items():

        provider = div.split("-")[0]
        for dict in list_of_dicts:
            dict["provider"] = provider
            dict["division"] = div
            dict["resource_type"] = dict["resource_name"].split(".")[0]

        complete_data_list_of_dicts.extend(list_of_dicts)

    df = pd.DataFrame(complete_data_list_of_dicts)

    # Is the resource new or not?
    resource_name_set = set()
    for name in new_resources.keys():
        resource_name_set.add(".".join(name.split(".")[1:]))

    df["is_new_resource"] = df["resource_name"].isin(resource_name_set)
    df[["monthly_cost", "monthly_quantity"]] = (
        df[["monthly_cost", "monthly_quantity"]].replace("", 0).astype(float)
    )

    return df


def _calculate_aggregate_costs_across_scan(df: pd.DataFrame) -> pd.DataFrame:
    """
    Calculate aggregate cloud costs by provider split into whether the costs are controlled by Terraform
    or not.
    """
    grouped_by_provider_tf_status_df = (
        df.groupby(by=["provider", "is_new_resource"])
        .agg(Cost=pd.NamedAgg(aggfunc="sum", column="monthly_cost"))
        .reset_index()
    )

    grouped_uncontrolled_df = (
        grouped_by_provider_tf_status_df.query("is_new_resource == False")
        .rename(columns={"Cost": "Terraform Controlled Resources Monthly Cost"})
        .drop(columns=["is_new_resource"])
    )
    grouped_controlled_df = (
        grouped_by_provider_tf_status_df.query("is_new_resource == True")
        .rename(columns={"Cost": "Uncontrolled Resources Monthly Cost"})
        .drop(columns=["is_new_resource"])
    )

    combined_cost_summary_df = grouped_controlled_df.merge(
        grouped_uncontrolled_df, how="outer", on="provider"
    ).fillna(0)
    assert len(combined_cost_summary_df) <= len(grouped_controlled_df) + len(
        grouped_uncontrolled_df
    )

    combined_cost_summary_df[
        [
            "Terraform Controlled Resources Monthly Cost",
            "Uncontrolled Resources Monthly Cost",
        ]
    ] = combined_cost_summary_df[
        [
            "Terraform Controlled Resources Monthly Cost",
            "Uncontrolled Resources Monthly Cost",
        ]
    ].round(
        2
    )

    combined_cost_summary_df["Total Cost"] = (
        combined_cost_summary_df["Terraform Controlled Resources Monthly Cost"]
        + combined_cost_summary_df["Uncontrolled Resources Monthly Cost"]
    )

    for cost_column in [
        "Total Cost",
        "Terraform Controlled Resources Monthly Cost",
        "Uncontrolled Resources Monthly Cost",
    ]:
        combined_cost_summary_df[cost_column] = "$" + combined_cost_summary_df[
            cost_column
        ].astype(str)

    return combined_cost_summary_df


def _uncontrolled_cost_by_div_by_type(df: pd.DataFrame) -> pd.DataFrame:
    """Calculated uncontrolled cost by division and by resource type"""
    uncontrolled_cost_by_div_by_type_df = (
        df.query("is_new_resource == True")
        .groupby(by=["provider", "resource_type"])
        .agg(
            num_cost_components=pd.NamedAgg(aggfunc="nunique", column="cost_component"),
            monthly_cost=pd.NamedAgg(aggfunc="sum", column="monthly_cost"),
            is_usage_based=pd.NamedAgg(aggfunc="first", column="is_usage_based"),
        )
        .reset_index()
    )

    uncontrolled_cost_by_div_by_type_df[
        "monthly_cost"
    ] = "$" + uncontrolled_cost_by_div_by_type_df["monthly_cost"].round(2).astype(str)
    uncontrolled_cost_by_div_by_type_df.loc[
        uncontrolled_cost_by_div_by_type_df["is_usage_based"], "monthly_cost"
    ] = (
        uncontrolled_cost_by_div_by_type_df.loc[
            uncontrolled_cost_by_div_by_type_df["is_usage_based"], "monthly_cost"
        ]
        + "*"
    )

    return uncontrolled_cost_by_div_by_type_df


def create_markdown_table_cost_summary(
    cost_summary_df: pd.DataFrame, markdown_file: MdUtils
) -> MdUtils:
    """Create a new Markdown table out of cost_summary_df"""
    list_of_strings = [
        "Cloud Provider",
        "Uncontrolled Resources Cost",
        "Terraform Controlled Resources Cost",
    ]

    if not cost_summary_df.empty:
        for record in cost_summary_df.to_dict("records"):
            list_of_strings.extend(
                [
                    record["provider"],
                    record["Uncontrolled Resources Monthly Cost"],
                    record["Terraform Controlled Resources Monthly Cost"],
                ]
            )

    _ = markdown_file.new_table(
        columns=3,
        rows=len(cost_summary_df) + 1,
        text=list_of_strings,
        text_align="center",
    )
    return markdown_file
