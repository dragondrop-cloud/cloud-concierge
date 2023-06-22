package hclcreate

import (
	"reflect"
	"testing"
)

func TestConvertTerraformerResourceName(t *testing.T) {
	input := "tfer--xyz-war_lkj"

	output := ConvertTerraformerResourceName(input)

	expectedOutput := "xyz_war_lkj"

	if output != expectedOutput {
		t.Errorf("got %v, expected %v", output, expectedOutput)
	}
}

func TestDecodeMigrationHistory(t *testing.T) {
	var envVar MigrationHistory

	err := envVar.Decode(`{"storageType": "s3", "bucket": "xyz", "region": "us-east1"}`)
	if err != nil {
		t.Errorf("unexpected error in envVar.Decode: %v", err)
	}

	expectedEnvVar := MigrationHistory{
		Bucket:      "xyz",
		Region:      "us-east1",
		StorageType: "s3",
	}
	if !reflect.DeepEqual(envVar, expectedEnvVar) {
		t.Errorf("got %v, expected %v", envVar, expectedEnvVar)
	}

	err = envVar.Decode(`{"storageType": "S3", "bucket": "xyz", "region": "us-east1"}`)
	if err != nil {
		t.Errorf("unexpected error in envVar.Decode: %v", err)
	}

	expectedEnvVar = MigrationHistory{
		Bucket:      "xyz",
		Region:      "us-east1",
		StorageType: "s3",
	}
	if !reflect.DeepEqual(envVar, expectedEnvVar) {
		t.Errorf("got %v, expected %v", envVar, expectedEnvVar)
	}

	err = envVar.Decode(`{"storageType": "Google Storage Bucket", "bucket": "xyz", "region": "us-east1"}`)
	if err != nil {
		t.Errorf("unexpected error in envVar.Decode: %v", err)
	}

	expectedEnvVar = MigrationHistory{
		Bucket:      "xyz",
		Region:      "us-east1",
		StorageType: "gcs",
	}
	if !reflect.DeepEqual(envVar, expectedEnvVar) {
		t.Errorf("got %v, expected %v", envVar, expectedEnvVar)
	}

	// Error case 1 - wrong storage type
	err = envVar.Decode(`{"storageType": "azureBlob", "bucket": "xyz", "region": "us-east1"}`)
	expectedError := "only types of 's3' and 'gcs' are currently supported. Attempted azureBlob"
	if err.Error() != expectedError {
		t.Errorf("got error %v, expected error %v", err, expectedError)
	}

	// Error case 2 - no region specified
	err = envVar.Decode(`{"storageType": "s3", "bucket": "xyz", "region": ""}`)
	expectedError = "region variable cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("got error %v, expected error %v", err, expectedError)
	}

	// Error case 3 - no bucket specified
	err = envVar.Decode(`{"storageType": "s3", "bucket": "", "region": "us-east1"}`)
	expectedError = "the required field `bucket` is not present"
	if err.Error() != expectedError {
		t.Errorf("got error %v, expected error %v", err, expectedError)
	}
}
