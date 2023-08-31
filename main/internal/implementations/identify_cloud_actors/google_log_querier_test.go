package identifyCloudActors

//
//import (
//	"reflect"
//	"testing"
//
//	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
//)
//
//func TestGenerateLogFilter(t *testing.T) {
//	// Given
//	glc := GoogleLogQuerier{
//		division: terraformValueObjects.Division("test-div"),
//	}
//
//	// When
//	output := glc.generateLogFilter("my-id")
//
//	// Then
//	expectedOutput := "logName=projects/test-div/logs/cloudaudit.googleapis.com%2Factivity AND protoPayload.resourceName=my-id"
//
//	if output != expectedOutput {
//		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
//	}
//}
//
//func TestExtractDataFromResourceResult_ModifierThenCreator(t *testing.T) {
//	// Given
//	glc := GoogleLogQuerier{}
//	inputResourceResult := []byte(`{
//    "entries": [
//        {
//            "protoPayload": {
//                "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
//                "status": {},
//                "authenticationInfo": {
//                    "principalEmail": "goodman.benjamin@dragondrop.cloud"
//                },
//                "serviceName": "storage.googleapis.com",
//                "methodName": "storage.buckets.update",
//                "authorizationInfo": [
//                    {
//                        "resource": "projects/_/buckets/testing-out-this-bucket",
//                        "permission": "storage.buckets.update",
//                        "granted": true,
//                        "resourceAttributes": {}
//                    }
//                ],
//                "resourceName": "projects/_/buckets/testing-out-this-bucket"
//            },
//            "insertId": "yyo9azd9x9y",
//            "resource": {
//                "type": "gcs_bucket",
//                "labels": {
//                    "project_id": "dragondrop-dev",
//                    "bucket_name": "testing-out-this-bucket",
//                    "location": "us"
//                }
//            },
//            "timestamp": "2023-03-08T17:24:53.274706482Z",
//            "severity": "NOTICE",
//            "logName": "projects/dragondrop-dev/logs/cloudaudit.googleapis.com%2Factivity",
//            "receiveTimestamp": "2023-03-11T17:24:54.248853501Z"
//        },
//        {
//            "protoPayload": {
//                "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
//                "status": {},
//                "authenticationInfo": {
//                    "principalEmail": "goodman.benjamin@dragondrop.cloud"
//                },
//                "serviceName": "storage.googleapis.com",
//                "methodName": "storage.buckets.update",
//                "authorizationInfo": [
//                    {
//                        "resource": "projects/_/buckets/testing-out-this-bucket",
//                        "permission": "storage.buckets.update",
//                        "granted": true,
//                        "resourceAttributes": {}
//                    }
//                ],
//                "resourceName": "projects/_/buckets/testing-out-this-bucket"
//            },
//            "timestamp": "2023-03-08T17:24:02.675064542Z",
//            "severity": "NOTICE",
//            "logName": "projects/dragondrop-dev/logs/cloudaudit.googleapis.com%2Factivity",
//            "receiveTimestamp": "2023-03-08T17:24:03.418303345Z"
//        },
//        {
//            "protoPayload": {
//                "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
//                "status": {},
//                "authenticationInfo": {
//                    "principalEmail": "goodman.benjamin@dragondrop.cloud"
//                },
//                "serviceName": "storage.googleapis.com",
//                "methodName": "storage.buckets.create",
//                "authorizationInfo": [
//                    {
//                        "resource": "projects/_/buckets/testing-out-this-bucket",
//                        "permission": "storage.buckets.create",
//                        "granted": true,
//                        "resourceAttributes": {}
//                    }
//                ],
//                "resourceName": "projects/_/buckets/testing-out-this-bucket"
//            },
//            "insertId": "8c1kgad5b9c",
//            "resource": {
//                "type": "gcs_bucket",
//                "labels": {
//                    "bucket_name": "testing-out-this-bucket",
//                    "location": "us",
//                    "project_id": "dragondrop-dev"
//                }
//            },
//            "timestamp": "2023-02-25T20:31:16.322417060Z",
//            "severity": "NOTICE",
//            "logName": "projects/dragondrop-dev/logs/cloudaudit.googleapis.com%2Factivity",
//            "receiveTimestamp": "2023-02-25T20:31:17.243309785Z"
//        }
//    ]
//}`)
//	// When
//	output, err := glc.ExtractDataFromResourceResult(inputResourceResult, "", true)
//	if err != nil {
//		t.Errorf("unexpected error in test: %v", err)
//	}
//
//	// Then
//	expectedOutput := terraformValueObjects.ResourceActions{
//		Creator: terraformValueObjects.CloudActorTimeStamp{
//			Actor:     terraformValueObjects.CloudActor("goodman.benjamin@dragondrop.cloud"),
//			Timestamp: terraformValueObjects.Timestamp("2023-02-25"),
//		},
//		Modifier: terraformValueObjects.CloudActorTimeStamp{
//			Actor:     terraformValueObjects.CloudActor("goodman.benjamin@dragondrop.cloud"),
//			Timestamp: terraformValueObjects.Timestamp("2023-03-11"),
//		},
//	}
//
//	if !reflect.DeepEqual(output, expectedOutput) {
//		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
//	}
//}
//
//func TestExtractDataFromResourceResult_CreatorThenModifier(t *testing.T) {
//	glc := GoogleLogQuerier{}
//	inputResourceResult := []byte(`{
//    "entries": [
//        {
//            "protoPayload": {
//                "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
//                "status": {},
//                "authenticationInfo": {
//                    "principalEmail": "goodman.benjamin@dragondrop.cloud"
//                },
//                "serviceName": "storage.googleapis.com",
//                "methodName": "storage.buckets.create",
//                "authorizationInfo": [
//                    {
//                        "resource": "projects/_/buckets/testing-out-this-bucket",
//                        "permission": "storage.buckets.create",
//                        "granted": true,
//                        "resourceAttributes": {}
//                    }
//                ],
//                "resourceName": "projects/_/buckets/testing-out-this-bucket"
//            },
//            "insertId": "8c1kgad5b9c",
//            "resource": {
//                "type": "gcs_bucket",
//                "labels": {
//                    "bucket_name": "testing-out-this-bucket",
//                    "location": "us",
//                    "project_id": "dragondrop-dev"
//                }
//            },
//            "timestamp": "2023-02-25T20:31:16.322417060Z",
//            "severity": "NOTICE",
//            "logName": "projects/dragondrop-dev/logs/cloudaudit.googleapis.com%2Factivity",
//            "receiveTimestamp": "2023-02-25T20:31:17.243309785Z"
//        },
//        {
//            "protoPayload": {
//                "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
//                "status": {},
//                "authenticationInfo": {
//                    "principalEmail": "goodman.benjamin@dragondrop.cloud"
//                },
//                "serviceName": "storage.googleapis.com",
//                "methodName": "storage.buckets.update",
//                "authorizationInfo": [
//                    {
//                        "resource": "projects/_/buckets/testing-out-this-bucket",
//                        "permission": "storage.buckets.update",
//                        "granted": true,
//                        "resourceAttributes": {}
//                    }
//                ],
//                "resourceName": "projects/_/buckets/testing-out-this-bucket"
//            },
//            "insertId": "yyo9azd9x9y",
//            "resource": {
//                "type": "gcs_bucket",
//                "labels": {
//                    "project_id": "dragondrop-dev",
//                    "bucket_name": "testing-out-this-bucket",
//                    "location": "us"
//                }
//            },
//            "timestamp": "2023-03-08T17:24:53.274706482Z",
//            "severity": "NOTICE",
//            "logName": "projects/dragondrop-dev/logs/cloudaudit.googleapis.com%2Factivity",
//            "receiveTimestamp": "2023-03-11T17:24:54.248853501Z"
//        },
//        {
//            "protoPayload": {
//                "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
//                "status": {},
//                "authenticationInfo": {
//                    "principalEmail": "goodman.benjamin@dragondrop.cloud"
//                },
//                "serviceName": "storage.googleapis.com",
//                "methodName": "storage.buckets.update",
//                "authorizationInfo": [
//                    {
//                        "resource": "projects/_/buckets/testing-out-this-bucket",
//                        "permission": "storage.buckets.update",
//                        "granted": true,
//                        "resourceAttributes": {}
//                    }
//                ],
//                "resourceName": "projects/_/buckets/testing-out-this-bucket"
//            },
//            "timestamp": "2023-03-08T17:24:02.675064542Z",
//            "severity": "NOTICE",
//            "logName": "projects/dragondrop-dev/logs/cloudaudit.googleapis.com%2Factivity",
//            "receiveTimestamp": "2023-03-08T17:24:03.418303345Z"
//        }
//    ]
//}`)
//	// When
//	output, err := glc.ExtractDataFromResourceResult(inputResourceResult, "", true)
//	if err != nil {
//		t.Errorf("unexpected error in test: %v", err)
//	}
//
//	// Then
//	expectedOutput := terraformValueObjects.ResourceActions{
//		Creator: terraformValueObjects.CloudActorTimeStamp{
//			Actor:     terraformValueObjects.CloudActor("goodman.benjamin@dragondrop.cloud"),
//			Timestamp: terraformValueObjects.Timestamp("2023-02-25"),
//		},
//	}
//
//	if !reflect.DeepEqual(output, expectedOutput) {
//		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
//	}
//}
//
//func TestExtractDataFromResourceResult_OnlyModifier(t *testing.T) {
//	// Given
//	glc := GoogleLogQuerier{}
//	inputResourceResult := []byte(`{
//    "entries": [
//        {
//            "protoPayload": {
//                "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
//                "status": {},
//                "authenticationInfo": {
//                    "principalEmail": "goodman.benjamin@dragondrop.cloud"
//                },
//                "serviceName": "storage.googleapis.com",
//                "methodName": "storage.buckets.update",
//                "authorizationInfo": [
//                    {
//                        "resource": "projects/_/buckets/testing-out-this-bucket",
//                        "permission": "storage.buckets.update",
//                        "granted": true,
//                        "resourceAttributes": {}
//                    }
//                ],
//                "resourceName": "projects/_/buckets/testing-out-this-bucket"
//            },
//            "insertId": "yyo9azd9x9y",
//            "resource": {
//                "type": "gcs_bucket",
//                "labels": {
//                    "project_id": "dragondrop-dev",
//                    "bucket_name": "testing-out-this-bucket",
//                    "location": "us"
//                }
//            },
//            "timestamp": "2023-03-08T17:24:53.274706482Z",
//            "severity": "NOTICE",
//            "logName": "projects/dragondrop-dev/logs/cloudaudit.googleapis.com%2Factivity",
//            "receiveTimestamp": "2023-03-11T17:24:54.248853501Z"
//        },
//        {
//            "protoPayload": {
//                "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
//                "status": {},
//                "authenticationInfo": {
//                    "principalEmail": "goodman.benjamin@dragondrop.cloud"
//                },
//                "serviceName": "storage.googleapis.com",
//                "methodName": "storage.buckets.update",
//                "authorizationInfo": [
//                    {
//                        "resource": "projects/_/buckets/testing-out-this-bucket",
//                        "permission": "storage.buckets.update",
//                        "granted": true,
//                        "resourceAttributes": {}
//                    }
//                ],
//                "resourceName": "projects/_/buckets/testing-out-this-bucket"
//            },
//            "timestamp": "2023-03-08T17:24:02.675064542Z",
//            "severity": "NOTICE",
//            "logName": "projects/dragondrop-dev/logs/cloudaudit.googleapis.com%2Factivity",
//            "receiveTimestamp": "2023-03-08T17:24:03.418303345Z"
//        }
//    ]
//}`)
//	// When
//	output, err := glc.ExtractDataFromResourceResult(inputResourceResult, "", true)
//	if err != nil {
//		t.Errorf("unexpected error in test: %v", err)
//	}
//
//	// Then
//	expectedOutput := terraformValueObjects.ResourceActions{
//		Modifier: terraformValueObjects.CloudActorTimeStamp{
//			Actor:     terraformValueObjects.CloudActor("goodman.benjamin@dragondrop.cloud"),
//			Timestamp: terraformValueObjects.Timestamp("2023-03-11"),
//		},
//	}
//
//	if !reflect.DeepEqual(output, expectedOutput) {
//		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
//	}
//}
//
//func TestExtractDataFromResourceResult_ManagedDriftRecordResponse(t *testing.T) {
//	// Given
//	glc := GoogleLogQuerier{}
//	inputResourceResult := []byte(`{
//    "entries": [
//        {
//            "protoPayload": {
//                "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
//                "status": {},
//                "authenticationInfo": {
//                    "principalEmail": "goodman.benjamin@dragondrop.cloud"
//                },
//                "serviceName": "storage.googleapis.com",
//                "methodName": "storage.buckets.update",
//                "authorizationInfo": [
//                    {
//                        "resource": "projects/_/buckets/testing-out-this-bucket",
//                        "permission": "storage.buckets.update",
//                        "granted": true,
//                        "resourceAttributes": {}
//                    }
//                ],
//                "resourceName": "projects/_/buckets/testing-out-this-bucket"
//            },
//            "insertId": "yyo9azd9x9y",
//            "resource": {
//                "type": "gcs_bucket",
//                "labels": {
//                    "project_id": "dragondrop-dev",
//                    "bucket_name": "testing-out-this-bucket",
//                    "location": "us"
//                }
//            },
//            "timestamp": "2023-03-08T17:24:53.274706482Z",
//            "severity": "NOTICE",
//            "logName": "projects/dragondrop-dev/logs/cloudaudit.googleapis.com%2Factivity",
//            "receiveTimestamp": "2023-03-11T17:24:54.248853501Z"
//        }
//    ]
//}`)
//	// When
//	output, err := glc.ExtractDataFromResourceResult(inputResourceResult, "", false)
//	if err != nil {
//		t.Errorf("unexpected error in test: %v", err)
//	}
//
//	// Then
//	expectedOutput := terraformValueObjects.ResourceActions{
//		Modifier: terraformValueObjects.CloudActorTimeStamp{
//			Actor:     terraformValueObjects.CloudActor("goodman.benjamin@dragondrop.cloud"),
//			Timestamp: terraformValueObjects.Timestamp("2023-03-11"),
//		},
//	}
//
//	if !reflect.DeepEqual(output, expectedOutput) {
//		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
//	}
//}
