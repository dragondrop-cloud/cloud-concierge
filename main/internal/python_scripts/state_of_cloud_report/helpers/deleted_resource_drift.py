from typing import List, Tuple
from mdutils.mdutils import MdUtils


def create_deleted_resource_tabular_breakdowns_with_cost(
        markdown_file: MdUtils,
        deleted_resources_list_of_dicts: List[dict]
) -> MdUtils:
    """
    Function that coordinates the creation of all tabular breakdown
    of deleted resources identified by dragondrop.
    """
    markdown_file.new_line()
    markdown_file, _ = create_markdown_table_deleted_resources(
        deleted_resources_list_of_dicts=deleted_resources_list_of_dicts,
        markdown_file=markdown_file,
    )

    return markdown_file


def create_markdown_table_deleted_resources(
        deleted_resources_list_of_dicts: List[dict],
        markdown_file: MdUtils,
) -> Tuple[MdUtils, str]:
    """Create a new Markdown table out of deleted_resources"""
    list_of_strings = [
        "Type",
        "Name",
        "Module",
        "State File",
    ]

    for deleted_resources in deleted_resources_list_of_dicts:
        list_of_strings.extend(
            [
                deleted_resources["ResourceType"],
                deleted_resources["ResourceName"],
                deleted_resources["ModuleName"],
                deleted_resources["StateFileName"],
            ]
        )

    new_table_str = markdown_file.new_table(
        columns=4,
        rows=len(deleted_resources_list_of_dicts) + 1,
        text=list_of_strings,
        text_align="center",
    )
    return markdown_file, new_table_str
