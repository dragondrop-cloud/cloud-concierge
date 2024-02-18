package markdowncreation

import (
	"os"
	"testing"
	"time"

	"github.com/atsushinee/go-markdown-generator/doc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdownCreator_initData(t *testing.T) {
	err := initializeTestFiles()
	require.NoError(t, err)

	defer os.RemoveAll("outputs")

	markdownCreator := NewMarkdownCreator()

	err = markdownCreator.initData()
	require.NoError(t, err)

	require.Equal(t, 3, len(markdownCreator.newResources))
	require.Equal(t, "terraform name of tfer ResourceID1 and type aws lb listener.", markdownCreator.newResources["aws_lb_listener.tfer--ResourceID1"])
	require.Equal(t, "terraform name of tfer ResourceID2 and type aws db subnet group.", markdownCreator.newResources["aws_db_subnet_group.tfer--ResourceID2"])
	require.Equal(t, "terraform name of tfer ResourceID3 and type aws s3 bucket.", markdownCreator.newResources["aws_s3_bucket.tfer--ResourceID3"])

	require.Equal(t, 3, len(markdownCreator.resourcesToCloudActions))
	require.Contains(t, markdownCreator.resourcesToCloudActions, "ResourceID1")
	require.Contains(t, markdownCreator.resourcesToCloudActions, "ResourceID2")
	require.Contains(t, markdownCreator.resourcesToCloudActions, "ResourceID3")
	require.Equal(t, "root", markdownCreator.resourcesToCloudActions["ResourceID1"]["modified"].Actor)
	require.Equal(t, "2023-08-09", markdownCreator.resourcesToCloudActions["ResourceID1"]["modified"].Timestamp)

	require.Equal(t, 3, len(markdownCreator.costEstimates))
	require.Equal(t, "16.425", markdownCreator.costEstimates[0].MonthlyCost)

	require.Equal(t, 3, len(markdownCreator.securityScan))
	require.Equal(t, "HIGH", markdownCreator.securityScan[0].Severity)

	require.Equal(t, 1, len(markdownCreator.managedDrift))
	require.Equal(t, "example_lb", markdownCreator.managedDrift[0].ResourceName)
	require.Equal(t, "enable_tls_version_and_cipher_suite_headers", markdownCreator.managedDrift[0].AttributeName)

	require.Equal(t, 3, len(markdownCreator.deletedResources))
	require.Equal(t, "aws_secretsmanager_secret", markdownCreator.deletedResources[0].ResourceType)
	require.Equal(t, "secret", markdownCreator.deletedResources[0].ResourceName)
}

func initializeTestFiles() error {
	err := os.Mkdir("outputs", 0o755)
	if err != nil {
		return err
	}

	newResources := `{
	  "aws_lb_listener.tfer--ResourceID1": "terraform name of tfer ResourceID1 and type aws lb listener.",
	  "aws_db_subnet_group.tfer--ResourceID2": "terraform name of tfer ResourceID2 and type aws db subnet group.",
	  "aws_s3_bucket.tfer--ResourceID3": "terraform name of tfer ResourceID3 and type aws s3 bucket."
	}`
	err = os.WriteFile("outputs/new-resources-to-documents.json", []byte(newResources), 0o600)
	if err != nil {
		return err
	}

	resourcesToCloudActions := `{
		"ResourceID1": {
			"modified": {
				"actor": "root",
				"timestamp": "2023-08-09"
			}
		},
		"ResourceID2": {
			"creation": {
				"actor": "root",
				"timestamp": "2023-08-09"
			}
		},
		"ResourceID3": {
			"creation": {
				"actor": "root",
				"timestamp": "2023-08-09"
			}
		}
	}`
	err = os.WriteFile("outputs/resources-to-cloud-actions.json", []byte(resourcesToCloudActions), 0o600)
	if err != nil {
		return err
	}

	costEstimates := `[
		{
			"cost_component": "Application load balancer",
			"is_usage_based": false,
			"monthly_cost": "16.425",
			"monthly_quantity": "730",
			"price": "0.0225",
			"resource_name": "aws_lb.tfer--resource-managed",
			"sub_resource_name": "",
			"unit": "hours"
		},
		{
			"cost_component": "Load balancer capacity units",
			"is_usage_based": false,
			"monthly_cost": "",
			"monthly_quantity": "",
			"price": "5.84",
			"resource_name": "aws_lb.tfer--tf-resource-managed",
			"sub_resource_name": "",
			"unit": "LCU"
		},
		{
			"cost_component": "S3 Bucket Storage",
			"is_usage_based": false,
			"monthly_cost": "",
			"monthly_quantity": "",
			"price": "10.02",
			"resource_name": "aws_s3_bucket.tfer--ResourceID3",
			"sub_resource_name": "",
			"unit": "LCU"
		}
	]`
	err = os.WriteFile("outputs/cost-estimates.json", []byte(costEstimates), 0o600)
	if err != nil {
		return err
	}

	securityScan := `{
	  "results": [
		{
		  "id": "ResourceID10",
		  "rule_id": "AVD-AWS-0053",
		  "long_id": "ResourceID10",
		  "rule_description": "Load balancer is exposed to the internet.",
		  "rule_provider": "aws",
		  "rule_service": "elb",
		  "impact": "The load balancer is exposed on the internet",
		  "resolution": "Switch to an internal load balancer or add a tfsec ignore",
		  "links": [
			"https://aquasecurity.github.io/tfsec/v1.28.1/checks/aws/elb/alb-not-public/",
			"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb"
		  ],
		  "description": "Load balancer is exposed publicly.",
		  "severity": "HIGH",
		  "warning": false,
		  "status": 0,
		  "resource": "aws_lb.tfer--ResourceID10",
		  "location": {
			"file_name": "",
			"start_line": 22,
			"end_line": 22
		  }
		},
		{
		  "id": "ResourceID11",
		  "rule_id": "AVD-AWS-0052",
		  "long_id": "ResourceID11",
		  "rule_description": "Load balancers should drop invalid headers",
		  "rule_provider": "aws",
		  "rule_service": "elb",
		  "impact": "Invalid headers being passed through to the target of the load balance may exploit vulnerabilities",
		  "resolution": "Set drop_invalid_header_fields to true",
		  "links": [
			"https://aquasecurity.github.io/tfsec/v1.28.1/checks/aws/elb/drop-invalid-headers/",
			"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb#drop_invalid_header_fields"
		  ],
		  "description": "Application load balancer is not set to drop invalid headers.",
		  "severity": "HIGH",
		  "warning": false,
		  "status": 0,
		  "resource": "aws_lb.tfer--ResourceID11",
		  "location": {
			"file_name": "",
			"start_line": 14,
			"end_line": 14
		  }
		},
		{
		  "id": "ResourceID12",
		  "rule_id": "AVD-AWS-0054",
		  "long_id": "ResourceID12",
		  "rule_description": "Use of plain HTTP.",
		  "rule_provider": "aws",
		  "rule_service": "elb",
		  "impact": "Your traffic is not protected",
		  "resolution": "Switch to HTTPS to benefit from TLS security features",
		  "links": [
			"https://aquasecurity.github.io/tfsec/v1.28.1/checks/aws/elb/http-not-used/",
			"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener"
		  ],
		  "description": "Listener for application load balancer does not use HTTPS.",
		  "severity": "CRITICAL",
		  "warning": false,
		  "status": 0,
		  "resource": "aws_lb_listener.tfer--ResourceID12",
		  "location": {
			"file_name": "",
			"start_line": 67,
			"end_line": 67
		  }
		}
	  ]
	}`
	err = os.WriteFile("outputs/security-scan.json", []byte(securityScan), 0o600)
	if err != nil {
		return err
	}

	resourcesDifferences := `[
	  {
		"RecentActor": "root",
		"RecentActionTimestamp": "2023-08-09",
		"AttributeName": "enable_tls_version_and_cipher_suite_headers",
		"TerraformValue": "false",
		"CloudValue": "true",
		"InstanceID": "i-0c0c0c0c0c0c0c0c0",
		"InstanceRegion": "us-east-1",
		"StateFileName": "aws-networking-dev",
		"ModuleName": "root",
		"ResourceType": "aws_lb",
		"ResourceName": "example_lb"
	  }
	]`
	err = os.WriteFile("outputs/drift-resources-differences.json", []byte(resourcesDifferences), 0o600)
	if err != nil {
		return err
	}

	resourcesDeleted := `[
	  {
		"InstanceID": "ResourceID20",
		"StateFileName": "state-file-1",
		"ModuleName": "module1",
		"ResourceType": "aws_secretsmanager_secret",
		"ResourceName": "secret"
	  },
	  {
		"InstanceID": "ResourceID21",
		"StateFileName": "state-file-2",
		"ModuleName": "module2",
		"ResourceType": "aws_s3_bucket",
		"ResourceName": "secret"
	  },
	  {
		"InstanceID": "ResourceID22",
		"StateFileName": "state-file-3",
		"ModuleName": "module3",
		"ResourceType": "aws_subnet",
		"ResourceName": "secret"
	  }
	]`
	err = os.WriteFile("outputs/drift-resources-deleted.json", []byte(resourcesDeleted), 0o600)
	if err != nil {
		return err
	}

	return nil
}

func TestMarkdownCreator_setGeneralData(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()

	// When
	markdownCreator.setGeneralData(report, "Job Name")

	// Then
	expectedMarkdown := "Job Name - State of Scanned Cloud Resources\n========================================================\n\n" +
		"# How to Read this Report\n\n'Job Name' has run. Of the resources the execution scans, at least one resource was identified to have drifted or be " +
		"outside of Terraform control. While code has been generated of the Terraform code and corresponding import statements needed to bring these " +
		"resources under Terraform control, below you will find a summary of the gaps identified in your current IaC posture.\n\n"
	assert.Equal(t, expectedMarkdown, report.String())
}

func TestMarkdownCreator_setFooter(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()

	// When
	markdownCreator.setFooter(report)

	// Then
	currentTime := time.Now().UTC()
	formattedTime := currentTime.Format("Created by Cloud Concierge at 15:04 UTC on 2006-01-02")

	expectedMarkdown := "#### Disclaimer\n\n*Indicates that a resource's cost is usage based. Since we currently do not infer/have knowledge of usage," +
		" costs may be material although indicated as 0 here.\n\nThis report presents information on the state of your cloud at a point in " +
		"time and as best Cloud Concierge is able to determine. Cloud Concierge does not currently scan every cloud resource for every cloud provider." +
		" For a list of supported resources, please see our [documentation](https://www.docs.dragondrop.cloud/).\n\n" +
		formattedTime
	assert.Equal(t, expectedMarkdown, report.String())
}
