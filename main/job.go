package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	costEstimation "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/cost_estimation"
	dragonDrop "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/dragon_drop"
	identifyCloudActors "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/identify_cloud_actors"
	resourcesCalculator "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/resources_calculator"
	resourcesWriter "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/resources_writer"
	terraformImportMigrationGenerator "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_import_migration_generator"
	terraformManagedResourcesDriftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector"
	terraformSecurity "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_security"
	terraformWorkspace "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_workspace"
	terraformerExecutor "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraformer_executor"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/vcs"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// Job is an instance of a runnable dragondrop job.
type Job struct {
	// vcs is the implementation of interfaces.VCS for interacting with a remote version control system
	vcs interfaces.VCS

	// terraformWorkspace is the implementation of interfaces.TerraformWorkspace for pulling information from
	// remote terraform state files.
	terraformWorkspace interfaces.TerraformWorkspace

	// terraformerExecutor is the implementation of interfaces.TerraformerExecutor for pulling information on
	// the current cloud environment using terraformer programmatically.
	terraformerExecutor interfaces.TerraformerExecutor

	// terraformerExecutor is the implementation of interfaces.TerraformImportMigrationGenerator for
	// generating terraform state import statements corresponding to identified resources.
	terraformImportMigrationGenerator interfaces.TerraformImportMigrationGenerator

	// resourcesCalculator is the implementation of interfaces.ResourcesCalculator for matching identified resources
	// to the appropriate state file.
	resourcesCalculator interfaces.ResourcesCalculator

	// resourcesWriter is the implementation of interfaces.ResourcesWriter for writing out resource
	// definitions programmatically.
	resourcesWriter interfaces.ResourcesWriter

	// dragonDrop is the implementation of interfaces.DragonDrop for interacting with
	// dragondrop API endpoints.
	dragonDrop interfaces.DragonDrop

	// identifyCloudActors is the implementation of interfaces.IdentifyCloudActors for determining
	// which cloud actor made resource changes outside of Terraform control.
	identifyCloudActors interfaces.IdentifyCloudActors

	// costEstimator is the implementation of interfaces.CostEstimation for calculating monthly cost
	// estimates of identified cloud resources.
	costEstimator interfaces.CostEstimation

	// driftDetector is the implementation of interfaces.TerraformManagedResourcesDriftDetector for detect the
	// drifted resources that already exists in the cloud state.
	driftDetector interfaces.TerraformManagedResourcesDriftDetector

	// terraformSecurity
	terraformSecurity interfaces.TerraformSecurity

	// name is the name of the current job
	name string

	// noNewResources is a flag to know if there are no new resources
	noNewResources bool

	// config is the configuration to run successfully the job
	config JobConfig
}

// Authorize ensures that the Job is valid by checking against the dragondrop
// API.
func (j *Job) Authorize(ctx context.Context) error {
	// For a Job managed by the dragondrop platform, we authorize and update the job name
	if j.config.JobID != "empty" && j.config.JobID != "" {
		err := j.dragonDrop.CheckLoggerAndToken(ctx)
		if err != nil {
			return fmt.Errorf("[create_job][error checking logger and token][%w]", err)
		}

		err = j.dragonDrop.InformStarted(ctx)
		if err != nil {
			return fmt.Errorf("[create_job][error informing started job][%w]", err)
		}

		jobName, err := j.dragonDrop.AuthorizeManagedJob(ctx)
		if err != nil {
			return fmt.Errorf("[create_job][error authorizing managed job][%w]", err)
		}
		j.name = jobName
		j.dragonDrop.PostLog(ctx, "Authorized against billing plan.")
	} else {
		err := j.dragonDrop.AuthorizeJob(ctx)
		if err != nil {
			fmt.Printf("Error authenticating the job run, please get an Organization token by signing up at https://app.dragondrop.cloud.")
			return fmt.Errorf("[create_job][error authorizing job][%w]", err)
		}
		j.name = j.config.JobName
	}
	return nil
}

// Run runs an instance of the Job struct to completion by coordinating calls to different
// interface implementations within the Job.
func (j *Job) Run(ctx context.Context) error {
	err := j.vcs.Clone()
	if err != nil {
		return fmt.Errorf("[run_job][error clonning repo][%w]", err)
	}

	err = j.dragonDrop.InformRepositoryCloned(ctx)
	if err != nil {
		return fmt.Errorf("[run_job][error posting cloned status]%w", err)
	}

	workspaceToDirectory, err := j.terraformWorkspace.FindTerraformWorkspaces(ctx)
	if err != nil {
		return fmt.Errorf("[run_job][error finding terraform workspaces][%w]", err)
	}

	err = j.terraformWorkspace.DownloadWorkspaceState(ctx, workspaceToDirectory)
	if err != nil {
		return fmt.Errorf("[run_job][error downloading workspace state][%w]", err)
	}

	err = j.terraformerExecutor.Execute(ctx)
	if err != nil {
		return fmt.Errorf("[run_job][error setting up terraformer executor][%w]", err)
	}

	err = j.terraformImportMigrationGenerator.Execute(ctx)
	if err != nil {
		return fmt.Errorf("[run_job][error executing terraform import][%w]", err)
	}

	if !j.config.IsManagedDriftOnly {
		err = j.resourcesCalculator.Execute(ctx, workspaceToDirectory)
		if err != nil {
			if errors.Unwrap(errors.Unwrap(err)) != resourcesCalculator.ErrNoNewResources {
				return fmt.Errorf("[run_job][error calculating resources][%w]", err)
			}

			j.noNewResources = true
			log.Warnf("Did not find new resources, but scanning for drifted resources")
		}
	} else {
		j.noNewResources = true
	}

	driftedResourcesIdentified, err := j.driftDetector.Execute(ctx, workspaceToDirectory)
	if err != nil {
		return fmt.Errorf("[run_job][error detecting drifted resources]%w", err)
	}

	err = j.dragonDrop.InformCloudActorIdentification(ctx)
	if err != nil {
		return fmt.Errorf("[run_job][error posting cloud actor identification status]%w", err)
	}

	err = j.identifyCloudActors.Execute(ctx)
	if err != nil {
		return fmt.Errorf("[run_job][error identifying cloud actors]%w", err)
	}

	err = j.dragonDrop.InformCostEstimation(ctx)
	if err != nil {
		return fmt.Errorf("[run_job][error posting cost estimation status]%w", err)
	}

	err = j.costEstimator.Execute(ctx)
	if err != nil {
		return fmt.Errorf("[run_job][error estimating cost for identified resources]%w", err)
	}

	err = j.dragonDrop.InformSecurityScan(ctx)
	if err != nil {
		return fmt.Errorf("[run_job][error posting security scan status]%w", err)
	}

	err = j.terraformSecurity.ExecuteScan(ctx)
	if err != nil {
		return fmt.Errorf("[run_job][error executing the tfsec command]%w", err)
	}

	createDummyFile := driftedResourcesIdentified && j.noNewResources
	prURL, err := j.resourcesWriter.Execute(ctx, j.name, createDummyFile, workspaceToDirectory)
	if err != nil {
		return fmt.Errorf("[run_job][error writing resources on vcs][%w]", err)
	}

	err = j.dragonDrop.PutJobPullRequestURL(ctx, prURL)
	if err != nil {
		return fmt.Errorf("[run_job][error putting job pull request URL][%v]", err)
	}

	err = j.dragonDrop.InformComplete(ctx)
	if err != nil {
		return fmt.Errorf("[run_job][error informing complete status][%w]", err)
	}

	return nil
}

// InitializeJobDependencies instantiates interface implementations for all needed interfaces
// and configures by pulling in environment variables.
func InitializeJobDependencies(ctx context.Context, env string) (*Job, error) {
	var jobConfig JobConfig
	err := envconfig.Process("CLOUDCONCIERGE", &jobConfig)
	if err != nil {
		return nil, fmt.Errorf("[cannot create job config]%w", err)
	}

	err = validateJobConfig(jobConfig)
	if err != nil {
		return nil, fmt.Errorf("[invalid job config]%w", err)
	}

	inferredData, err := getInferredData(jobConfig)
	if err != nil {
		log.Errorf("[cannot create job config]%s", err.Error())
		return nil, fmt.Errorf("[cannot create job config]%w", err)
	}

	dragonDropInstance, err := (&dragonDrop.Factory{}).Instantiate(env, jobConfig.getDragonDropConfig())
	if err != nil {
		return nil, err
	}
	vcsInstance, err := (&vcs.Factory{}).Instantiate(ctx, env, dragonDropInstance, jobConfig.getVCSConfig(), inferredData.VCSSystem)
	if err != nil {
		return nil, err
	}
	workspace, err := (&terraformWorkspace.Factory{}).Instantiate(ctx, env, dragonDropInstance, jobConfig.getTerraformWorkspaceConfig())
	if err != nil {
		return nil, err
	}
	executor, err := (&terraformerExecutor.Factory{}).Instantiate(ctx, env, dragonDropInstance, inferredData.Provider,
		jobConfig.getHCLCreateConfig(), jobConfig.getTerraformerConfig(), jobConfig.getTerraformerCLIConfig())
	if err != nil {
		return nil, err
	}
	instantiate, err := (&terraformImportMigrationGenerator.Factory{}).Instantiate(ctx, env, dragonDropInstance, inferredData.Provider,
		jobConfig.getTerraformImportMigrationGeneratorConfig())
	if err != nil {
		return nil, err
	}
	calculator, err := (&resourcesCalculator.Factory{}).Instantiate(ctx, env, dragonDropInstance, inferredData.Provider)
	if err != nil {
		return nil, err
	}
	costEstimator, err := (&costEstimation.Factory{}).Instantiate(env, inferredData.Provider, jobConfig.getCostEstimationConfig())
	if err != nil {
		return nil, err
	}
	identifier, err := (&identifyCloudActors.Factory{}).Instantiate(ctx, env, dragonDropInstance, inferredData.Provider, jobConfig.getIdentifyCloudActorsConfig())
	if err != nil {
		return nil, err
	}
	writer, err := (&resourcesWriter.Factory{}).Instantiate(ctx, env, vcsInstance, dragonDropInstance, inferredData.Provider, jobConfig.getHCLCreateConfig())
	if err != nil {
		return nil, err
	}
	driftDetector, err := (&terraformManagedResourcesDriftDetector.Factory{}).Instantiate(ctx, env, inferredData.Provider)
	if err != nil {
		return nil, err
	}
	tfSec, err := (&terraformSecurity.Factory{}).Instantiate(ctx, env, inferredData.Provider)
	if err != nil {
		return nil, err
	}

	return &Job{
		vcs:                               vcsInstance,
		terraformWorkspace:                workspace,
		terraformerExecutor:               executor,
		terraformImportMigrationGenerator: instantiate,
		resourcesCalculator:               calculator,
		resourcesWriter:                   writer,
		dragonDrop:                        dragonDropInstance,
		costEstimator:                     costEstimator,
		identifyCloudActors:               identifier,
		driftDetector:                     driftDetector,
		config:                            jobConfig,
		terraformSecurity:                 tfSec,
	}, nil
}
