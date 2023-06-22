package terraformImportMigrationGenerator

import (
	"testing"

	"github.com/Jeffail/gabs/v2"
	"github.com/stretchr/testify/assert"

	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
)

func TestGetResourceLocationFormatted_AWS_S3(t *testing.T) {
	// Given
	provider := terraformValueObjects.Provider("aws")
	resourceType := ResourceType("aws_s3_bucket")
	resourcesJSON := []byte(`{
		"name": "tfer--dragondrop-example-2",
		"instances": [
			{
				"attributes_flat": {
			  		"id": "dragondrop-example-2",
			  		"arn": "arn:aws:s3:::dragondrop-example-2"
				}
			}
		]
	}`)
	resourcesParsed, err := gabs.ParseJSON(resourcesJSON)
	assert.Nil(t, err)

	// When
	resourceFormatted, err := GetRemoteCloudReference(resourcesParsed, provider, resourceType)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "dragondrop-example-2", resourceFormatted)
}

func TestGetRemoteCloudReference_AWS(t *testing.T) {
	tests := []struct {
		name           string
		resourceType   string
		inputJSON      string
		expectedOutput string
		expectedErr    error
	}{
		{
			name:         "aws_accessanalyzer_analyzer",
			resourceType: "aws_accessanalyzer_analyzer",
			inputJSON: `{
                "name": "my-analyzer",
                "instances": [
                    {
                        "attributes_flat": {
                            "arn": "arn:aws:accessanalyzer:us-west-2:123456789012:analyzer/my-analyzer"
                        }
                    }
                ]
            }`,
			expectedOutput: "arn:aws:accessanalyzer:us-west-2:123456789012:analyzer/my-analyzer",
		},
		{
			name:         "aws_acm_certificate",
			resourceType: "aws_acm_certificate",
			inputJSON: `{
                "name": "my-certificate",
                "instances": [
                    {
                        "attributes_flat": {
                            "arn": "arn:aws:acm:us-west-2:123456789012:certificate/12345678-1234-1234-1234-123456789012"
                        }
                    }
                ]
            }`,
			expectedOutput: "arn:aws:acm:us-west-2:123456789012:certificate/12345678-1234-1234-1234-123456789012",
		},
		{
			name:         "aws_lb",
			resourceType: "aws_lb",
			inputJSON: `{
                "name": "my-load-balancer",
                "instances": [
                    {
                        "attributes_flat": {
                            "arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/1234567890abcdef"
                        }
                    }
                ]
            }`,
			expectedOutput: "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/1234567890abcdef",
		},
		{
			name:         "aws_lb_listener",
			resourceType: "aws_lb_listener",
			inputJSON: `{
                "name": "my-listener",
                "instances": [
                    {
                        "attributes_flat": {
                            "load_balancer_arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/1234567890abcdef",
                            "port": "80"
                        }
                    }
                ]
            }`,
			expectedOutput: "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/1234567890abcdef/80",
		},
		{
			name:         "aws_lb_listener_rule",
			resourceType: "aws_lb_listener_rule",
			inputJSON: `{
                "name": "my-listener-rule",
                "instances": [
                    {
                        "attributes_flat": {
                            "load_balancer_arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/1234567890abcdef",
                            "listener_arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:listener/app/my-load-balancer/1234567890abcdef/1234567890abcdef",
                            "arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:listener-rule/app/my-load-balancer/1234567890abcdef/1234567890abcdef/1234567890abcdef"
                        }
                    }
                ]
            }`,
			expectedOutput: "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/1234567890abcdef/arn:aws:elasticloadbalancing:us-west-2:123456789012:listener/app/my-load-balancer/1234567890abcdef/1234567890abcdef/arn:aws:elasticloadbalancing:us-west-2:123456789012:listener-rule/app/my-load-balancer/1234567890abcdef/1234567890abcdef/1234567890abcdef",
		},
		{
			name:         "aws_workspaces_workspace",
			resourceType: "aws_workspaces_workspace",
			inputJSON: `{
                "name": "my-workspace",
                "instances": [
                    {
                        "attributes_flat": {
                            "workspace_id": "ws-0123456789abcdef0"
                        }
                    }
                ]
            }`,
			expectedOutput: "ws-0123456789abcdef0",
		},
		{
			name:         "aws_xray_sampling_rule",
			resourceType: "aws_xray_sampling_rule",
			inputJSON: `{
                "name": "my-sampling-rule",
                "instances": [
                    {
                        "attributes_flat": {
                            "id": "abcd1234"
                        }
                    }
                ]
            }`,
			expectedOutput: "abcd1234",
		},
		{
			name:         "aws_workspaces_directory",
			resourceType: "aws_workspaces_directory",
			inputJSON: `{
                "name": "my-directory",
                "instances": [
                    {
                        "attributes_flat": {
                            "directory_id": "d-1234567890"
                        }
                    }
                ]
            }`,
			expectedOutput: "d-1234567890",
		},
		{
			name:         "aws_workspaces_ip_group",
			resourceType: "aws_workspaces_ip_group",
			inputJSON: `{
                "name": "my-ip-group",
                "instances": [
                    {
                        "attributes_flat": {
                            "group_id": "wsipg-0123456789abcdef"
                        }
                    }
                ]
            }`,
			expectedOutput: "wsipg-0123456789abcdef",
		},
		{
			name:         "aws_wafv2_web_acl",
			resourceType: "aws_wafv2_web_acl",
			inputJSON: `{
				"name": "my-web-acl",
				"instances": [
					{
						"attributes_flat": {
							"id": "webacl-0123456789abcdef"
						}
					}
				]
			}`,
			expectedOutput: "webacl-0123456789abcdef",
		},
		{
			name:         "aws_wafv2_web_acl_association",
			resourceType: "aws_wafv2_web_acl_association",
			inputJSON: `{
				"name": "my-web-acl-association",
				"instances": [
					{
						"attributes_flat": {
							"web_acl_arn": "arn:aws:wafv2:us-west-2:123456789012:global/webacl/my-web-acl/01234567-89ab-cdef-0123-456789abcdef",
							"resource_type": "ALB",
							"resource_id": "my-alb"
						}
					}
				]
			}`,
			expectedOutput: "arn:aws:wafv2:us-west-2:123456789012:global/webacl/my-web-acl/01234567-89ab-cdef-0123-456789abcdef/ALB/my-alb",
		},
		{
			name:         "aws_wafv2_web_acl_logging_configuration",
			resourceType: "aws_wafv2_web_acl_logging_configuration",
			inputJSON: `{
				"name": "my-web-acl-logging-configuration",
				"instances": [
					{
						"attributes_flat": {
							"resource_arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-alb/0123456789abcdef",
							"id": "logging"
						}
					}
				]
			}`,
			expectedOutput: "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-alb/0123456789abcdef/logging",
		},
		{
			name:         "aws_wafregional_xss_match_set",
			resourceType: "aws_wafregional_xss_match_set",
			inputJSON: `{
            "name": "my-xss-match-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "xss-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "xss-12345678",
		},
		{
			name:         "aws_wafv2_ip_set",
			resourceType: "aws_wafv2_ip_set",
			inputJSON: `{
            "name": "my-ip-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "ipset-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "ipset-12345678",
		},
		{
			name:         "aws_wafv2_regex_pattern_set",
			resourceType: "aws_wafv2_regex_pattern_set",
			inputJSON: `{
            "name": "my-regex-pattern-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "regexpatternset-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "regexpatternset-12345678",
		},
		{
			name:         "aws_wafregional_size_constraint_set",
			resourceType: "aws_wafregional_size_constraint_set",
			inputJSON: `{
				"instances": [
					{
						"attributes_flat": {
							"id": "my-id"
						}
					}
				]
			}`,
			expectedOutput: "my-id",
		},
		{
			name:         "aws_wafregional_sql_injection_match_set",
			resourceType: "aws_wafregional_sql_injection_match_set",
			inputJSON: `{
				"instances": [
					{
						"attributes_flat": {
							"id": "my-id"
						}
					}
				]
			}`,
			expectedOutput: "my-id",
		},
		{
			name:         "aws_wafregional_web_acl",
			resourceType: "aws_wafregional_web_acl",
			inputJSON: `{
				"instances": [
					{
						"attributes_flat": {
							"id": "my-id"
						}
					}
				]
			}`,
			expectedOutput: "my-id",
		},
		{
			name:         "aws_wafregional_regex_pattern_set",
			resourceType: "aws_wafregional_regex_pattern_set",
			inputJSON: `{
			"instances": [
				{
					"attributes_flat": {
						"id": "my-id"
					}
				}
			]
		}`,
			expectedOutput: "my-id",
		},
		{
			name:         "aws_wafregional_rule",
			resourceType: "aws_wafregional_rule",
			inputJSON: `{
			"instances": [
				{
					"attributes_flat": {
						"id": "my-id"
					}
				}
			]
		}`,
			expectedOutput: "my-id",
		},
		{
			name:         "aws_wafregional_rule_group",
			resourceType: "aws_wafregional_rule_group",
			inputJSON: `{
			"instances": [
				{
					"attributes_flat": {
						"id": "my-id"
					}
				}
			]
		}`,
			expectedOutput: "my-id",
		},
		{
			name:         "aws_wafregional_ipset",
			resourceType: "aws_wafregional_ipset",
			inputJSON: `{
            "name": "my-ip-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "ipset-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "ipset-12345678",
		},
		{
			name:         "aws_wafregional_rate_based_rule",
			resourceType: "aws_wafregional_rate_based_rule",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-id"
                    }
                }
            ]
        }`,
			expectedOutput: "my-id",
		},
		{
			name:         "aws_wafregional_regex_match_set",
			resourceType: "aws_wafregional_regex_match_set",
			inputJSON: `{
            "name": "my-regex-match-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "regexmatchset-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "regexmatchset-12345678",
		},
		{
			name:         "aws_waf_xss_match_set",
			resourceType: "aws_waf_xss_match_set",
			inputJSON: `{
            "name": "my-xss-match-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "abcd1234"
                    }
                }
            ]
        }`,
			expectedOutput: "abcd1234",
		},
		{
			name:         "aws_wafregional_byte_match_set",
			resourceType: "aws_wafregional_byte_match_set",
			inputJSON: `{
            "name": "my-byte-match-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "bytematchset-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "bytematchset-12345678",
		},
		{
			name:         "aws_wafregional_geo_match_set",
			resourceType: "aws_wafregional_geo_match_set",
			inputJSON: `{
            "name": "my-geo-match-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "geomatchset-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "geomatchset-12345678",
		},
		{
			name:         "aws_waf_size_constraint_set",
			resourceType: "aws_waf_size_constraint_set",
			inputJSON: `{
            "name": "my-size-constraint-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "size-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "size-12345678",
		},
		{
			name:         "aws_waf_sql_injection_match_set",
			resourceType: "aws_waf_sql_injection_match_set",
			inputJSON: `{
            "name": "my-sql-injection-match-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "sqlinjectionmatchset-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "sqlinjectionmatchset-12345678",
		},
		{
			name:         "aws_waf_web_acl",
			resourceType: "aws_waf_web_acl",
			inputJSON: `{
            "name": "my-web-acl",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "webacl-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "webacl-12345678",
		},
		{
			name:         "aws_waf_regex_pattern_set",
			resourceType: "aws_waf_regex_pattern_set",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-regex-pattern-set"
                    }
                }
            ]
        }`,
			expectedOutput: "my-regex-pattern-set",
		},
		{
			name:         "aws_waf_rule",
			resourceType: "aws_waf_rule",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-rule"
                    }
                }
            ]
        }`,
			expectedOutput: "my-rule",
		},
		{
			name:         "aws_waf_rule_group",
			resourceType: "aws_waf_rule_group",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-rule-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-rule-group",
		},
		{
			name:         "aws_waf_geo_match_set",
			resourceType: "aws_waf_geo_match_set",
			inputJSON: `{
            "name": "my-geo-match-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "geomatchset-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "geomatchset-12345678",
		},
		{
			name:         "aws_waf_ipset",
			resourceType: "aws_waf_ipset",
			inputJSON: `{
            "name": "my-ip-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "ipset-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "ipset-12345678",
		},
		{
			name:         "aws_waf_rate_based_rule",
			resourceType: "aws_waf_rate_based_rule",
			inputJSON: `{
            "name": "my-rate-based-rule",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "ratebasedrule-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "ratebasedrule-12345678",
		},
		{
			name:         "aws_waf_regex_match_set",
			resourceType: "aws_waf_regex_match_set",
			inputJSON: `{
            "name": "my-regex-match-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "regexmatchset-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "regexmatchset-12345678",
		},
		{
			name:         "aws_ec2_transit_gateway_vpc_attachment",
			resourceType: "aws_ec2_transit_gateway_vpc_attachment",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "tgw-attach-0123456789abcdef"
                    }
                }
            ]
        }`,
			expectedOutput: "tgw-attach-0123456789abcdef",
		},
		{
			name:         "aws_vpc",
			resourceType: "aws_vpc",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "vpc-0123456789abcdef"
                    }
                }
            ]
        }`,
			expectedOutput: "vpc-0123456789abcdef",
		},
		{
			name:         "aws_vpc_peering_connection",
			resourceType: "aws_vpc_peering_connection",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "pcx-0123456789abcdef"
                    }
                }
            ]
        }`,
			expectedOutput: "pcx-0123456789abcdef",
		},
		{
			name:         "aws_vpn_connection",
			resourceType: "aws_vpn_connection",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "vpn-0123456789abcdef"
                    }
                }
            ]
        }`,
			expectedOutput: "vpn-0123456789abcdef",
		},
		{
			name:         "aws_vpn_gateway",
			resourceType: "aws_vpn_gateway",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "vgw-0123456789abcdef"
                    }
                }
            ]
        }`,
			expectedOutput: "vgw-0123456789abcdef",
		},
		{
			name:         "aws_waf_byte_match_set",
			resourceType: "aws_waf_byte_match_set",
			inputJSON: `{
            "name": "my-byte-match-set",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "bytematchset-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "bytematchset-12345678",
		},
		{
			name:         "aws_ssm_parameter",
			resourceType: "aws_ssm_parameter",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-parameter"
                    }
                }
            ]
        }`,
			expectedOutput: "my-parameter",
		},
		{
			name:         "aws_subnet",
			resourceType: "aws_subnet",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "subnet-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "subnet-12345678",
		},
		{
			name:         "aws_swf_domain",
			resourceType: "aws_swf_domain",
			inputJSON: `{
            "name": "my-domain",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-domain"
                    }
                }
            ]
        }`,
			expectedOutput: "my-domain",
		},
		{
			name:         "aws_ec2_transit_gateway_route_table",
			resourceType: "aws_ec2_transit_gateway_route_table",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "tgw-rtb-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "tgw-rtb-12345678",
		},
		{
			name:         "aws_sfn_state_machine",
			resourceType: "aws_sfn_state_machine",
			inputJSON: `{
            "name": "my-state-machine",
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:states:us-west-2:012345678901:stateMachine:my-state-machine"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:states:us-west-2:012345678901:stateMachine:my-state-machine",
		},
		{
			name:         "aws_security_group",
			resourceType: "aws_security_group",
			inputJSON: `{
            "name": "my-security-group",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "sg-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "sg-12345678",
		},
		{
			name:         "aws_security_group_rule",
			resourceType: "aws_security_group_rule",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "security_group_id": "sg-12345678",
                        "id": "sgrule-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "sg-12345678/sgrule-12345678",
		},
		{
			name:         "aws_sns_topic",
			resourceType: "aws_sns_topic",
			inputJSON: `{
            "name": "my-topic",
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:sns:us-west-2:012345678901:my-topic"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:sns:us-west-2:012345678901:my-topic",
		},
		{
			name:         "aws_sns_topic_subscription",
			resourceType: "aws_sns_topic_subscription",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "topic_arn": "arn:aws:sns:us-west-2:012345678901:my-topic",
                        "subscription_arn": "arn:aws:sns:us-west-2:012345678901:my-topic:12345678-1234-1234-1234-123456789012"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:sns:us-west-2:012345678901:my-topic/arn:aws:sns:us-west-2:012345678901:my-topic:12345678-1234-1234-1234-123456789012",
		},
		{
			name:         "aws_sqs_queue",
			resourceType: "aws_sqs_queue",
			inputJSON: `{
            "name": "my-queue",
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:sqs:us-west-2:012345678901:my-queue"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:sqs:us-west-2:012345678901:my-queue",
		},
		{
			name:         "aws_ses_configuration_set",
			resourceType: "aws_ses_configuration_set",
			inputJSON: `{
            "name": "my-configuration-set",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-configuration-set"
                    }
                }
            ]
        }`,
			expectedOutput: "my-configuration-set",
		},
		{
			name:         "aws_ses_domain_identity",
			resourceType: "aws_ses_domain_identity",
			inputJSON: `{
            "name": "my-domain-identity",
            "instances": [
                {
                    "attributes_flat": {
                        "domain": "my-domain-identity"
                    }
                }
            ]
        }`,
			expectedOutput: "my-domain-identity",
		},
		{
			name:         "aws_ses_email_identity",
			resourceType: "aws_ses_email_identity",
			inputJSON: `{
            "name": "my-email-identity",
            "instances": [
                {
                    "attributes_flat": {
                        "email": "my-email-identity"
                    }
                }
            ]
        }`,
			expectedOutput: "my-email-identity",
		},
		{
			name:         "aws_ses_receipt_rule",
			resourceType: "aws_ses_receipt_rule",
			inputJSON: `{
            "rule_set_name": "my-rule-set",
            "name": "my-receipt-rule",
            "instances": [
                {
                    "attributes_flat": {
                        "rule_set_name": "my-rule-set",
                        "name": "my-receipt-rule"
                    }
                }
            ]
        }`,
			expectedOutput: "my-rule-set/my-receipt-rule",
		},
		{
			name:         "aws_ses_receipt_rule_set",
			resourceType: "aws_ses_receipt_rule_set",
			inputJSON: `{
            "name": "my-receipt-rule-set",
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:ses:us-west-2:123456789012:receipt-rule-set/my-receipt-rule-set"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:ses:us-west-2:123456789012:receipt-rule-set/my-receipt-rule-set",
		},
		{
			name:         "aws_ses_template",
			resourceType: "aws_ses_template",
			inputJSON: `{
            "name": "my-template",
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:ses:us-west-2:123456789012:template/my-template"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:ses:us-west-2:123456789012:template/my-template",
		},
		{
			name:         "aws_sfn_activity",
			resourceType: "aws_sfn_activity",
			inputJSON: `{
            "name": "my-activity",
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:states:us-west-2:123456789012:activity:my-activity"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:states:us-west-2:123456789012:activity:my-activity",
		},
		{
			name:         "aws_ses_configuration_set",
			resourceType: "aws_ses_configuration_set",
			inputJSON: `{
            "name": "my-configuration-set",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-configuration-set"
                    }
                }
            ]
        }`,
			expectedOutput: "my-configuration-set",
		},
		{
			name:         "aws_ses_domain_identity",
			resourceType: "aws_ses_domain_identity",
			inputJSON: `{
            "domain": "example.com",
            "instances": [
                {
                    "attributes_flat": {
                        "domain": "example.com"
                    }
                }
            ]
        }`,
			expectedOutput: "example.com",
		},
		{
			name:         "aws_ses_email_identity",
			resourceType: "aws_ses_email_identity",
			inputJSON: `{
            "email": "user@example.com",
            "instances": [
                {
                    "attributes_flat": {
                        "email": "user@example.com"
                    }
                }
            ]
        }`,
			expectedOutput: "user@example.com",
		},
		{
			name:         "aws_ses_receipt_rule",
			resourceType: "aws_ses_receipt_rule",
			inputJSON: `{
            "rule_set_name": "my-rule-set",
            "name": "my-rule",
            "instances": [
                {
                    "attributes_flat": {
                        "rule_set_name": "my-rule-set",
                        "name": "my-rule"
                    }
                }
            ]
        }`,
			expectedOutput: "my-rule-set/my-rule",
		},
		{
			name:         "aws_ses_receipt_rule_set",
			resourceType: "aws_ses_receipt_rule_set",
			inputJSON: `{
            "arn": "arn:aws:ses:us-west-2:123456789012:rule-set/my-rule-set",
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:ses:us-west-2:123456789012:rule-set/my-rule-set"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:ses:us-west-2:123456789012:rule-set/my-rule-set",
		},
		{
			name:         "aws_ses_template",
			resourceType: "aws_ses_template",
			inputJSON: `{
            "arn": "arn:aws:ses:us-west-2:123456789012:template/my-template",
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:ses:us-west-2:123456789012:template/my-template"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:ses:us-west-2:123456789012:template/my-template",
		},
		{
			name:         "aws_route53_record",
			resourceType: "aws_route53_record",
			inputJSON: `{
            "zone_id": "my-zone-id",
            "name": "my-record",
            "instances": [
                {
                    "attributes_flat": {
                        "zone_id": "my-zone-id",
                        "name": "my-record"
                    }
                }
            ]
        }`,
			expectedOutput: "my-zone-id/my-record",
		},
		{
			name:         "aws_route_table",
			resourceType: "aws_route_table",
			inputJSON: `{
            "name": "my-route-table",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-route-table"
                    }
                }
            ]
        }`,
			expectedOutput: "my-route-table",
		},
		{
			name:         "aws_main_route_table_association",
			resourceType: "aws_main_route_table_association",
			inputJSON: `{
            "name": "my-main-route-table-association",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-main-route-table-association"
                    }
                }
            ]
        }`,
			expectedOutput: "my-main-route-table-association",
		},
		{
			name:         "aws_route_table_association",
			resourceType: "aws_route_table_association",
			inputJSON: `{
            "name": "my-route-table-association",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-route-table-association"
                    }
                }
            ]
        }`,
			expectedOutput: "my-route-table-association",
		},
		{
			name:         "aws_accessanalyzer_analyzer",
			resourceType: "aws_accessanalyzer_analyzer",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:accessanalyzer:us-east-1:123456789012:analyzer/my-analyzer"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:accessanalyzer:us-east-1:123456789012:analyzer/my-analyzer",
		},
		{
			name:         "aws_acm_certificate",
			resourceType: "aws_acm_certificate",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012",
		},
		{
			name:         "aws_lb",
			resourceType: "aws_lb",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/12345678",
		},
		{
			name:         "aws_lb_listener",
			resourceType: "aws_lb_listener",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "load_balancer_arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/12345678",
                        "port": "80"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/12345678/80",
		},
		{
			name:         "aws_lb_listener_rule",
			resourceType: "aws_lb_listener_rule",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "load_balancer_arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/12345678",
                        "listener_arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:listener/app/my-load-balancer/12345678/abcdefgh12345678",
                        "arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:listener-rule/app/my-load-balancer/12345678/abcdefgh12345678/12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/12345678/arn:aws:elasticloadbalancing:us-west-2:123456789012:listener/app/my-load-balancer/12345678/abcdefgh12345678/arn:aws:elasticloadbalancing:us-west-2:123456789012:listener-rule/app/my-load-balancer/12345678/abcdefgh12345678/12345678",
		},
		{
			name:         "aws_lb_listener_certificate",
			resourceType: "aws_lb_listener_certificate",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "load_balancer_arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/12345678",
                        "listener_arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:listener/app/my-load-balancer/12345678/abcdefgh12345678",
                        "certificate_arn": "arn:aws:acm:us-west-2:123456789012:certificate/12345678-1234-1234-1234-123456789012"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/12345678/arn:aws:elasticloadbalancing:us-west-2:123456789012:listener/app/my-load-balancer/12345678/abcdefgh12345678/arn:aws:acm:us-west-2:123456789012:certificate/12345678-1234-1234-1234-123456789012",
		},
		{
			name:         "aws_lb_target_group",
			resourceType: "aws_lb_target_group",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-target-group/12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-target-group/12345678",
		},
		{
			name:         "aws_lb_target_group_attachment",
			resourceType: "aws_lb_target_group_attachment",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "target_group_arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-target-group/12345678",
                        "target_id": "i-0123456789abcdef0"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-target-group/12345678/i-0123456789abcdef0",
		},
		{
			name:         "aws_api_gateway_authorizer",
			resourceType: "aws_api_gateway_authorizer",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "rest_api_id": "1234567890",
                        "id": "abcdefghij"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890/abcdefghij",
		},
		{
			name:         "aws_api_gateway_api_key",
			resourceType: "aws_api_gateway_api_key",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "12345678",
		},
		{
			name:         "aws_api_gateway_documentation_part",
			resourceType: "aws_api_gateway_documentation_part",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "rest_api_id": "1234567890",
                        "id": "abcdefghij"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890/abcdefghij",
		},
		{
			name:         "aws_api_gateway_gateway_response",
			resourceType: "aws_api_gateway_gateway_response",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "rest_api_id": "1234567890",
                        "response_type": "DEFAULT_4XX"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890/DEFAULT_4XX",
		},
		{
			name:         "aws_api_gateway_integration",
			resourceType: "aws_api_gateway_integration",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "rest_api_id": "1234567890",
                        "resource_id": "abcdefghij"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890/abcdefghij",
		},
		{
			name:         "aws_api_gateway_integration_response",
			resourceType: "aws_api_gateway_integration_response",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "rest_api_id": "1234567890",
                        "resource_id": "abcdefghij",
                        "http_method": "GET"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890/abcdefghij/GET",
		},
		{
			name:         "aws_api_gateway_method",
			resourceType: "aws_api_gateway_method",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "rest_api_id": "1234567890",
                        "resource_id": "abcdefghij",
                        "http_method": "GET"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890/abcdefghij/GET",
		},
		{
			name:         "aws_api_gateway_method_response",
			resourceType: "aws_api_gateway_method_response",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "rest_api_id": "1234567890",
                        "resource_id": "abcdefghij",
                        "http_method": "GET",
                        "status_code": "200"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890/abcdefghij/GET/200",
		},
		{
			name:         "aws_api_gateway_model",
			resourceType: "aws_api_gateway_model",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "rest_api_id": "1234567890",
                        "name": "my_model"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890/my_model",
		},
		{
			name:         "aws_api_gateway_resource",
			resourceType: "aws_api_gateway_resource",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "rest_api_id": "1234567890",
                        "path": "my-resource"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890/my-resource",
		},
		{
			name:         "aws_api_gateway_rest_api",
			resourceType: "aws_api_gateway_rest_api",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "1234567890"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890",
		},
		{
			name:         "aws_api_gateway_stage",
			resourceType: "aws_api_gateway_stage",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "rest_api_id": "1234567890",
                        "stage_name": "my-stage"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890/my-stage",
		},
		{
			name:         "aws_api_gateway_usage_plan",
			resourceType: "aws_api_gateway_usage_plan",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "1234567890"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890",
		},
		{
			name:         "aws_api_gateway_vpc_link",
			resourceType: "aws_api_gateway_vpc_link",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "1234567890"
                    }
                }
            ]
        }`,
			expectedOutput: "1234567890",
		},
		{
			name:         "aws_appsync_graphql_api",
			resourceType: "aws_appsync_graphql_api",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:appsync:us-west-2:012345678901:apis/abcdefghij"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:appsync:us-west-2:012345678901:apis/abcdefghij",
		},
		{
			name:         "aws_autoscaling_group",
			resourceType: "aws_autoscaling_group",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-asg"
                    }
                }
            ]
        }`,
			expectedOutput: "my-asg",
		},
		{
			name:         "aws_launch_configuration",
			resourceType: "aws_launch_configuration",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-launch-config"
                    }
                }
            ]
        }`,
			expectedOutput: "my-launch-config",
		},
		{
			name:         "aws_launch_template",
			resourceType: "aws_launch_template",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "lt-0123456789abcdef0"
                    }
                }
            ]
        }`,
			expectedOutput: "lt-0123456789abcdef0",
		},
		{
			name:         "aws_batch_compute_environment",
			resourceType: "aws_batch_compute_environment",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:batch:us-west-2:012345678901:compute-environment/my-ce"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:batch:us-west-2:012345678901:compute-environment/my-ce",
		},
		{
			name:         "aws_batch_job_definition",
			resourceType: "aws_batch_job_definition",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:batch:us-west-2:012345678901:job-definition/my-job-def:1"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:batch:us-west-2:012345678901:job-definition/my-job-def:1",
		},
		{
			name:         "aws_batch_job_queue",
			resourceType: "aws_batch_job_queue",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:batch:us-west-2:012345678901:job-queue/my-job-queue"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:batch:us-west-2:012345678901:job-queue/my-job-queue",
		},
		{
			name:         "aws_budgets_budget",
			resourceType: "aws_budgets_budget",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-budget"
                    }
                }
            ]
        }`,
			expectedOutput: "my-budget",
		},
		{
			name:         "aws_cloud9_environment_ec2",
			resourceType: "aws_cloud9_environment_ec2",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-environment"
                    }
                }
            ]
        }`,
			expectedOutput: "my-environment",
		},
		{
			name:         "aws_cloudformation_stack",
			resourceType: "aws_cloudformation_stack",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-stack-id"
                    }
                }
            ]
        }`,
			expectedOutput: "my-stack-id",
		},
		{
			name:         "aws_cloudformation_stack_set",
			resourceType: "aws_cloudformation_stack_set",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-stack-set-id"
                    }
                }
            ]
        }`,
			expectedOutput: "my-stack-set-id",
		},
		{
			name:         "aws_cloudformation_stack_set_instance",
			resourceType: "aws_cloudformation_stack_set_instance",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-stack-set-instance-id"
                    }
                }
            ]
        }`,
			expectedOutput: "my-stack-set-instance-id",
		},
		{
			name:         "aws_cloudfront_distribution",
			resourceType: "aws_cloudfront_distribution",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-cloudfront-distribution-id"
                    }
                }
            ]
        }`,
			expectedOutput: "my-cloudfront-distribution-id",
		},
		{
			name:         "aws_cloudfront_cache_policy",
			resourceType: "aws_cloudfront_cache_policy",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-cloudfront-cache-policy"
                    }
                }
            ]
        }`,
			expectedOutput: "my-cloudfront-cache-policy",
		},
		{
			name:         "aws_cloudhsm_v2_cluster",
			resourceType: "aws_cloudhsm_v2_cluster",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-cloudhsm-cluster"
                    }
                }
            ]
        }`,
			expectedOutput: "my-cloudhsm-cluster",
		},
		{
			name:         "aws_cloudhsm_v2_hsm",
			resourceType: "aws_cloudhsm_v2_hsm",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-cloudhsm-hsm"
                    }
                }
            ]
        }`,
			expectedOutput: "my-cloudhsm-hsm",
		},
		{
			name:         "aws_cloudtrail",
			resourceType: "aws_cloudtrail",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-cloudtrail"
                    }
                }
            ]
        }`,
			expectedOutput: "my-cloudtrail",
		},
		{
			name:         "aws_cloudwatch_dashboard",
			resourceType: "aws_cloudwatch_dashboard",
			inputJSON: `{
            "dashboard_name": "my-dashboard",
            "instances": [
                {
                    "attributes_flat": {
                        "dashboard_name": "arn:aws:cloudwatch:us-west-2:012345678901:dashboard/my-dashboard"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:cloudwatch:us-west-2:012345678901:dashboard/my-dashboard",
		},
		{
			name:         "aws_cloudwatch_event_rule",
			resourceType: "aws_cloudwatch_event_rule",
			inputJSON: `{
            "name": "my-rule",
            "instances": [
                {
                    "attributes_flat": {
						"name": "my-rule",
                        "arn": "arn:aws:events:us-west-2:012345678901:rule/my-rule"
                    }
                }
            ]
        }`,
			expectedOutput: "my-rule",
		},
		{
			name:         "aws_cloudwatch_event_target",
			resourceType: "aws_cloudwatch_event_target",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "rule": "my-rule",
                        "arn": "arn:aws:lambda:us-west-2:012345678901:function:my-function"
                    }
                }
            ]
        }`,
			expectedOutput: "my-rule",
		},
		{
			name:         "aws_cloudwatch_metric_alarm",
			resourceType: "aws_cloudwatch_metric_alarm",
			inputJSON: `{
            "alarm_name": "my-alarm",
            "instances": [
                {
                    "attributes_flat": {
						"alarm_name": "my-alarm",
                        "alarm_arn": "arn:aws:cloudwatch:us-west-2:012345678901:alarm:my-alarm"
                    }
                }
            ]
        }`,
			expectedOutput: "my-alarm",
		},
		{
			name:         "aws_codebuild_project",
			resourceType: "aws_codebuild_project",
			inputJSON: `{
            "name": "my-project",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-project"
                    }
                }
            ]
        }`,
			expectedOutput: "my-project",
		},
		{
			name:         "aws_codecommit_repository",
			resourceType: "aws_codecommit_repository",
			inputJSON: `{
            "name": "my-repo",
            "instances": [
                {
                    "attributes_flat": {
                        "repository_name": "my-repo"
                    }
                }
            ]
        }`,
			expectedOutput: "my-repo",
		},
		{
			name:         "aws_codedeploy_app",
			resourceType: "aws_codedeploy_app",
			inputJSON: `{
            "name": "my-app",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-app"
                    }
                }
            ]
        }`,
			expectedOutput: "my-app",
		},
		{
			name:         "aws_codepipeline",
			resourceType: "aws_codepipeline",
			inputJSON: `{
            "name": "my-pipeline",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-pipeline"
                    }
                }
            ]
        }`,
			expectedOutput: "my-pipeline",
		},
		{
			name:         "aws_codepipeline_webhook",
			resourceType: "aws_codepipeline_webhook",
			inputJSON: `{
            "name": "my-webhook",
            "instances": [
                {
                    "attributes_flat": {
						"name": "my-webhook",
                        "arn": "arn:aws:codepipeline:us-west-2:012345678901:webhook:my-webhook"
                    }
                }
            ]
        }`,
			expectedOutput: "my-webhook",
		},
		{
			name:         "aws_cognito_identity_pool",
			resourceType: "aws_cognito_identity_pool",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
						"identity_pool_name": "my-identity-pool",
                        "id": "us-west-2:01234567-89ab-cdef-0123-456789abcdef"
                    }
                }
            ]
        }`,
			expectedOutput: "my-identity-pool",
		},
		{
			name:         "aws_cognito_user_pool",
			resourceType: "aws_cognito_user_pool",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
						"name": "my-user-pool",
                        "id": "us-west-2_012345678"
                    }
                }
            ]
        }`,
			expectedOutput: "my-user-pool",
		},
		{
			name:         "aws_config_config_rule",
			resourceType: "aws_config_config_rule",
			inputJSON: `{
            "name": "my-config-rule",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-config-rule"
                    }
                }
            ]
        }`,
			expectedOutput: "my-config-rule",
		},
		{
			name:         "aws_config_configuration_recorder",
			resourceType: "aws_config_configuration_recorder",
			inputJSON: `{
            "name": "my-configuration-recorder",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-configuration-recorder"
                    }
                }
            ]
        }`,
			expectedOutput: "my-configuration-recorder",
		},
		{
			name:         "aws_config_delivery_channel",
			resourceType: "aws_config_delivery_channel",
			inputJSON: `{
            "name": "my-delivery-channel",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-delivery-channel"
                    }
                }
            ]
        }`,
			expectedOutput: "my-delivery-channel",
		},
		{
			name:         "aws_customer_gateway",
			resourceType: "aws_customer_gateway",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "cgw-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "cgw-12345678",
		},
		{
			name:         "aws_datapipeline_pipeline",
			resourceType: "aws_datapipeline_pipeline",
			inputJSON: `{
            "name": "my-pipeline",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "df-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "df-12345678",
		},
		{
			name:         "aws_devicefarm_project",
			resourceType: "aws_devicefarm_project",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
						"name": "my-project",
                        "arn": "arn:aws:devicefarm:us-west-2:012345678901:project:12345678-abcd-1234-abcd-123456789012"
                    }
                }
            ]
        }`,
			expectedOutput: "my-project",
		},
		{
			name:         "aws_docdb_cluster",
			resourceType: "aws_docdb_cluster",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "docdb-cluster-1234"
                    }
                }
            ]
        }`,
			expectedOutput: "docdb-cluster-1234",
		},
		{
			name:         "aws_docdb_cluster_instance",
			resourceType: "aws_docdb_cluster_instance",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "docdb-cluster-instance-1234"
                    }
                }
            ]
        }`,
			expectedOutput: "docdb-cluster-instance-1234",
		},
		{
			name:         "aws_docdb_cluster_parameter_group",
			resourceType: "aws_docdb_cluster_parameter_group",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-docdb-cluster-parameter-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-docdb-cluster-parameter-group",
		},
		{
			name:         "aws_docdb_subnet_group",
			resourceType: "aws_docdb_subnet_group",
			inputJSON: `{
            "name": "my-docdb-subnet-group",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-docdb-subnet-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-docdb-subnet-group",
		},
		{
			name:         "aws_dynamodb_table",
			resourceType: "aws_dynamodb_table",
			inputJSON: `{
            "name": "my-dynamodb-table",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-dynamodb-table"
                    }
                }
            ]
        }`,
			expectedOutput: "my-dynamodb-table",
		},
		{
			name:         "aws_ebs_volume",
			resourceType: "aws_ebs_volume",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "vol-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "vol-12345678",
		},
		{
			name:         "aws_volume_attachment",
			resourceType: "aws_volume_attachment",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "att-0123456789abcdef"
                    }
                }
            ]
        }`,
			expectedOutput: "att-0123456789abcdef",
		},
		{
			name:         "aws_instance",
			resourceType: "aws_instance",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "i-0123456789abcdef"
                    }
                }
            ]
        }`,
			expectedOutput: "i-0123456789abcdef",
		},
		{
			name:         "aws_ecr_lifecycle_policy",
			resourceType: "aws_ecr_lifecycle_policy",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
						"registry_id": "012345678901",
            			"repository": "my-repo",
                        "policy_id": "my-policy-id"
                    }
                }
            ]
        }`,
			expectedOutput: "my-repo",
		},
		{
			name:         "aws_ecr_repository",
			resourceType: "aws_ecr_repository",
			inputJSON: `{
            "name": "my-repo",
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:ecr:us-west-2:012345678901:repository/my-repo",
                        "name": "my-repo"
                    }
                }
            ]
        }`,
			expectedOutput: "my-repo",
		},
		{
			name:         "aws_ecr_repository_policy",
			resourceType: "aws_ecr_repository_policy",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "registry_id": "012345678901",
                        "repository": "my-repo"
                    }
                }
            ]
        }`,
			expectedOutput: "my-repo",
		},
		{
			name:         "aws_ecrpublic_repository",
			resourceType: "aws_ecrpublic_repository",
			inputJSON: `{
            "name": "my-public-repo",
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:ecr-public:us-west-2:012345678901:repository/my-public-repo",
                        "name": "my-public-repo"
                    }
                }
            ]
        }`,
			expectedOutput: "my-public-repo",
		},
		{
			name:         "aws_ecs_cluster",
			resourceType: "aws_ecs_cluster",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-cluster"
                    }
                }
            ]
        }`,
			expectedOutput: "my-cluster",
		},
		{
			name:         "aws_ecs_service",
			resourceType: "aws_ecs_service",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-service"
                    }
                }
            ]
        }`,
			expectedOutput: "my-service",
		},
		{
			name:         "aws_ecs_task_definition",
			resourceType: "aws_ecs_task_definition",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "aws::arn"
                    }
                }
            ]
        }`,
			expectedOutput: "aws::arn",
		},
		{
			name:         "aws_efs_access_point",
			resourceType: "aws_efs_access_point",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "fsap-0123456789abcdef0"
                    }
                }
            ]
        }`,
			expectedOutput: "fsap-0123456789abcdef0",
		},
		{
			name:         "aws_efs_file_system",
			resourceType: "aws_efs_file_system",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "fs-01234567"
                    }
                }
            ]
        }`,
			expectedOutput: "fs-01234567",
		},
		{
			name:         "aws_efs_file_system_policy",
			resourceType: "aws_efs_file_system_policy",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "file_system_id": "fs-01234567"
                    }
                }
            ]
        }`,
			expectedOutput: "fs-01234567",
		},
		{
			name:         "aws_efs_access_point",
			resourceType: "aws_efs_access_point",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "fsap-0123456789abcdef0"
                    }
                }
            ]
        }`,
			expectedOutput: "fsap-0123456789abcdef0",
		},
		{
			name:         "aws_efs_file_system",
			resourceType: "aws_efs_file_system",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "fs-01234567"
                    }
                }
            ]
        }`,
			expectedOutput: "fs-01234567",
		},
		{
			name:         "aws_efs_file_system_policy",
			resourceType: "aws_efs_file_system_policy",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "file_system_id": "fs-01234567"
                    }
                }
            ]
        }`,
			expectedOutput: "fs-01234567",
		},
		{
			name:         "aws_eks_node_group",
			resourceType: "aws_eks_node_group",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "cluster_name": "my-cluster",
                        "node_group_name": "my-node-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-cluster:my-node-group",
		},
		{
			name:         "aws_elasticache_cluster",
			resourceType: "aws_elasticache_cluster",
			inputJSON: `{
            "name": "my-elasticache-cluster",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-elasticache-cluster"
                    }
                }
            ]
        }`,
			expectedOutput: "my-elasticache-cluster",
		},
		{
			name:         "aws_elasticache_parameter_group",
			resourceType: "aws_elasticache_parameter_group",
			inputJSON: `{
            "name": "my-parameter-group",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-parameter-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-parameter-group",
		},
		{
			name:         "aws_elasticache_subnet_group",
			resourceType: "aws_elasticache_subnet_group",
			inputJSON: `{
            "name": "my-subnet-group",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-subnet-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-subnet-group",
		},
		{
			name:         "aws_elasticache_replication_group",
			resourceType: "aws_elasticache_replication_group",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-replication-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-replication-group",
		},
		{
			name:         "aws_elastic_beanstalk_application",
			resourceType: "aws_elastic_beanstalk_application",
			inputJSON: `{
            "name": "my-application",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-application"
                    }
                }
            ]
        }`,
			expectedOutput: "my-application",
		},
		{
			name:         "aws_elastic_beanstalk_environment",
			resourceType: "aws_elastic_beanstalk_environment",
			inputJSON: `{
            "name": "my-environment",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-environment"
                    }
                }
            ]
        }`,
			expectedOutput: "my-environment",
		},
		{
			name:         "aws_elb",
			resourceType: "aws_elb",
			inputJSON: `{
            "name": "my-load-balancer",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-load-balancer"
                    }
                }
            ]
        }`,
			expectedOutput: "my-load-balancer",
		},
		{
			name:         "aws_emr_cluster",
			resourceType: "aws_emr_cluster",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "j-123456789012"
                    }
                }
            ]
        }`,
			expectedOutput: "j-123456789012",
		},
		{
			name:         "aws_emr_security_configuration",
			resourceType: "aws_emr_security_configuration",
			inputJSON: `{
            "name": "my-security-configuration",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-security-configuration"
                    }
                }
            ]
        }`,
			expectedOutput: "my-security-configuration",
		},
		{
			name:         "aws_network_interface",
			resourceType: "aws_network_interface",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "eni-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "eni-12345678",
		},
		{
			name:         "aws_elasticsearch_domain",
			resourceType: "aws_elasticsearch_domain",
			inputJSON: `{
            "domain_name": "my-elasticsearch-domain",
            "instances": [
                {
                    "attributes_flat": {
                        "domain_name": "my-elasticsearch-domain"
                    }
                }
            ]
        }`,
			expectedOutput: "my-elasticsearch-domain",
		},
		{
			name:         "aws_kinesis_firehose_delivery_stream",
			resourceType: "aws_kinesis_firehose_delivery_stream",
			inputJSON: `{
            "name": "my-firehose-stream",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-firehose-stream"
                    }
                }
            ]
        }`,
			expectedOutput: "my-firehose-stream",
		},
		{
			name:         "aws_glue_crawler",
			resourceType: "aws_glue_crawler",
			inputJSON: `{
            "name": "my-crawler",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-crawler"
                    }
                }
            ]
        }`,
			expectedOutput: "my-crawler",
		},
		{
			name:         "aws_glue_catalog_database",
			resourceType: "aws_glue_catalog_database",
			inputJSON: `{
            "name": "my-database",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-database"
                    }
                }
            ]
        }`,
			expectedOutput: "my-database",
		},
		{
			name:         "aws_glue_catalog_table",
			resourceType: "aws_glue_catalog_table",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
						"id": "my-id",
                        "name": "my-table",
                        "database_name": "my-database"
                    }
                }
            ]
        }`,
			expectedOutput: "my-id",
		},
		{
			name:         "aws_glue_job",
			resourceType: "aws_glue_job",
			inputJSON: `{
            "name": "my-job",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-job"
                    }
                }
            ]
        }`,
			expectedOutput: "my-job",
		},
		{
			name:         "aws_glue_trigger",
			resourceType: "aws_glue_trigger",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-trigger"
                    }
                }
            ]
        }`,
			expectedOutput: "my-trigger",
		},
		{
			name:         "aws_iam_access_key",
			resourceType: "aws_iam_access_key",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "user": "my-user",
                        "id": "MyAccessKey"
                    }
                }
            ]
        }`,
			expectedOutput: "MyAccessKey",
		},
		{
			name:         "aws_iam_group",
			resourceType: "aws_iam_group",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-group",
		},
		{
			name:         "aws_iam_group_policy",
			resourceType: "aws_iam_group_policy",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "group": "my-group",
                        "name": "my-policy"
                    }
                }
            ]
        }`,
			expectedOutput: "my-group:my-policy",
		},
		{
			name:         "aws_iam_group_policy_attachment",
			resourceType: "aws_iam_group_policy_attachment",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "group": "my-group",
                        "policy_arn": "arn:aws:iam::012345678901:policy/my-policy"
                    }
                }
            ]
        }`,
			expectedOutput: "my-group/arn:aws:iam::012345678901:policy/my-policy",
		},
		{
			name:         "aws_iam_instance_profile",
			resourceType: "aws_iam_instance_profile",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-instance-profile"
                    }
                }
            ]
        }`,
			expectedOutput: "my-instance-profile",
		},
		{
			name:         "aws_iam_policy",
			resourceType: "aws_iam_policy",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:iam::012345678901:policy/my-policy"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:iam::012345678901:policy/my-policy",
		},
		{
			name:         "aws_iam_role",
			resourceType: "aws_iam_role",
			inputJSON: `{
            "name": "my-role",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-role"
                    }
                }
            ]
        }`,
			expectedOutput: "my-role",
		},
		{
			name:         "aws_iam_role_policy",
			resourceType: "aws_iam_role_policy",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "role": "my-role",
                        "name": "my-policy"
                    }
                }
            ]
        }`,
			expectedOutput: "my-role/my-policy",
		},
		{
			name:         "aws_iam_role_policy_attachment",
			resourceType: "aws_iam_role_policy_attachment",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "role": "my-role",
                        "policy_arn": "arn:aws:iam::012345678901:policy/my-policy"
                    }
                }
            ]
        }`,
			expectedOutput: "my-role/arn:aws:iam::012345678901:policy/my-policy",
		},
		{
			name:         "aws_iam_user",
			resourceType: "aws_iam_user",
			inputJSON: `{
            "name": "my-user",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-user"
                    }
                }
            ]
        }`,
			expectedOutput: "my-user",
		},
		{
			name:         "aws_iam_user_group_membership",
			resourceType: "aws_iam_user_group_membership",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "user": "my-user",
                        "group": "my-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-user/my-group",
		},
		{
			name:         "aws_iam_user_policy",
			resourceType: "aws_iam_user_policy",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "user": "my-user",
                        "name": "my-policy"
                    }
                }
            ]
        }`,
			expectedOutput: "my-user:my-policy",
		},
		{
			name:         "aws_iam_user_policy_attachment",
			resourceType: "aws_iam_user_policy_attachment",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "user": "my-user",
                        "policy_arn": "arn:aws:iam::012345678901:policy/my-policy"
                    }
                }
            ]
        }`,
			expectedOutput: "my-user/arn:aws:iam::012345678901:policy/my-policy",
		},
		{
			name:         "aws_internet_gateway",
			resourceType: "aws_internet_gateway",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "igw-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "igw-12345678",
		},
		{
			name:         "aws_iot_thing",
			resourceType: "aws_iot_thing",
			inputJSON: `{
            "name": "my-thing",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-thing"
                    }
                }
            ]
        }`,
			expectedOutput: "my-thing",
		},
		{
			name:         "aws_iot_thing_type",
			resourceType: "aws_iot_thing_type",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-thing-type"
                    }
                }
            ]
        }`,
			expectedOutput: "my-thing-type",
		},
		{
			name:         "aws_iot_topic_rule",
			resourceType: "aws_iot_topic_rule",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-topic-rule"
                    }
                }
            ]
        }`,
			expectedOutput: "my-topic-rule",
		},
		{
			name:         "aws_iot_role_alias",
			resourceType: "aws_iot_role_alias",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-role-alias"
                    }
                }
            ]
        }`,
			expectedOutput: "my-role-alias",
		},
		{
			name:         "aws_kinesis_stream",
			resourceType: "aws_kinesis_stream",
			inputJSON: `{
            "name": "my-stream",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-stream"
                    }
                }
            ]
        }`,
			expectedOutput: "my-stream",
		},
		{
			name:         "aws_kms_key",
			resourceType: "aws_kms_key",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-key-id"
                    }
                }
            ]
        }`,
			expectedOutput: "my-key-id",
		},
		{
			name:         "aws_kms_alias",
			resourceType: "aws_kms_alias",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-alias"
                    }
                }
            ]
        }`,
			expectedOutput: "my-alias",
		},
		{
			name:         "aws_kms_grant",
			resourceType: "aws_kms_grant",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "key_id": "my-key-id",
                        "grant_id": "my-grant-id"
                    }
                }
            ]
        }`,
			expectedOutput: "my-key-id:my-grant-id",
		},
		{
			name:         "aws_lambda_event_source_mapping",
			resourceType: "aws_lambda_event_source_mapping",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "uuid": "my-uuid"
                    }
                }
            ]
        }`,
			expectedOutput: "my-uuid",
		},
		{
			name:         "aws_lambda_function",
			resourceType: "aws_lambda_function",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "function_name": "my-function-name"
                    }
                }
            ]
        }`,
			expectedOutput: "my-function-name",
		},
		{
			name:         "aws_lambda_function_event_invoke_config",
			resourceType: "aws_lambda_function_event_invoke_config",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "function_name": "my-function"
                    }
                }
            ]
        }`,
			expectedOutput: "my-function",
		},
		{
			name:         "aws_lambda_layer_version",
			resourceType: "aws_lambda_layer_version",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "layer_name": "my-layer",
                        "version": "1",
						"arn": "aws::lambda::arn"
                    }
                }
            ]
        }`,
			expectedOutput: "aws::lambda::arn",
		},
		{
			name:         "aws_lambda_permission",
			resourceType: "aws_lambda_permission",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "statement_id": "my-statement"
                    }
                }
            ]
        }`,
			expectedOutput: "my-statement",
		},
		{
			name:         "aws_cloudwatch_log_group",
			resourceType: "aws_cloudwatch_log_group",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
						"name": "/aws/lambda/my-function",
                        "arn": "arn:aws:logs:us-west-2:123456789012:log-group:/aws/lambda/my-function:*"
                    }
                }
            ]
        }`,
			expectedOutput: "/aws/lambda/my-function",
		},
		{
			name:         "aws_media_package_channel",
			resourceType: "aws_media_package_channel",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-channel"
                    }
                }
            ]
        }`,
			expectedOutput: "my-channel",
		},
		{
			name:         "aws_media_store_container",
			resourceType: "aws_media_store_container",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
						"name": "my-container",
                        "arn": "arn:aws:mediastore:us-west-2:123456789012:container/my-container"
                    }
                }
            ]
        }`,
			expectedOutput: "my-container",
		},
		{
			name:         "aws_medialive_channel",
			resourceType: "aws_medialive_channel",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-channel"
                    }
                }
            ]
        }`,
			expectedOutput: "my-channel",
		},
		{
			name:         "aws_medialive_input",
			resourceType: "aws_medialive_input",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-input"
                    }
                }
            ]
        }`,
			expectedOutput: "my-input",
		},
		{
			name:         "aws_medialive_input_security_group",
			resourceType: "aws_medialive_input_security_group",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-security-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-security-group",
		},
		{
			name:         "aws_msk_cluster",
			resourceType: "aws_msk_cluster",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "cluster_name": "my-msk-cluster"
                    }
                }
            ]
        }`,
			expectedOutput: "my-msk-cluster",
		},
		{
			name:         "aws_network_acl",
			resourceType: "aws_network_acl",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-network-acl"
                    }
                }
            ]
        }`,
			expectedOutput: "my-network-acl",
		},
		{
			name:         "aws_nat_gateway",
			resourceType: "aws_nat_gateway",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-nat-gateway"
                    }
                }
            ]
        }`,
			expectedOutput: "my-nat-gateway",
		},
		{
			name:         "aws_opsworks_application",
			resourceType: "aws_opsworks_application",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "app_id": "my-app"
                    }
                }
            ]
        }`,
			expectedOutput: "my-app",
		},
		{
			name:         "aws_opsworks_custom_layer",
			resourceType: "aws_opsworks_custom_layer",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "layer_id": "my-layer"
                    }
                }
            ]
        }`,
			expectedOutput: "my-layer",
		},
		{
			name:         "aws_opsworks_instance",
			resourceType: "aws_opsworks_instance",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "instance_id": "my-instance"
                    }
                }
            ]
        }`,
			expectedOutput: "my-instance",
		},
		{
			name:         "aws_opsworks_java_app_layer",
			resourceType: "aws_opsworks_java_app_layer",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "layer_id": "my-layer"
                    }
                }
            ]
        }`,
			expectedOutput: "my-layer",
		},
		{
			name:         "aws_opsworks_php_app_layer",
			resourceType: "aws_opsworks_php_app_layer",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "layer_id": "my-layer"
                    }
                }
            ]
        }`,
			expectedOutput: "my-layer",
		},
		{
			name:         "aws_opsworks_rds_db_instance",
			resourceType: "aws_opsworks_rds_db_instance",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "rds_db_instance_arn": "my-instance"
                    }
                }
            ]
        }`,
			expectedOutput: "my-instance",
		},
		{
			name:         "aws_opsworks_stack",
			resourceType: "aws_opsworks_stack",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "stack_id": "my-stack"
                    }
                }
            ]
        }`,
			expectedOutput: "my-stack",
		},
		{
			name:         "aws_opsworks_static_web_layer",
			resourceType: "aws_opsworks_static_web_layer",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "layer_id": "my-layer"
                    }
                }
            ]
        }`,
			expectedOutput: "my-layer",
		},
		{
			name:         "aws_opsworks_user_profile",
			resourceType: "aws_opsworks_user_profile",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "user_arn": "my-user"
                    }
                }
            ]
        }`,
			expectedOutput: "my-user",
		},
		{
			name:         "aws_organizations_account",
			resourceType: "aws_organizations_account",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-account"
                    }
                }
            ]
        }`,
			expectedOutput: "my-account",
		},
		{
			name:         "aws_organizations_organization",
			resourceType: "aws_organizations_organization",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "o-0123456789abcdefg"
                    }
                }
            ]
        }`,
			expectedOutput: "o-0123456789abcdefg",
		},
		{
			name:         "aws_organizations_organizational_unit",
			resourceType: "aws_organizations_organizational_unit",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "ou-0123456789abcdefg-0123456789abcdefg"
                    }
                }
            ]
        }`,
			expectedOutput: "ou-0123456789abcdefg-0123456789abcdefg",
		},
		{
			name:         "aws_organizations_policy",
			resourceType: "aws_organizations_policy",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "p-0123456789abcdefg"
                    }
                }
            ]
        }`,
			expectedOutput: "p-0123456789abcdefg",
		},
		{
			name:         "aws_organizations_policy_attachment",
			resourceType: "aws_organizations_policy_attachment",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "policy-attachment-id"
                    }
                }
            ]
        }`,
			expectedOutput: "policy-attachment-id",
		},
		{
			name:         "aws_qldb_ledger",
			resourceType: "aws_qldb_ledger",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
						"name": "my-ledger",
                        "arn": "arn:aws:qldb:us-east-1:123456789012:ledger/my-ledger"
                    }
                }
            ]
        }`,
			expectedOutput: "my-ledger",
		},
		{
			name:         "aws_db_instance",
			resourceType: "aws_db_instance",
			inputJSON: `{
            "name": "my-db-instance",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "db-instance-id"
                    }
                }
            ]
        }`,
			expectedOutput: "db-instance-id",
		},
		{
			name:         "aws_sfn_state_machine",
			resourceType: "aws_sfn_state_machine",
			inputJSON: `{
            "name": "my-state-machine",
            "instances": [
                {
                    "attributes_flat": {
                        "arn": "arn:aws:states:us-west-2:012345678901:stateMachine:my-state-machine"
                    }
                }
            ]
        }`,
			expectedOutput: "arn:aws:states:us-west-2:012345678901:stateMachine:my-state-machine",
		},
		{
			name:         "aws_security_group",
			resourceType: "aws_security_group",
			inputJSON: `{
            "name": "my-security-group",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "sg-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "sg-12345678",
		},
		{
			name:         "aws_security_group_rule",
			resourceType: "aws_security_group_rule",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "security_group_id": "sg-12345678",
                        "id": "sgrule-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "sg-12345678/sgrule-12345678",
		},
		{
			name:         "aws_db_parameter_group",
			resourceType: "aws_db_parameter_group",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-db-parameter-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-db-parameter-group",
		},
		{
			name:         "aws_db_snapshot",
			resourceType: "aws_db_snapshot",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-db-snapshot"
                    }
                }
            ]
        }`,
			expectedOutput: "my-db-snapshot",
		},
		{
			name:         "aws_db_subnet_group",
			resourceType: "aws_db_subnet_group",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-db-subnet-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-db-subnet-group",
		},
		{
			name:         "aws_db_option_group",
			resourceType: "aws_db_option_group",
			inputJSON: `{
            "name": "my-option-group",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-option-group"
                    }
                }
            ]
        }`,
			expectedOutput: "my-option-group",
		},
		{
			name:         "aws_db_event_subscription",
			resourceType: "aws_db_event_subscription",
			inputJSON: `{
            "name": "my-event-subscription",
            "instances": [
                {
                    "attributes_flat": {
                        "name": "my-event-subscription"
                    }
                }
            ]
        }`,
			expectedOutput: "my-event-subscription",
		},
		{
			name:         "aws_rds_global_cluster",
			resourceType: "aws_rds_global_cluster",
			inputJSON: `{
            "instances": [
                {
                    "attributes_flat": {
                        "id": "my-global-cluster"
                    }
                }
            ]
        }`,
			expectedOutput: "my-global-cluster",
		},
		{
			name:         "aws_resourcegroups_group",
			resourceType: "aws_resourcegroups_group",
			inputJSON: `{
            "name": "my-group",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "rg-12345678"
                    }
                }
            ]
        }`,
			expectedOutput: "rg-12345678",
		},
		{
			name:         "aws_route53_zone",
			resourceType: "aws_route53_zone",
			inputJSON: `{
            "name": "my-zone",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "/hostedzone/Z0123456789ABCDEF0123456"
                    }
                }
            ]
        }`,
			expectedOutput: "/hostedzone/Z0123456789ABCDEF0123456",
		},
		{
			name:         "aws_route_table",
			resourceType: "aws_route_table",
			inputJSON: `{
            "id": "rtb-0123456789abcdef0",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "rtb-0123456789abcdef0"
                    }
                }
            ]
        }`,
			expectedOutput: "rtb-0123456789abcdef0",
		},
		{
			name:         "aws_main_route_table_association",
			resourceType: "aws_main_route_table_association",
			inputJSON: `{
            "id": "rtbassoc-0123456789abcdef0",
            "instances": [
                {
                    "attributes_flat": {
                        "id": "rtbassoc-0123456789abcdef0"
                    }
                }
            ]
        }`,
			expectedOutput: "rtbassoc-0123456789abcdef0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider := terraformValueObjects.Provider("aws")
			resourceType := ResourceType(test.resourceType)
			resourcesParsed, err := gabs.ParseJSON([]byte(test.inputJSON))
			assert.Nil(t, err)

			output, err := GetRemoteCloudReference(resourcesParsed, provider, resourceType)
			assert.Equal(t, test.expectedOutput, output)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
