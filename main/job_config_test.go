package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/hclcreate"
	costEstimation "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/cost_estimation"
	identifyCloudActors "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/identify_cloud_actors"
	terraformImportMigrationGenerator "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_import_migration_generator"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	terraformWorkspace "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_workspace"
	terraformerCli "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraformer_executor/terraformer_cli"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/vcs"
)

func validJobConfig() *JobConfig {
	return &JobConfig{
		IsManagedDriftOnly: false,
		CloudRegions:       terraformValueObjects.CloudRegionsDecoder{"us-east1"},
		CloudCredential:    "{}",
		JobID:              "JobID",
		MigrationHistoryStorage: hclcreate.MigrationHistory{
			StorageType: "S3",
			Bucket:      "Bucket",
			Region:      "Region",
		},
		TerraformVersion:           "TerraformVersion",
		StateBackend:               "StateBackend",
		TerraformCloudOrganization: "TerraformCloudOrganization",
		TerraformCloudToken:        "TerraformCloudToken",
		WorkspaceDirectories:       terraformWorkspace.WorkspaceDirectoriesDecoder{"workspace1", "workspace2"},
		Provider: map[terraformValueObjects.Provider]string{
			"aws": "~>4.57.0",
		},
		VCSRepo:            "VCSRepo",
		VCSPat:             "xyz",
		InfracostToken:     "ico-mytoken",
		PullReviewers:      []string{"PullReviewer1", "PullReviewer2"},
		ResourcesWhiteList: terraformValueObjects.ResourceNameList{"Resource1", "Resource2"},
		ResourcesBlackList: terraformValueObjects.ResourceNameList{"Resource3", "Resource4"},
	}
}

func TestGetVCSConfig(t *testing.T) {
	// Given
	jobConfig := validJobConfig()

	// When
	got := jobConfig.getVCSConfig()

	// Then
	want := vcs.Config{
		VCSRepo:       jobConfig.VCSRepo,
		VCSPat:        jobConfig.VCSPat,
		PullReviewers: jobConfig.PullReviewers,
	}

	assert.Equal(t, want, got, "VCS Config should be equal")
}

func TestGetTerraformWorkspaceConfig(t *testing.T) {
	// Given
	jobConfig := validJobConfig()

	// When
	got := jobConfig.getTerraformWorkspaceConfig()

	// Then
	want := terraformWorkspace.TfStackConfig{
		Region:                     "us-east1",
		CloudCredential:            "{}",
		StateBackend:               jobConfig.StateBackend,
		TerraformCloudOrganization: jobConfig.TerraformCloudOrganization,
		TerraformCloudToken:        jobConfig.TerraformCloudToken,
		WorkspaceDirectories:       jobConfig.WorkspaceDirectories,
	}

	assert.Equal(t, want, got, "TerraformWorkspaceConfig should be equal")
}

func TestGetHCLCreateConfig(t *testing.T) {
	// Given
	jobConfig := validJobConfig()

	// When
	got := jobConfig.getHCLCreateConfig()

	// Then
	want := hclcreate.Config{
		MigrationHistoryStorage: jobConfig.MigrationHistoryStorage,
		TerraformVersion:        jobConfig.TerraformVersion,
	}

	assert.Equal(t, want, got, "HCLCreateConfig should be equal")
}

func TestGetTerraformerConfig(t *testing.T) {
	// Given
	jobConfig := validJobConfig()

	// When
	got := jobConfig.getTerraformerConfig()

	// Then
	want := terraformerCli.TerraformerExecutorConfig{
		CloudCredential:  "{}",
		Provider:         jobConfig.Provider,
		TerraformVersion: terraformValueObjects.Version(jobConfig.TerraformVersion),
		CloudRegions:     jobConfig.CloudRegions,
	}

	assert.Equal(t, want, got, "TerraformerExecutorConfig should be equal")
}

func TestGetTerraformerCLIConfig(t *testing.T) {
	// Given
	jobConfig := validJobConfig()

	// When
	got := jobConfig.getTerraformerCLIConfig()

	// Then
	want := terraformerCli.Config{
		ResourcesWhiteList: jobConfig.ResourcesWhiteList,
		ResourcesBlackList: jobConfig.ResourcesBlackList,
	}

	assert.Equal(t, want, got, "TerraformerCLIConfig should be equal")
}

func TestGetTerraformImportMigrationGeneratorConfig(t *testing.T) {
	// Given
	jobConfig := validJobConfig()

	// When
	got := jobConfig.getTerraformImportMigrationGeneratorConfig()

	// Then
	want := terraformImportMigrationGenerator.Config{
		CloudCredential: "{}",
	}

	assert.Equal(t, want, got, "TerraformImportMigrationGeneratorConfig should be equal")
}

func TestGetCostEstimationConfig(t *testing.T) {
	// Given
	jobConfig := validJobConfig()

	// When
	got := jobConfig.getCostEstimationConfig()

	// Then
	want := costEstimation.CostEstimatorConfig{
		CloudCredential:   "{}",
		InfracostAPIToken: jobConfig.InfracostToken,
	}

	assert.Equal(t, want, got, "CostEstimationConfig should be equal")
}

func TestGetIdentifyCloudActorsConfig(t *testing.T) {
	// Given
	jobConfig := validJobConfig()

	// When
	got := jobConfig.getIdentifyCloudActorsConfig()

	// Then
	want := identifyCloudActors.Config{
		CloudCredential: "{}",
	}

	assert.Equal(t, want, got, "IdentifyCloudActorsConfig should be equal")
}
