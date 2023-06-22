"""
Helper functions for formatting managed resource drift results.
"""
from typing import Tuple
import pandas as pd
from mdutils.mdutils import MdUtils


def create_markdown_table_resource_attribute_changes(
    instance_attribute_changes_df: pd.DataFrame, markdown_file: MdUtils
) -> Tuple[MdUtils, str]:
    """Create a new Markdown table out of actor_action_count_df"""
    list_of_strings = ["Attribute", "Terraform Value", "Cloud Value"]
    for record in instance_attribute_changes_df.to_dict("records"):
        list_of_strings.extend(
            [
                record["AttributeName"],
                record["TerraformValue"],
                record["CloudValue"],
            ]
        )

    new_table_str = markdown_file.new_table(
        columns=3,
        rows=len(instance_attribute_changes_df) + 1,
        text=list_of_strings,
        text_align="center",
    )
    return markdown_file, new_table_str


def create_managed_drift_markdown(
    managed_drift_df: pd.DataFrame, markdown_file: MdUtils
) -> MdUtils:
    """Create structured tables of managed drift data."""
    for state_file in managed_drift_df["StateFileName"].unique():
        markdown_file.new_header(
            level=2,
            title=f"State File `{state_file}`",
            add_table_of_contents="n",
        )

        current_state_file_df = managed_drift_df[
            managed_drift_df["StateFileName"] == state_file
        ].sort_values("ResourcePath")

        for resource in current_state_file_df["ResourcePath"].unique():
            markdown_file.new_header(
                level=3, title=f"Resource: {resource}", add_table_of_contents="n"
            )
            for instance_id in current_state_file_df["InstanceID"].unique():
                instance_attribute_changes_df = current_state_file_df[
                    (current_state_file_df["InstanceID"] == instance_id)
                    & (current_state_file_df["ResourcePath"] == resource)
                ]
                if len(instance_attribute_changes_df) == 0:
                    continue
                actor = instance_attribute_changes_df["RecentActor"].unique()[0]
                actor = "Not Known" if actor is "" else actor

                timestamp = instance_attribute_changes_df[
                    "RecentActionTimestamp"
                ].unique()[0]
                timestamp = "Not Known" if timestamp is "" else timestamp

                markdown_file.new_line(f"**Instance ID**: `{instance_id}`")
                markdown_file.new_line("")
                markdown_file.new_line(
                    f"**Most Recent Non-Terraform Actor**: `{actor}`"
                )
                markdown_file.new_line(f"**Most Recent Action Date**: `{timestamp}`")
                markdown_file.new_line("")
                markdown_file.new_line(f"- [ ] Completed")
                markdown_file.new_line("")

                markdown_file, _ = create_markdown_table_resource_attribute_changes(
                    instance_attribute_changes_df=instance_attribute_changes_df,
                    markdown_file=markdown_file,
                )

    return markdown_file
