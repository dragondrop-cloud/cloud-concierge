"""
Unit tests for helpers in security scanning formatting.
"""
import pandas as pd
from main.internal.python_scripts.state_of_cloud_report.helpers.security_scanning import (
    security_scan_to_df,
)


def test_division_to_security_scan_to_df_dict():
    """
    Unit test for division_to_security_scan_to_df_dict
    """
    input_security_results_list = [
        {
            "id": "arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302",
            "rule_id": "AVD-AWS-0053",
            "long_id": "aws-elb-alb-not-public",
            "rule_description": "Load balancer is exposed to the internet.",
            "rule_provider": "aws",
            "rule_service": "elb",
            "impact": "The load balancer is exposed on the internet",
            "resolution": "Switch to an internal load balancer or add a tfsec ignore",
            "links": [
                "https://aquasecurity.github.io/tfsec/v1.28.1/checks/aws/elb/alb-not-public/",
                "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb",
            ],
            "description": "Load balancer is exposed publicly.",
            "severity": "HIGH",
            "warning": False,
            "status": 0,
            "resource": "aws_lb.tfer--tf-managed-demo-alb",
            "location": {"file_name": "", "start_line": 11, "end_line": 11},
        },
        {
            "id": "arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302",
            "rule_id": "AVD-AWS-0053",
            "long_id": "aws-elb-alb-not-public",
            "rule_description": "Load balancer is exposed to the internet.",
            "rule_provider": "aws",
            "rule_service": "elb",
            "impact": "The load balancer is exposed on the internet",
            "resolution": "Switch to an internal load balancer or add a tfsec ignore",
            "links": [
                "https://aquasecurity.github.io/tfsec/v1.28.1/checks/aws/elb/alb-not-public/",
                "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb",
            ],
            "description": "Load balancer is exposed publicly.",
            "severity": "CRITICAL",
            "warning": False,
            "status": 0,
            "resource": "aws_lb.tfer--tf-managed-demo-alb",
            "location": {"file_name": "", "start_line": 11, "end_line": 11},
        },
    ]

    output_df = security_scan_to_df(list_of_dicts=input_security_results_list)

    expected_output_df = pd.DataFrame(
        [
            {
                "id": "arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302",
                "rule_description": "Load balancer is exposed to the internet.",
                "impact": "The load balancer is exposed on the internet",
                "resolution": "Switch to an internal load balancer or add a tfsec ignore",
                "links": [
                    "https://aquasecurity.github.io/tfsec/v1.28.1/checks/aws/elb/alb-not-public/",
                    "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb",
                ],
                "description": "Load balancer is exposed publicly.",
                "severity": "CRITICAL",
                "resource": "aws_lb.tfer--tf-managed-demo-alb",
            },
            {
                "id": "arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302",
                "rule_description": "Load balancer is exposed to the internet.",
                "impact": "The load balancer is exposed on the internet",
                "resolution": "Switch to an internal load balancer or add a tfsec ignore",
                "links": [
                    "https://aquasecurity.github.io/tfsec/v1.28.1/checks/aws/elb/alb-not-public/",
                    "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb",
                ],
                "description": "Load balancer is exposed publicly.",
                "severity": "HIGH    ",
                "resource": "aws_lb.tfer--tf-managed-demo-alb",
            },
        ]
    )

    pd.testing.assert_frame_equal(output_df, expected_output_df)
