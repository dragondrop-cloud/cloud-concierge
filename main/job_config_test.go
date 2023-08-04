package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/hclcreate"
	costEstimation "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/cost_estimation"
	dragonDrop "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/dragon_drop"
	identifyCloudActors "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/identify_cloud_actors"
	terraformImportMigrationGenerator "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_import_migration_generator"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	terraformWorkspace "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_workspace"
	terraformerCli "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraformer_executor/terraformer_cli"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/vcs"
)

func validJobConfig() *JobConfig {
	return &JobConfig{
		IsManagedDriftOnly:         false,
		InfracostAPIToken:          "InfracostAPIToken",
		APIPath:                    "https://api.dragondrop.cloud",
		JobID:                      "JobID",
		OrgToken:                   "OrgToken",
		MigrationHistoryStorage:    hclcreate.MigrationHistory{ /* Valor necesario */ },
		TerraformVersion:           "TerraformVersion",
		StateBackend:               "StateBackend",
		TerraformCloudOrganization: "TerraformCloudOrganization",
		TerraformCloudToken:        "TerraformCloudToken",
		WorkspaceDirectories:       terraformWorkspace.WorkspaceDirectoriesDecoder{ /* Valor necesario */ },
		Provider: map[terraformValueObjects.Provider]string{
			"aws": "~>4.57.0",
		},
		VCSToken:           "VCSToken",
		VCSUser:            "VCSUser",
		VCSRepo:            "VCSRepo",
		PullReviewers:      []string{"PullReviewer1", "PullReviewer2"},
		ResourcesWhiteList: terraformValueObjects.ResourceNameList{ /* Valor necesario */ },
		ResourcesBlackList: terraformValueObjects.ResourceNameList{ /* Valor necesario */ },
	}
}

func TestGetDragonDropConfig(t *testing.T) {
	// Given
	jobConfig := validJobConfig()

	// When
	dragonDropConfig := jobConfig.getDragonDropConfig()

	// Then
	assert.Equal(t, jobConfig.APIPath, dragonDropConfig.APIPath, "APIPath should be equal")
	assert.Equal(t, jobConfig.JobID, dragonDropConfig.JobID, "JobID should be equal")
	assert.Equal(t, jobConfig.OrgToken, dragonDropConfig.OrgToken, "OrgToken should be equal")

	want := dragonDrop.HTTPDragonDropClientConfig{
		APIPath:  jobConfig.APIPath,
		JobID:    jobConfig.JobID,
		OrgToken: jobConfig.OrgToken,
	}

	assert.Equal(t, want, dragonDropConfig, "HTTPDragonDropClientConfig should be equal")
}

func TestGetVCSConfig(t *testing.T) {
	// Given
	jobConfig := validJobConfig()

	// When
	got := jobConfig.getVCSConfig()

	// Then
	want := vcs.Config{
		VCSRepo:       jobConfig.VCSRepo,
		VCSToken:      jobConfig.VCSToken,
		VCSUser:       jobConfig.VCSUser,
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
	want := terraformWorkspace.TerraformCloudConfig{
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
	want := terraformImportMigrationGenerator.Config{}

	assert.Equal(t, want, got, "TerraformImportMigrationGeneratorConfig should be equal")
}

func TestGetCostEstimationConfig(t *testing.T) {
	// Given
	jobConfig := validJobConfig()

	// When
	got := jobConfig.getCostEstimationConfig()

	// Then
	want := costEstimation.CostEstimatorConfig{
		InfracostAPIToken: jobConfig.InfracostAPIToken,
	}

	assert.Equal(t, want, got, "CostEstimationConfig should be equal")
}

func TestGetIdentifyCloudActorsConfig(t *testing.T) {
	// Given
	jobConfig := validJobConfig()

	// When
	got := jobConfig.getIdentifyCloudActorsConfig()

	// Then
	want := identifyCloudActors.Config{}

	assert.Equal(t, want, got, "IdentifyCloudActorsConfig should be equal")
}
