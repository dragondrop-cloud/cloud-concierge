"""
Helper functions for formatting security scan results.
"""
from typing import List
import pandas as pd
from mdutils.mdutils import MdUtils


def security_scan_to_df(list_of_dicts: List[dict]) -> pd.DataFrame:
    """Converts an input list of dicts representing a security scan"""
    output_df = pd.DataFrame(list_of_dicts)

    concise_output_df = output_df.drop(
        columns=[
            "rule_id",
            "long_id",
            "rule_provider",
            "rule_service",
            "warning",
            "status",
            "location",
        ]
    )
    data_df = concise_output_df.sort_values(by=["id", "severity"]).reset_index(
        drop=True
    )
    data_df["severity"] = data_df["severity"].apply(
        lambda x: x + (8 - len(x)) * " " if len(x) < 8 else x
    )

    return data_df


def create_markdown_table_security_scans(
    security_df: pd.DataFrame, markdown_file: MdUtils
) -> MdUtils:
    """Create a new Markdown table out of security_df"""
    resources_ids = security_df["id"].unique()

    if len(resources_ids) > 0:
        for resource in resources_ids:
            markdown_file.new_line(f"**Instance ID**: `{resource}`")
            current_security_df = security_df[security_df["id"] == resource]
            list_of_strings = [
                "Rule Description",
                "Severity",
                "Resolution",
                "Doc Links",
            ]
            for record in current_security_df.to_dict("records"):
                list_of_strings.extend(
                    [
                        record["rule_description"],
                        record["severity"],
                        record["resolution"],
                        f'[Rule]({record["links"][0]}), [Tf Doc]({record["links"][1]})'
                        if len(record["links"]) > 1
                        else f'[Rule]({record["links"][0]})',
                    ]
                )

            _ = markdown_file.new_table(
                columns=4,
                rows=len(current_security_df) + 1,
                text=list_of_strings,
                text_align="center",
            )
    return markdown_file
