package main

import (
	"fmt"
	"strings"

	nlpenginerequestor "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/nlp_engine_requester"
	"github.com/sirupsen/logrus"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/hclcreate"
	costEstimation "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/cost_estimation"
	identifyCloudActors "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/identify_cloud_actors"
	terraformImportMigrationGenerator "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_import_migration_generator"
	driftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector/drift_detector"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	terraformWorkspace "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_workspace"
	terraformerCli "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraformer_executor/terraformer_cli"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/vcs"
)

// JobConfig is the configuration for the Job that contains the variables to run successfully
type JobConfig struct {
	// IsManagedDriftOnly represents the option for the user to only scan drifted resources and not new resources
	IsManagedDriftOnly bool `default:"false"`

	// CloudCredential is a cloud credential that is used to authenticate with a cloud provider. Credential should
	// only require read-only access.
	CloudCredential terraformValueObjects.Credential `required:"false"`

	// Division is the name of a cloud division. In AWS this is an account, in GCP this is a project name, and in Azure this is a subscription.
	Division terraformValueObjects.Division `required:"true"`

	// InfracostToken is the token for accessing Infracost's cloud-pricing API.
	InfracostToken string `required:"true"`

	// JobID is the unique identification string for the current job run.
	JobID string `default:"empty"`

	// JobName is the name of the job.
	JobName string `default:"Cloud Concierge Report"`

	// MigrationHistoryStorage is a map containing information needed for specifying tfmigrate
	// history storage appropriately.
	MigrationHistoryStorage hclcreate.MigrationHistory

	// NLPEndpoint is the endpoint for the NLP service used by cloud-concierge to match uncontrolled resources
	// to the right state files.
	NLPEndpoint string `default:"https://us-east4-dragondrop-prod.cloudfunctions.net/nlpengine-endpoint-prod"`

	// TerraformVersion is the version of Terraform used.
	TerraformVersion string `required:"true"`

	// StateBackend is the name of the backend used for storing State.
	StateBackend string `required:"true"`

	// TerraformCloudOrganization is the name of the organization within Terraform Cloud
	TerraformCloudOrganization string

	// TerraformCloudToken is the auth token to access Terraform Cloud programmatically.
	TerraformCloudToken string

	// WorkspaceDirectories is a slice of directories that contains terraform workspaces within the user repo.
	WorkspaceDirectories terraformWorkspace.WorkspaceDirectoriesDecoder `required:"true"`

	// Provider is a map between a cloud provider and the version for that provider.
	Provider map[terraformValueObjects.Provider]string `required:"true"`

	// VCSRepo is the full path of the repo containing a customer's infrastructure specification.
	// At the moment, must be a valid GitHub repository URL.
	VCSRepo string `required:"true"`

	// VCSPat is the personal access token for the VCS where a Pull Request should be output.
	VCSPat string `required:"true"`

	// PullReviewers is the name of the pull request reviewer who will be tagged on the opened pull request.
	PullReviewers []string `default:"NoReviewer"`

	// ResourcesWhiteList represents the list of resource names that will be exclusively considered for inclusion in the import statement.
	ResourcesWhiteList terraformValueObjects.ResourceNameList

	// ResourcesBlackList represents the list of resource names that will be excluded from consideration for inclusion in the import statement.
	ResourcesBlackList terraformValueObjects.ResourceNameList

	// CloudRegions represents the list of cloud regions that will be considered for inclusion in the import statement.
	CloudRegions terraformValueObjects.CloudRegionsDecoder `default:"['us-east-1]"`
}

// validateJobConfig validates the JobConfig struct with the values as expected.
func validateJobConfig(config JobConfig) error {
	logrus.Debugf("Validating job config: %+v", config)

	if strings.ToLower(config.StateBackend) == "terraformcloud" {
		if config.TerraformCloudOrganization == "" {
			return fmt.Errorf("[terraform cloud organization is required when using terraform cloud as state backend]")
		}
		if config.TerraformCloudToken == "" {
			return fmt.Errorf("[terraform cloud token is required when using terraform cloud as state backend]")
		}
	}
	return nil
}

func (c JobConfig) getVCSConfig() vcs.Config {
	return vcs.Config{
		VCSRepo:       c.VCSRepo,
		VCSPat:        c.VCSPat,
		PullReviewers: c.PullReviewers,
	}
}

func (c JobConfig) getTerraformWorkspaceConfig() terraformWorkspace.TfStackConfig {
	return terraformWorkspace.TfStackConfig{
		Region:                     string(c.CloudRegions[0]),
		CloudCredential:            c.CloudCredential,
		StateBackend:               c.StateBackend,
		TerraformCloudOrganization: c.TerraformCloudOrganization,
		TerraformCloudToken:        c.TerraformCloudToken,
		WorkspaceDirectories:       c.WorkspaceDirectories,
	}
}

func (c JobConfig) getHCLCreateConfig() hclcreate.Config {
	return hclcreate.Config{
		MigrationHistoryStorage: c.MigrationHistoryStorage,
		TerraformVersion:        c.TerraformVersion,
	}
}

func (c JobConfig) getTerraformerConfig() terraformerCli.TerraformerExecutorConfig {
	return terraformerCli.TerraformerExecutorConfig{
		CloudCredential:  c.CloudCredential,
		Division:         c.Division,
		Provider:         c.Provider,
		TerraformVersion: terraformValueObjects.Version(c.TerraformVersion),
		CloudRegions:     c.CloudRegions,
	}
}

func (c JobConfig) getTerraformerCLIConfig() terraformerCli.Config {
	return terraformerCli.Config{
		ResourcesWhiteList: c.ResourcesWhiteList,
		ResourcesBlackList: c.ResourcesBlackList,
	}
}

func (c JobConfig) getTerraformImportMigrationGeneratorConfig() terraformImportMigrationGenerator.Config {
	return terraformImportMigrationGenerator.Config{
		CloudCredential: c.CloudCredential,
		Division:        c.Division,
	}
}

func (c JobConfig) getCostEstimationConfig() costEstimation.CostEstimatorConfig {
	return costEstimation.CostEstimatorConfig{
		CloudCredential: c.CloudCredential,
	}
}

func (c JobConfig) getIdentifyCloudActorsConfig() identifyCloudActors.Config {
	return identifyCloudActors.Config{
		CloudCredential: c.CloudCredential,
		Division:        c.Division,
	}
}

func (c JobConfig) getManagedResourceDriftDetectorConfig() driftDetector.ManagedResourceDriftDetectorConfig {
	return driftDetector.ManagedResourceDriftDetectorConfig{
		ResourcesWhiteList: c.ResourcesWhiteList,
		ResourcesBlackList: c.ResourcesBlackList,
	}
}

func (c JobConfig) getNLPEngineConfig() nlpenginerequestor.HTTPNLPEngineClientConfig {
	return nlpenginerequestor.HTTPNLPEngineClientConfig{
		NLPEndpoint:          c.NLPEndpoint,
		VCSRepo:              c.VCSRepo,
		WorkspaceDirectories: c.WorkspaceDirectories,
	}
}
