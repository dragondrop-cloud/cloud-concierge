package hclcreate

import (
	"reflect"
	"strconv"
	"testing"
)

func TestResourceToIdentifierStruct(t *testing.T) {
	h := hclCreate{}
	input := "google-project_id.google_storage_bucket.tfer--resource-id"

	output := h.resourceToIdentifierStruct(input)

	expectedOutput := ResourceIdentifier{
		division:     "google-project_id",
		resourceType: "google_storage_bucket",
		resourceName: "tfer--resource-id",
	}

	reflect.DeepEqual(output, expectedOutput)
}

func TestGenerateImportStatement(t *testing.T) {
	h := hclCreate{}

	inputResource := "google-dev.tf_type_abc.tfer--tf_name_xyz"

	resourceImportsByDivision := ResourceImportsByDivision{
		"google-dev": {
			"tf_type_123.tfer--tf_name_xyz": {
				TerraformConfigLocation: "tf_type_123.tfer--tf_name_xyz",
				RemoteCloudReference:    "import_1",
			},
			"tf_type_abc.tfer--tf_name_123": {
				TerraformConfigLocation: "tf_type_abc.tfer--tf_name_123",
				RemoteCloudReference:    "import_2",
			},
			"tf_type_abc.tfer--tf_name_xyz": {
				TerraformConfigLocation: "tf_type_abc.tfer--tf_name_xyz",
				RemoteCloudReference:    "import_3",
			},
		},
		"google-prod": {
			"tf_type_123.tfer--tf_name_xyz": {
				TerraformConfigLocation: "tf_type_123.tfer--tf_name_xyz",
				RemoteCloudReference:    "import_1",
			},
		},
	}

	expectedOutput := "import tf_type_abc.tf_name_xyz import_3"

	output, err := h.generateImportStatement(inputResource, resourceImportsByDivision)
	if err != nil {
		t.Errorf("unexpected error in h.generateImportStatement: %v", err)
	}

	if expectedOutput != output {
		t.Errorf("got: %v\n\nexpected: %v", output, expectedOutput)
	}
}

func TestGenerateImportStatementText(t *testing.T) {
	h := hclCreate{}

	resourceCloudID := "example-id"

	inputResourceID := ResourceIdentifier{
		resourceName: "tfer--example-name_broski",
		resourceType: "resource_type",
	}

	output := h.generateImportStatementText(resourceCloudID, inputResourceID)
	expectedOutput := "import resource_type.example_name_broski example-id"

	if output != expectedOutput {
		t.Errorf("got %v, expected %v", output, expectedOutput)
	}

}

func TestIndividualTFMigrateConfigS3(t *testing.T) {
	h := hclCreate{
		config: Config{
			MigrationHistoryStorage: MigrationHistory{
				StorageType: "s3",
				Bucket:      "example-bucket",
				Region:      "us-east1",
			},
		},
	}

	expectedOutput := "tfmigrate {\n  migration_dir              = \"./cloud-concierge/tfmigrate/\"\n  " +
		"is_backend_terraform_cloud = true\n  " +
		"history {\n    storage \"s3\" {\n      bucket = \"example-bucket\"\n      " +
		"key    = \"exampleWorkspace/history.json\"\n      region = \"us-east1\"\n    }\n  }\n}\n"

	output, err := h.individualTFMigrateConfig(
		"exampleWorkspace",
	)
	if err != nil {
		t.Errorf("unexpected error in h.IndividualTFMigrateConfig: %v", err)
	}

	outputString := string(output)

	if expectedOutput != outputString {
		t.Errorf(
			"got:\n%v\n\nexpected:\n%v",
			strconv.Quote(outputString),
			strconv.Quote(expectedOutput),
		)
	}
}

func TestIndividualTFMigrateConfigGCS(t *testing.T) {
	h := hclCreate{
		config: Config{
			MigrationHistoryStorage: MigrationHistory{
				StorageType: "gcs",
				Bucket:      "example-bucket",
				Region:      "",
			},
		},
	}

	expectedOutput := "tfmigrate {\n  migration_dir              = \"./cloud-concierge/tfmigrate/\"\n  is_backend_terraform_cloud = true\n  " +
		"history {\n    storage \"gcs\" {\n      bucket = \"example-bucket\"\n      " +
		"name   = \"exampleWorkspace/history.json\"\n    }\n  }\n}\n"
	output, err := h.individualTFMigrateConfig(
		"exampleWorkspace",
	)
	if err != nil {
		t.Errorf("unexpected error in h.IndividualTFMigrateConfig: %v", err)
	}

	outputString := string(output)

	if expectedOutput != outputString {
		t.Errorf(
			"got:\n%v\n\nexpected:\n%v",
			strconv.Quote(outputString),
			strconv.Quote(expectedOutput),
		)
	}
}

func TestIndividualTFMigrateMigration(t *testing.T) {
	h := hclCreate{}

	resourceImportsByDivision := ResourceImportsByDivision{
		"google-dev": {
			"tf_type_123.tfer--tf_name_xyz": {
				TerraformConfigLocation: "tf_type_123.tfer--tf_name_xyz",
				RemoteCloudReference:    "import_1a",
			},
			"tf_type_abc.tfer--tf_name_123": {
				TerraformConfigLocation: "tf_type_abc.tfer--tf_name_123",
				RemoteCloudReference:    "import_2",
			},
			"tf_type_abc.tfer--tf_name_xyz": {
				TerraformConfigLocation: "tf_type_abc.tfer--tf_name_xyz",
				RemoteCloudReference:    "import_3",
			},
		},
		"google-prod": {
			"tf_type_123.tfer--tf_name_xyz_prod": {
				TerraformConfigLocation: "tf_type_123.tfer--tf_name_xyz",
				RemoteCloudReference:    "import_1b",
			},
		},
	}

	newResourceToWorkspace := NewResourceToWorkspace{
		"google-dev.tf_type_123.tfer--tf_name_xyz":       "workspace_1",
		"google-dev.tf_type_abc.tfer--tf_name_123":       "workspace_2",
		"google-dev.tf_type_abc.tfer--tf_name_xyz":       "workspace_2",
		"google-prod.tf_type_123.tfer--tf_name_xyz_prod": "workspace_1",
	}

	expectedOutput := "migration \"state\" \"import\" {\n  dir       = \"/github/workspace/xyz\"\n  workspace = \"workspace_1\"\n  actions   = [\"import tf_type_123.tf_name_xyz_prod import_1b\", \"import tf_type_123.tf_name_xyz import_1a\"]\n}\n"

	expectedOutputTwo := "migration \"state\" \"import\" {\n  dir       = \"/github/workspace/xyz\"\n  workspace = \"workspace_1\"\n  actions   = [\"import tf_type_123.tf_name_xyz import_1a\", \"import tf_type_123.tf_name_xyz_prod import_1b\"]\n}\n"

	output, err := h.individualTFMigrateMigration("/xyz", "workspace_1", resourceImportsByDivision, newResourceToWorkspace)
	if err != nil {
		t.Errorf("unexpected error in h.individualTFMigrateMigration(): %v", err)
	}

	stringOutput := string(output)

	if stringOutput != expectedOutput && stringOutput != expectedOutputTwo {
		t.Errorf("got:\n%v\n\nexpected:\n%v", strconv.Quote(stringOutput), strconv.Quote(expectedOutput))
	}
}
