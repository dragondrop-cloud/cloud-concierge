package identifycloudactors

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	queryParamData "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/identify_cloud_actors/query_param_data"
	driftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector/drift_detector"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func TestAWSLogQuerier_ExtractDataFromResourceResultManagedByTerraform(t *testing.T) {
	// Given
	alc := AWSLogQuerier{
		resourceToCloudTrailType: queryParamData.NewAWSResourceToCloudTrailLookup(),
	}
	timeVal, _ := time.Parse("2006-01-02", "2023-05-17")

	inputResourceResult := []*cloudtrail.Event{
		{
			EventId:     aws.String("929c77ea-c9f8-478f-a458-3c9d3644a78d"),
			EventName:   aws.String("CreateListener"),
			ReadOnly:    aws.String("false"),
			AccessKeyId: aws.String("ASJ"),
			EventTime:   &timeVal,
			EventSource: aws.String("elasticloadbalancing.amazonaws.com"),
			Username:    aws.String("test@dragondrop.cloud"),
			Resources: []*cloudtrail.Resource{
				{
					ResourceType: aws.String("AWS::ElasticLoadBalancingV2::LoadBalancer"),
					ResourceName: aws.String("arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302"),
				},
				{
					ResourceType: aws.String("AWS::ElasticLoadBalancingV2::Listener"),
					ResourceName: aws.String("arn:aws:elasticloadbalancing:us-east-1:682649898103:listener/app/tf-managed-demo-alb/4c89e21113613302/2e7801118647d076"),
				},
			},
			CloudTrailEvent: aws.String("{\"eventVersion\":\"1.08\",\"userIdentity\":{\"type\":\"AssumedRole\",\"principalId\":\"AROAZ54I2ZR3WBDHIUTRP:goodman.benjamin@dragondrop.cloud\",\"arn\":\"arn:aws:sts::682649898103:assumed-role/AWSReservedSSO_PowerUserAccess_ed3fdc44f35e9600/goodman.benjamin@dragondrop.cloud\",\"accountId\":\"682649898103\",\"accessKeyId\":\"ASIAZ54I2ZR37X23RUNJ\",\"sessionContext\":{\"sessionIssuer\":{\"type\":\"Role\",\"principalId\":\"AROAZ54I2ZR3WBDHIUTRP\",\"arn\":\"arn:aws:iam::682649898103:role/aws-reserved/sso.amazonaws.com/AWSReservedSSO_PowerUserAccess_ed3fdc44f35e9600\",\"accountId\":\"682649898103\",\"userName\":\"AWSReservedSSO_PowerUserAccess_ed3fdc44f35e9600\"},\"webIdFederationData\":{},\"attributes\":{\"creationDate\":\"2023-05-18T14:17:17Z\",\"mfaAuthenticated\":\"false\"}}},\"eventTime\":\"2023-05-18T14:19:00Z\",\"eventSource\":\"elasticloadbalancing.amazonaws.com\",\"eventName\":\"CreateListener\",\"awsRegion\":\"us-east-1\",\"sourceIPAddress\":\"146.70.41.20\",\"userAgent\":\"AWS Internal\",\"requestParameters\":{\"loadBalancerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302\",\"protocol\":\"HTTP\",\"port\":80,\"defaultActions\":[{\"order\":1,\"type\":\"fixed-response\",\"fixedResponseConfig\":{\"contentType\":\"text/plain\",\"statusCode\":\"503\"}}]},\"responseElements\":{\"listeners\":[{\"protocol\":\"HTTP\",\"defaultActions\":[{\"order\":1,\"type\":\"fixed-response\",\"fixedResponseConfig\":{\"contentType\":\"text/plain\",\"statusCode\":\"503\"}}],\"listenerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:listener/app/tf-managed-demo-alb/4c89e21113613302/2e7801118647d076\",\"loadBalancerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302\",\"port\":80}]},\"requestID\":\"d41fe5e8-2830-44c1-918d-8f37b76a11c3\",\"eventID\":\"929c77ea-c9f8-478f-a458-3c9d3644a78d\",\"readOnly\":false,\"eventType\":\"AwsApiCall\",\"apiVersion\":\"2015-12-01\",\"managementEvent\":true,\"recipientAccountId\":\"682649898103\",\"eventCategory\":\"Management\",\"sessionCredentialFromConsole\":\"true\"}"),
		},
		{
			EventId:     aws.String("6adee639-a5d2-418a-9d99-60124476e849"),
			EventName:   aws.String("ModifyLoadBalancerAttributes"),
			ReadOnly:    aws.String("false"),
			AccessKeyId: aws.String("ASG"),
			EventTime:   &timeVal,
			EventSource: aws.String("elasticloadbalancing.amazonaws.com"),
			Username:    aws.String("root"),
			Resources: []*cloudtrail.Resource{
				{
					ResourceType: aws.String("AWS::ElasticLoadBalancingV2::LoadBalancer"),
					ResourceName: aws.String("arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302"),
				},
			},
			CloudTrailEvent: aws.String("{\"eventVersion\":\"1.08\",\"userIdentity\":{\"type\":\"Root\",\"principalId\":\"682649898103\",\"arn\":\"arn:aws:iam::682649898103:root\",\"accountId\":\"682649898103\",\"accessKeyId\":\"ASIAZ54I2ZR3RUTSUDJG\",\"sessionContext\":{\"sessionIssuer\":{},\"webIdFederationData\":{},\"attributes\":{\"creationDate\":\"2023-05-17T13:38:07Z\",\"mfaAuthenticated\":\"true\"}}},\"eventTime\":\"2023-05-17T18:39:59Z\",\"eventSource\":\"elasticloadbalancing.amazonaws.com\",\"eventName\":\"ModifyLoadBalancerAttributes\",\"awsRegion\":\"us-east-1\",\"sourceIPAddress\":\"146.70.41.20\",\"userAgent\":\"AWS Internal\",\"requestParameters\":{\"attributes\":[{\"key\":\"access_logs.s3.enabled\",\"value\":\"false\"},{\"key\":\"idle_timeout.timeout_seconds\",\"value\":\"60\"},{\"key\":\"deletion_protection.enabled\",\"value\":\"true\"},{\"key\":\"routing.http2.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.drop_invalid_header_fields.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_client_port.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.preserve_host_header.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_header_processing.mode\",\"value\":\"append\"},{\"key\":\"load_balancing.cross_zone.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.desync_mitigation_mode\",\"value\":\"defensive\"},{\"key\":\"waf.fail_open.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.x_amzn_tls_version_and_cipher_suite.enabled\",\"value\":\"true\"}],\"loadBalancerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302\"},\"responseElements\":{\"attributes\":[{\"key\":\"access_logs.s3.enabled\",\"value\":\"false\"},{\"key\":\"idle_timeout.timeout_seconds\",\"value\":\"60\"},{\"key\":\"deletion_protection.enabled\",\"value\":\"true\"},{\"key\":\"routing.http2.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.drop_invalid_header_fields.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_client_port.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.preserve_host_header.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_header_processing.mode\",\"value\":\"append\"},{\"key\":\"load_balancing.cross_zone.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.desync_mitigation_mode\",\"value\":\"defensive\"},{\"key\":\"waf.fail_open.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.x_amzn_tls_version_and_cipher_suite.enabled\",\"value\":\"true\"},{\"key\":\"access_logs.s3.bucket\",\"value\":\"\"},{\"key\":\"access_logs.s3.prefix\",\"value\":\"\"}]},\"requestID\":\"a33b51c0-07bd-407e-91be-503584be5303\",\"eventID\":\"6adee639-a5d2-418a-9d99-60124476e849\",\"readOnly\":false,\"eventType\":\"AwsApiCall\",\"apiVersion\":\"2015-12-01\",\"managementEvent\":true,\"recipientAccountId\":\"682649898103\",\"eventCategory\":\"Management\",\"sessionCredentialFromConsole\":\"true\"}"),
		},
	}

	inputResourceType := "aws_lb"

	// When
	output, err := alc.ExtractDataFromResourceResult(inputResourceResult, inputResourceType, false)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expectedOutput := terraformValueObjects.ResourceActions{
		Creator: nil,
		Modifier: &terraformValueObjects.CloudActorTimeStamp{
			Actor:     terraformValueObjects.CloudActor("root"),
			Timestamp: terraformValueObjects.Timestamp("2023-05-17"),
		},
	}

	// Then
	if !reflect.DeepEqual(&output, &expectedOutput) {
		t.Errorf("expected:\n%v\ngot:\n%v", &expectedOutput, &output)
	}
}

func TestAWSLogQuerier_ExtractDataFromResourceResultNewToTerraform(t *testing.T) {
	// Given
	alc := AWSLogQuerier{
		resourceToCloudTrailType: queryParamData.NewAWSResourceToCloudTrailLookup(),
	}
	timeVal, _ := time.Parse("2006-01-02", "2023-05-18")

	inputResourceResult := []*cloudtrail.Event{
		{
			EventId:     aws.String("929c77ea-c9f8-478f-a458-3c9d3644a78d"),
			EventName:   aws.String("CreateListener"),
			ReadOnly:    aws.String("false"),
			AccessKeyId: aws.String("ASINJ"),
			EventTime:   &timeVal,
			EventSource: aws.String("elasticloadbalancing.amazonaws.com"),
			Username:    aws.String("test@dragondrop.cloud"),
			Resources: []*cloudtrail.Resource{
				{
					ResourceType: aws.String("AWS::ElasticLoadBalancingV2::LoadBalancer"),
					ResourceName: aws.String("arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302"),
				},
				{
					ResourceType: aws.String("AWS::ElasticLoadBalancingV2::Listener"),
					ResourceName: aws.String("arn:aws:elasticloadbalancing:us-east-1:682649898103:listener/app/tf-managed-demo-alb/4c89e21113613302/2e7801118647d076"),
				},
			},
			CloudTrailEvent: aws.String("{\"eventVersion\":\"1.08\",\"userIdentity\":{\"type\":\"AssumedRole\",\"principalId\":\"AROAZ54I2ZR3WBDHIUTRP:goodman.benjamin@dragondrop.cloud\",\"arn\":\"arn:aws:sts::682649898103:assumed-role/AWSReservedSSO_PowerUserAccess_ed3fdc44f35e9600/goodman.benjamin@dragondrop.cloud\",\"accountId\":\"682649898103\",\"accessKeyId\":\"ASIAZ54I2ZR37X23RUNJ\",\"sessionContext\":{\"sessionIssuer\":{\"type\":\"Role\",\"principalId\":\"AROAZ54I2ZR3WBDHIUTRP\",\"arn\":\"arn:aws:iam::682649898103:role/aws-reserved/sso.amazonaws.com/AWSReservedSSO_PowerUserAccess_ed3fdc44f35e9600\",\"accountId\":\"682649898103\",\"userName\":\"AWSReservedSSO_PowerUserAccess_ed3fdc44f35e9600\"},\"webIdFederationData\":{},\"attributes\":{\"creationDate\":\"2023-05-18T14:17:17Z\",\"mfaAuthenticated\":\"false\"}}},\"eventTime\":\"2023-05-18T14:19:00Z\",\"eventSource\":\"elasticloadbalancing.amazonaws.com\",\"eventName\":\"CreateListener\",\"awsRegion\":\"us-east-1\",\"sourceIPAddress\":\"146.70.41.20\",\"userAgent\":\"AWS Internal\",\"requestParameters\":{\"loadBalancerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302\",\"protocol\":\"HTTP\",\"port\":80,\"defaultActions\":[{\"order\":1,\"type\":\"fixed-response\",\"fixedResponseConfig\":{\"contentType\":\"text/plain\",\"statusCode\":\"503\"}}]},\"responseElements\":{\"listeners\":[{\"protocol\":\"HTTP\",\"defaultActions\":[{\"order\":1,\"type\":\"fixed-response\",\"fixedResponseConfig\":{\"contentType\":\"text/plain\",\"statusCode\":\"503\"}}],\"listenerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:listener/app/tf-managed-demo-alb/4c89e21113613302/2e7801118647d076\",\"loadBalancerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302\",\"port\":80}]},\"requestID\":\"d41fe5e8-2830-44c1-918d-8f37b76a11c3\",\"eventID\":\"929c77ea-c9f8-478f-a458-3c9d3644a78d\",\"readOnly\":false,\"eventType\":\"AwsApiCall\",\"apiVersion\":\"2015-12-01\",\"managementEvent\":true,\"recipientAccountId\":\"682649898103\",\"eventCategory\":\"Management\",\"sessionCredentialFromConsole\":\"true\"}"),
		},
		{
			EventId:     aws.String("6adee639-a5d2-418a-9d99-60124476e849"),
			EventName:   aws.String("ModifyLoadBalancerAttributes"),
			ReadOnly:    aws.String("false"),
			AccessKeyId: aws.String("ASIADJG"),
			EventTime:   &timeVal,
			EventSource: aws.String("elasticloadbalancing.amazonaws.com"),
			Username:    aws.String("root"),
			Resources: []*cloudtrail.Resource{
				{
					ResourceType: aws.String("AWS::ElasticLoadBalancingV2::LoadBalancer"),
					ResourceName: aws.String("arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302"),
				},
			},
			CloudTrailEvent: aws.String("{\"eventVersion\":\"1.08\",\"userIdentity\":{\"type\":\"Root\",\"principalId\":\"682649898103\",\"arn\":\"arn:aws:iam::682649898103:root\",\"accountId\":\"682649898103\",\"accessKeyId\":\"ASIAZ54I2ZR3RUTSUDJG\",\"sessionContext\":{\"sessionIssuer\":{},\"webIdFederationData\":{},\"attributes\":{\"creationDate\":\"2023-05-17T13:38:07Z\",\"mfaAuthenticated\":\"true\"}}},\"eventTime\":\"2023-05-17T18:39:59Z\",\"eventSource\":\"elasticloadbalancing.amazonaws.com\",\"eventName\":\"ModifyLoadBalancerAttributes\",\"awsRegion\":\"us-east-1\",\"sourceIPAddress\":\"146.70.41.20\",\"userAgent\":\"AWS Internal\",\"requestParameters\":{\"attributes\":[{\"key\":\"access_logs.s3.enabled\",\"value\":\"false\"},{\"key\":\"idle_timeout.timeout_seconds\",\"value\":\"60\"},{\"key\":\"deletion_protection.enabled\",\"value\":\"true\"},{\"key\":\"routing.http2.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.drop_invalid_header_fields.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_client_port.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.preserve_host_header.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_header_processing.mode\",\"value\":\"append\"},{\"key\":\"load_balancing.cross_zone.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.desync_mitigation_mode\",\"value\":\"defensive\"},{\"key\":\"waf.fail_open.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.x_amzn_tls_version_and_cipher_suite.enabled\",\"value\":\"true\"}],\"loadBalancerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302\"},\"responseElements\":{\"attributes\":[{\"key\":\"access_logs.s3.enabled\",\"value\":\"false\"},{\"key\":\"idle_timeout.timeout_seconds\",\"value\":\"60\"},{\"key\":\"deletion_protection.enabled\",\"value\":\"true\"},{\"key\":\"routing.http2.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.drop_invalid_header_fields.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_client_port.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.preserve_host_header.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_header_processing.mode\",\"value\":\"append\"},{\"key\":\"load_balancing.cross_zone.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.desync_mitigation_mode\",\"value\":\"defensive\"},{\"key\":\"waf.fail_open.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.x_amzn_tls_version_and_cipher_suite.enabled\",\"value\":\"true\"},{\"key\":\"access_logs.s3.bucket\",\"value\":\"\"},{\"key\":\"access_logs.s3.prefix\",\"value\":\"\"}]},\"requestID\":\"a33b51c0-07bd-407e-91be-503584be5303\",\"eventID\":\"6adee639-a5d2-418a-9d99-60124476e849\",\"readOnly\":false,\"eventType\":\"AwsApiCall\",\"apiVersion\":\"2015-12-01\",\"managementEvent\":true,\"recipientAccountId\":\"682649898103\",\"eventCategory\":\"Management\",\"sessionCredentialFromConsole\":\"true\"}"),
		},
	}

	inputResourceType := "aws_lb"

	// When
	output, err := alc.ExtractDataFromResourceResult(inputResourceResult, inputResourceType, true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expectedOutput := terraformValueObjects.ResourceActions{
		Creator: &terraformValueObjects.CloudActorTimeStamp{
			Actor:     terraformValueObjects.CloudActor("test@dragondrop.cloud"),
			Timestamp: terraformValueObjects.Timestamp("2023-05-18"),
		},
		Modifier: nil,
	}

	// Then
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected:\n%v\ngot:\n%v", expectedOutput, output)
	}
}

func TestAWSLogQuerier_ExtractDataFromResourceResultNewToTerraformWithModification(t *testing.T) {
	// Given
	alc := AWSLogQuerier{
		resourceToCloudTrailType: queryParamData.NewAWSResourceToCloudTrailLookup(),
	}
	timeVal, _ := time.Parse("2006-01-02", "2023-05-17")

	inputResourceResult := []*cloudtrail.Event{
		{
			EventId:     aws.String("6adee639-a5d2-418a-9d99-60124476e849"),
			EventName:   aws.String("ModifyLoadBalancerAttributes"),
			ReadOnly:    aws.String("false"),
			AccessKeyId: aws.String("ASIAZ54I2ZR3RUTSUDJG"),
			EventTime:   &timeVal,
			EventSource: aws.String("elasticloadbalancing.amazonaws.com"),
			Username:    aws.String("root"),
			Resources: []*cloudtrail.Resource{
				{
					ResourceType: aws.String("AWS::ElasticLoadBalancingV2::LoadBalancer"),
					ResourceName: aws.String("arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302"),
				},
			},
			CloudTrailEvent: aws.String("{\"eventVersion\":\"1.08\",\"userIdentity\":{\"type\":\"Root\",\"principalId\":\"682649898103\",\"arn\":\"arn:aws:iam::682649898103:root\",\"accountId\":\"682649898103\",\"accessKeyId\":\"ASIAZ54I2ZR3RUTSUDJG\",\"sessionContext\":{\"sessionIssuer\":{},\"webIdFederationData\":{},\"attributes\":{\"creationDate\":\"2023-05-17T13:38:07Z\",\"mfaAuthenticated\":\"true\"}}},\"eventTime\":\"2023-05-17T18:39:59Z\",\"eventSource\":\"elasticloadbalancing.amazonaws.com\",\"eventName\":\"ModifyLoadBalancerAttributes\",\"awsRegion\":\"us-east-1\",\"sourceIPAddress\":\"146.70.41.20\",\"userAgent\":\"AWS Internal\",\"requestParameters\":{\"attributes\":[{\"key\":\"access_logs.s3.enabled\",\"value\":\"false\"},{\"key\":\"idle_timeout.timeout_seconds\",\"value\":\"60\"},{\"key\":\"deletion_protection.enabled\",\"value\":\"true\"},{\"key\":\"routing.http2.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.drop_invalid_header_fields.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_client_port.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.preserve_host_header.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_header_processing.mode\",\"value\":\"append\"},{\"key\":\"load_balancing.cross_zone.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.desync_mitigation_mode\",\"value\":\"defensive\"},{\"key\":\"waf.fail_open.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.x_amzn_tls_version_and_cipher_suite.enabled\",\"value\":\"true\"}],\"loadBalancerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302\"},\"responseElements\":{\"attributes\":[{\"key\":\"access_logs.s3.enabled\",\"value\":\"false\"},{\"key\":\"idle_timeout.timeout_seconds\",\"value\":\"60\"},{\"key\":\"deletion_protection.enabled\",\"value\":\"true\"},{\"key\":\"routing.http2.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.drop_invalid_header_fields.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_client_port.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.preserve_host_header.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_header_processing.mode\",\"value\":\"append\"},{\"key\":\"load_balancing.cross_zone.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.desync_mitigation_mode\",\"value\":\"defensive\"},{\"key\":\"waf.fail_open.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.x_amzn_tls_version_and_cipher_suite.enabled\",\"value\":\"true\"},{\"key\":\"access_logs.s3.bucket\",\"value\":\"\"},{\"key\":\"access_logs.s3.prefix\",\"value\":\"\"}]},\"requestID\":\"a33b51c0-07bd-407e-91be-503584be5303\",\"eventID\":\"6adee639-a5d2-418a-9d99-60124476e849\",\"readOnly\":false,\"eventType\":\"AwsApiCall\",\"apiVersion\":\"2015-12-01\",\"managementEvent\":true,\"recipientAccountId\":\"682649898103\",\"eventCategory\":\"Management\",\"sessionCredentialFromConsole\":\"true\"}"),
		},
		{
			EventId:     aws.String("929c77ea-c9f8-478f-a458-3c9d3644a78d"),
			EventName:   aws.String("CreateListener"),
			ReadOnly:    aws.String("false"),
			AccessKeyId: aws.String("ASIAZ54I2ZR37X23RUNJ"),
			EventTime:   &timeVal,
			EventSource: aws.String("elasticloadbalancing.amazonaws.com"),
			Username:    aws.String("test@dragondrop.cloud"),
			Resources: []*cloudtrail.Resource{
				{
					ResourceType: aws.String("AWS::ElasticLoadBalancingV2::LoadBalancer"),
					ResourceName: aws.String("arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302"),
				},
				{
					ResourceType: aws.String("AWS::ElasticLoadBalancingV2::Listener"),
					ResourceName: aws.String("arn:aws:elasticloadbalancing:us-east-1:682649898103:listener/app/tf-managed-demo-alb/4c89e21113613302/2e7801118647d076"),
				},
			},
			CloudTrailEvent: aws.String("{\"eventVersion\":\"1.08\",\"userIdentity\":{\"type\":\"AssumedRole\",\"principalId\":\"AROAZ54I2ZR3WBDHIUTRP:goodman.benjamin@dragondrop.cloud\",\"arn\":\"arn:aws:sts::682649898103:assumed-role/AWSReservedSSO_PowerUserAccess_ed3fdc44f35e9600/goodman.benjamin@dragondrop.cloud\",\"accountId\":\"682649898103\",\"accessKeyId\":\"ASIAZ54I2ZR37X23RUNJ\",\"sessionContext\":{\"sessionIssuer\":{\"type\":\"Role\",\"principalId\":\"AROAZ54I2ZR3WBDHIUTRP\",\"arn\":\"arn:aws:iam::682649898103:role/aws-reserved/sso.amazonaws.com/AWSReservedSSO_PowerUserAccess_ed3fdc44f35e9600\",\"accountId\":\"682649898103\",\"userName\":\"AWSReservedSSO_PowerUserAccess_ed3fdc44f35e9600\"},\"webIdFederationData\":{},\"attributes\":{\"creationDate\":\"2023-05-18T14:17:17Z\",\"mfaAuthenticated\":\"false\"}}},\"eventTime\":\"2023-05-18T14:19:00Z\",\"eventSource\":\"elasticloadbalancing.amazonaws.com\",\"eventName\":\"CreateListener\",\"awsRegion\":\"us-east-1\",\"sourceIPAddress\":\"146.70.41.20\",\"userAgent\":\"AWS Internal\",\"requestParameters\":{\"loadBalancerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302\",\"protocol\":\"HTTP\",\"port\":80,\"defaultActions\":[{\"order\":1,\"type\":\"fixed-response\",\"fixedResponseConfig\":{\"contentType\":\"text/plain\",\"statusCode\":\"503\"}}]},\"responseElements\":{\"listeners\":[{\"protocol\":\"HTTP\",\"defaultActions\":[{\"order\":1,\"type\":\"fixed-response\",\"fixedResponseConfig\":{\"contentType\":\"text/plain\",\"statusCode\":\"503\"}}],\"listenerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:listener/app/tf-managed-demo-alb/4c89e21113613302/2e7801118647d076\",\"loadBalancerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302\",\"port\":80}]},\"requestID\":\"d41fe5e8-2830-44c1-918d-8f37b76a11c3\",\"eventID\":\"929c77ea-c9f8-478f-a458-3c9d3644a78d\",\"readOnly\":false,\"eventType\":\"AwsApiCall\",\"apiVersion\":\"2015-12-01\",\"managementEvent\":true,\"recipientAccountId\":\"682649898103\",\"eventCategory\":\"Management\",\"sessionCredentialFromConsole\":\"true\"}"),
		},
		{
			EventId:     aws.String("bcc8056a-6a43-4140-bf54-0348491e634c"),
			EventName:   aws.String("ModifyLoadBalancerAttributes"),
			ReadOnly:    aws.String("false"),
			AccessKeyId: aws.String("ASIAZ54I2ZR3RUTSUDJG"),
			EventTime:   &timeVal,
			EventSource: aws.String("elasticloadbalancing.amazonaws.com"),
			Username:    aws.String("root"),
			Resources: []*cloudtrail.Resource{
				{
					ResourceType: aws.String("AWS::ElasticLoadBalancingV2::LoadBalancer"),
					ResourceName: aws.String("arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302"),
				},
			},
			CloudTrailEvent: aws.String("{\"eventVersion\":\"1.08\",\"userIdentity\":{\"type\":\"Root\",\"principalId\":\"682649898103\",\"arn\":\"arn:aws:iam::682649898103:root\",\"accountId\":\"682649898103\",\"accessKeyId\":\"ASIAZ54I2ZR3RUTSUDJG\",\"sessionContext\":{\"sessionIssuer\":{},\"webIdFederationData\":{},\"attributes\":{\"creationDate\":\"2023-05-17T13:38:07Z\",\"mfaAuthenticated\":\"true\"}}},\"eventTime\":\"2023-05-17T13:40:53Z\",\"eventSource\":\"elasticloadbalancing.amazonaws.com\",\"eventName\":\"ModifyLoadBalancerAttributes\",\"awsRegion\":\"us-east-1\",\"sourceIPAddress\":\"146.70.41.20\",\"userAgent\":\"AWS Internal\",\"requestParameters\":{\"attributes\":[{\"key\":\"access_logs.s3.enabled\",\"value\":\"false\"},{\"key\":\"idle_timeout.timeout_seconds\",\"value\":\"60\"},{\"key\":\"deletion_protection.enabled\",\"value\":\"true\"},{\"key\":\"routing.http2.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.drop_invalid_header_fields.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_client_port.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.preserve_host_header.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_header_processing.mode\",\"value\":\"append\"},{\"key\":\"load_balancing.cross_zone.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.desync_mitigation_mode\",\"value\":\"defensive\"},{\"key\":\"waf.fail_open.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.x_amzn_tls_version_and_cipher_suite.enabled\",\"value\":\"false\"}],\"loadBalancerArn\":\"arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302\"},\"responseElements\":{\"attributes\":[{\"key\":\"access_logs.s3.enabled\",\"value\":\"false\"},{\"key\":\"idle_timeout.timeout_seconds\",\"value\":\"60\"},{\"key\":\"deletion_protection.enabled\",\"value\":\"true\"},{\"key\":\"routing.http2.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.drop_invalid_header_fields.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_client_port.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.preserve_host_header.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.xff_header_processing.mode\",\"value\":\"append\"},{\"key\":\"load_balancing.cross_zone.enabled\",\"value\":\"true\"},{\"key\":\"routing.http.desync_mitigation_mode\",\"value\":\"defensive\"},{\"key\":\"waf.fail_open.enabled\",\"value\":\"false\"},{\"key\":\"routing.http.x_amzn_tls_version_and_cipher_suite.enabled\",\"value\":\"false\"},{\"key\":\"access_logs.s3.bucket\",\"value\":\"\"},{\"key\":\"access_logs.s3.prefix\",\"value\":\"\"}]},\"requestID\":\"dc533430-01ef-4106-a7e6-80bf97775fa5\",\"eventID\":\"bcc8056a-6a43-4140-bf54-0348491e634c\",\"readOnly\":false,\"eventType\":\"AwsApiCall\",\"apiVersion\":\"2015-12-01\",\"managementEvent\":true,\"recipientAccountId\":\"682649898103\",\"eventCategory\":\"Management\",\"sessionCredentialFromConsole\":\"true\"}"),
		},
	}

	inputResourceType := "aws_lb"

	// When
	output, err := alc.ExtractDataFromResourceResult(inputResourceResult, inputResourceType, true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expectedOutput := terraformValueObjects.ResourceActions{
		Creator: &terraformValueObjects.CloudActorTimeStamp{
			Actor:     terraformValueObjects.CloudActor("test@dragondrop.cloud"),
			Timestamp: terraformValueObjects.Timestamp("2023-05-17"),
		},
		Modifier: &terraformValueObjects.CloudActorTimeStamp{
			Actor:     terraformValueObjects.CloudActor("root"),
			Timestamp: terraformValueObjects.Timestamp("2023-05-17"),
		},
	}

	// Then
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected:\n%v\ngot:\n%v", expectedOutput, output)
	}
}

func TestAWSLogQuerier_UpdateManagedDriftAttributeDifferences(t *testing.T) {
	// Given
	alc := AWSLogQuerier{
		managedDriftAttributeDifferences: []driftDetector.AttributeDifference{
			{
				InstanceID: "dragondrop",
				AttributeDetail: driftDetector.AttributeDetail{
					ResourceType:  "type-1",
					ResourceName:  "name-1",
					StateFileName: "state-1",
				},
			},
			{
				InstanceID: "dragondrop",
				AttributeDetail: driftDetector.AttributeDetail{
					ResourceType:  "type-1",
					ResourceName:  "name-1",
					StateFileName: "state-2",
				},
			},
		},
	}

	inputDivisionResourceActions := terraformValueObjects.ResourceActionMap{
		"state-1.type-1.name-1.dragondrop": {
			Modifier: &terraformValueObjects.CloudActorTimeStamp{
				Actor:     terraformValueObjects.CloudActor("root"),
				Timestamp: terraformValueObjects.Timestamp("2023-05-17"),
			},
		},
		"state-2.type-1.name-1.dragondrop": {
			Modifier: &terraformValueObjects.CloudActorTimeStamp{
				Actor:     terraformValueObjects.CloudActor("jenny-from-the-block"),
				Timestamp: terraformValueObjects.Timestamp("2023-05-16"),
			},
		},
	}

	expectedOutput := []driftDetector.AttributeDifference{
		{
			InstanceID:            "dragondrop",
			RecentActionTimestamp: "2023-05-17",
			RecentActor:           "root",
			AttributeDetail: driftDetector.AttributeDetail{
				ResourceType:  "type-1",
				ResourceName:  "name-1",
				StateFileName: "state-1",
			},
		},
		{
			InstanceID:            "dragondrop",
			RecentActionTimestamp: "2023-05-16",
			RecentActor:           "jenny-from-the-block",
			AttributeDetail: driftDetector.AttributeDetail{
				ResourceType:  "type-1",
				ResourceName:  "name-1",
				StateFileName: "state-2",
			},
		},
	}

	// when
	alc.UpdateManagedDriftAttributeDifferences(inputDivisionResourceActions)

	// then
	if !reflect.DeepEqual(alc.managedDriftAttributeDifferences, expectedOutput) {
		t.Errorf("expected:\n%v\ngot:\n%v", expectedOutput, alc.managedDriftAttributeDifferences)
	}
}
