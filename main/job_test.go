package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	resourcesCalculator "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/resources_calculator"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	. "github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

func TestAuthorize_Success(t *testing.T) {
	// Given
	vcs := new(VCSMock)
	terraformWorkspace := new(TerraformWorkspaceMock)
	terraformerExecutor := new(TerraformerExecutorMock)
	terraformImportMigrationGenerator := new(TerraformImportMigrationGeneratorMock)
	resourcesCalculator := new(ResourcesCalculatorMock)
	resourcesWriter := new(ResourcesWriterMock)
	dragonDrop := new(DragonDropMock)
	identifyCloudActors := new(IdentifyCloudActorsMock)
	costEstimator := new(CostEstimationMock)

	ctx := context.Background()

	// When
	dragonDrop.On("CheckLoggerAndToken", ctx).Return(nil)
	dragonDrop.On("InformStarted", ctx).Return(nil)
	dragonDrop.On("AuthorizeManagedJob", ctx).Return("job_name", nil)

	job := Job{
		costEstimator:                     costEstimator,
		dragonDrop:                        dragonDrop,
		identifyCloudActors:               identifyCloudActors,
		resourcesCalculator:               resourcesCalculator,
		resourcesWriter:                   resourcesWriter,
		terraformerExecutor:               terraformerExecutor,
		terraformImportMigrationGenerator: terraformImportMigrationGenerator,
		terraformWorkspace:                terraformWorkspace,
		vcs:                               vcs,
	}
	err := job.Authorize(ctx)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, job)
}

func TestNotCreateJob_WithInvalidToken(t *testing.T) {
	// Given
	vcs := new(VCSMock)
	terraformWorkspace := new(TerraformWorkspaceMock)
	terraformerExecutor := new(TerraformerExecutorMock)
	terraformImportMigrationGenerator := new(TerraformImportMigrationGeneratorMock)
	resourcesCalculator := new(ResourcesCalculatorMock)
	resourcesWriter := new(ResourcesWriterMock)
	dragonDrop := new(DragonDropMock)
	identifyCloudActors := new(IdentifyCloudActorsMock)
	costEstimator := new(CostEstimationMock)

	ctx := context.Background()
	checkLoggerAndTokenErr := errors.New("error checking job token")

	// When
	dragonDrop.On("CheckLoggerAndToken", ctx).Return(checkLoggerAndTokenErr)

	job := Job{
		costEstimator:                     costEstimator,
		dragonDrop:                        dragonDrop,
		identifyCloudActors:               identifyCloudActors,
		resourcesCalculator:               resourcesCalculator,
		resourcesWriter:                   resourcesWriter,
		terraformerExecutor:               terraformerExecutor,
		terraformImportMigrationGenerator: terraformImportMigrationGenerator,
		terraformWorkspace:                terraformWorkspace,
		vcs:                               vcs,
	}
	err := job.Authorize(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, checkLoggerAndTokenErr, errors.Unwrap(err))
}

func TestNotCreateJob_CannotInformStarted(t *testing.T) {
	// Given
	vcs := new(VCSMock)
	terraformWorkspace := new(TerraformWorkspaceMock)
	terraformerExecutor := new(TerraformerExecutorMock)
	terraformImportMigrationGenerator := new(TerraformImportMigrationGeneratorMock)
	resourcesCalculator := new(ResourcesCalculatorMock)
	resourcesWriter := new(ResourcesWriterMock)
	dragonDrop := new(DragonDropMock)
	identifyCloudActors := new(IdentifyCloudActorsMock)
	costEstimator := new(CostEstimationMock)

	ctx := context.Background()
	informStartedErr := errors.New("informing job started error")

	// When
	dragonDrop.On("CheckLoggerAndToken", ctx).Return(nil)
	dragonDrop.On("InformStarted", ctx).Return(informStartedErr)

	job := Job{
		costEstimator:                     costEstimator,
		dragonDrop:                        dragonDrop,
		identifyCloudActors:               identifyCloudActors,
		resourcesCalculator:               resourcesCalculator,
		resourcesWriter:                   resourcesWriter,
		terraformerExecutor:               terraformerExecutor,
		terraformImportMigrationGenerator: terraformImportMigrationGenerator,
		terraformWorkspace:                terraformWorkspace,
		vcs:                               vcs,
	}
	err := job.Authorize(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, informStartedErr, errors.Unwrap(err))
}

func TestNotCreateJob_UnauthorizedJob(t *testing.T) {
	// Given
	vcs := new(VCSMock)
	terraformWorkspace := new(TerraformWorkspaceMock)
	terraformerExecutor := new(TerraformerExecutorMock)
	terraformImportMigrationGenerator := new(TerraformImportMigrationGeneratorMock)
	resourcesCalculator := new(ResourcesCalculatorMock)
	resourcesWriter := new(ResourcesWriterMock)
	dragonDrop := new(DragonDropMock)
	identifyCloudActors := new(IdentifyCloudActorsMock)
	costEstimator := new(CostEstimationMock)

	ctx := context.Background()

	authJobErr := errors.New("cannot authorize job")

	// When
	dragonDrop.On("CheckLoggerAndToken", ctx).Return(nil)
	dragonDrop.On("InformStarted", ctx).Return(nil)
	dragonDrop.On("AuthorizeManagedJob", ctx).Return("job_name", authJobErr)

	job := Job{
		costEstimator:                     costEstimator,
		dragonDrop:                        dragonDrop,
		identifyCloudActors:               identifyCloudActors,
		resourcesCalculator:               resourcesCalculator,
		resourcesWriter:                   resourcesWriter,
		terraformerExecutor:               terraformerExecutor,
		terraformImportMigrationGenerator: terraformImportMigrationGenerator,
		terraformWorkspace:                terraformWorkspace,
		vcs:                               vcs,
	}
	err := job.Authorize(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, authJobErr, errors.Unwrap(err))
}

type JobDependenciesMock struct {
	vcs                               *VCSMock
	terraformWorkspace                *TerraformWorkspaceMock
	terraformerExecutor               *TerraformerExecutorMock
	terraformImportMigrationGenerator *TerraformImportMigrationGeneratorMock
	resourcesCalculator               *ResourcesCalculatorMock
	resourcesWriter                   *ResourcesWriterMock
	dragonDrop                        *DragonDropMock
	identifyCloudActors               *IdentifyCloudActorsMock
	costEstimator                     *CostEstimationMock
	driftDetector                     *TerraformManagedResourcesDriftDetectorMock
	terraformSecurity                 *TerraformSecurityMock
}

func createValidJob(t *testing.T) (*JobDependenciesMock, *Job) {
	vcs := new(VCSMock)
	terraformWorkspace := new(TerraformWorkspaceMock)
	terraformerExecutor := new(TerraformerExecutorMock)
	terraformImportMigrationGenerator := new(TerraformImportMigrationGeneratorMock)
	resourcesCalculator := new(ResourcesCalculatorMock)
	resourcesWriter := new(ResourcesWriterMock)
	dragonDrop := new(DragonDropMock)
	identifyCloudActors := new(IdentifyCloudActorsMock)
	costEstimator := new(CostEstimationMock)
	driftDetector := new(TerraformManagedResourcesDriftDetectorMock)
	tfSec := new(TerraformSecurityMock)

	ctx := context.Background()
	dragonDrop.On("CheckLoggerAndToken", ctx).Return(nil)
	dragonDrop.On("InformStarted", ctx).Return(nil)
	dragonDrop.On("AuthorizeManagedJob", ctx).Return("job_name", nil)
	dragonDrop.On("InformRepositoryCloned", ctx).Return(nil)
	job := &Job{
		costEstimator:                     costEstimator,
		dragonDrop:                        dragonDrop,
		resourcesCalculator:               resourcesCalculator,
		resourcesWriter:                   resourcesWriter,
		terraformerExecutor:               terraformerExecutor,
		terraformImportMigrationGenerator: terraformImportMigrationGenerator,
		terraformWorkspace:                terraformWorkspace,
		vcs:                               vcs,
		identifyCloudActors:               identifyCloudActors,
		driftDetector:                     driftDetector,
		terraformSecurity:                 tfSec,
	}
	err := job.Authorize(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, job)

	return &JobDependenciesMock{
		costEstimator:                     costEstimator,
		vcs:                               vcs,
		terraformWorkspace:                terraformWorkspace,
		terraformerExecutor:               terraformerExecutor,
		terraformImportMigrationGenerator: terraformImportMigrationGenerator,
		resourcesCalculator:               resourcesCalculator,
		resourcesWriter:                   resourcesWriter,
		dragonDrop:                        dragonDrop,
		identifyCloudActors:               identifyCloudActors,
		driftDetector:                     driftDetector,
		terraformSecurity:                 tfSec,
	}, job
}

func TestRunJob_Success(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	// When
	mocks.dragonDrop.On("PutJobPullRequestURL", ctx, "").Return(nil)
	mocks.dragonDrop.On("InformComplete", ctx).Return(nil)
	mocks.dragonDrop.On("InformRepositoryCloned", ctx).Return(nil)
	mocks.dragonDrop.On("InformCloudActorIdentification", ctx).Return(nil)
	mocks.dragonDrop.On("InformCostEstimation", ctx).Return(nil)
	mocks.dragonDrop.On("InformSecurityScan", ctx).Return(nil)

	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute", ctx).Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(true, nil)
	mocks.terraformSecurity.On("ExecuteScan", ctx).Return(nil)

	err := job.Run(ctx)

	// Then
	assert.Nil(t, err)
	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 1)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 1)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 1)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 1)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 1)
}

func TestRunJob_CannotCloneRepo(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()

	vcsCloneError := errors.New("cannot clone repo")

	// When
	mocks.vcs.On("Clone").Return(vcsCloneError)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute", ctx).Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)
	mocks.dragonDrop.On("InformComplete", ctx).Return(nil)

	err := job.Run(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, vcsCloneError, errors.Unwrap(err))

	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 0)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 0)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 0)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 0)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 0)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 0)
}

func TestRunJob_CannotDownloadWorkspaceState(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	downloadWorkspaceErr := errors.New("cannot download workspace state")

	// When
	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(downloadWorkspaceErr)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute", ctx).Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)
	mocks.dragonDrop.On("InformComplete", ctx).Return(nil)
	mocks.dragonDrop.On("InformRepositoryCloned", ctx).Return(nil)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(true, nil)

	err := job.Run(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, downloadWorkspaceErr, errors.Unwrap(err))

	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 1)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 0)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 0)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 0)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 0)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 0)
}

func TestRunJob_CannotExecuteTerraformerExecutor(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	setUpTerraformerExecutorErr := errors.New("cannot set up terraformer executor")

	// When
	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(setUpTerraformerExecutorErr)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute", ctx).Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)
	mocks.dragonDrop.On("InformComplete", ctx).Return(nil)

	err := job.Run(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, setUpTerraformerExecutorErr, errors.Unwrap(err))

	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 1)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 0)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 0)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 0)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 0)
}

func TestRunJob_CannotExecuteTerraformImportMigrationGenerator(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	terraformImportMigrationGeneratorErr := errors.New("cannot execute terraform import")

	// When
	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(terraformImportMigrationGeneratorErr)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute", ctx).Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)
	mocks.dragonDrop.On("InformComplete", ctx).Return(nil)

	err := job.Run(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, terraformImportMigrationGeneratorErr, errors.Unwrap(err))

	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 1)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 0)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 0)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 0)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 0)
}

func TestRunJob_CannotCalculateResources(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	calculateResourcesErr := errors.New("cannot calculate resources")

	// When
	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(calculateResourcesErr)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute", ctx).Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)
	mocks.dragonDrop.On("InformComplete", ctx).Return(nil)

	err := job.Run(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, calculateResourcesErr, errors.Unwrap(err))

	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 1)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 0)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 0)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 0)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 0)
}

func TestRunJob_CannotDriftDetect(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	managedDriftDetectErr := errors.New("cannot do managed drift detection")

	// When
	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(false, managedDriftDetectErr)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute", ctx).Return(nil)
	mocks.resourcesWriter.On("Execute", ctx).Return("", nil)
	mocks.dragonDrop.On("InformComplete", ctx).Return(nil)

	err := job.Run(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, errors.Unwrap(err), managedDriftDetectErr)

	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 1)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.driftDetector.AssertNumberOfCalls(t, "Execute", 1)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 0)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 0)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 0)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 0)
}

func TestRunJob_CannotIdentifyCloudActors(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	identifyCloudActorsErr := errors.New("cannot identify cloud actors")

	// When
	mocks.dragonDrop.On("InformCloudActorIdentification", ctx).Return(nil)
	mocks.dragonDrop.On("InformCostEstimation", ctx).Return(nil)
	mocks.dragonDrop.On("InformSecurityScan", ctx).Return(nil)

	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(true, nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(identifyCloudActorsErr)
	mocks.costEstimator.On("Execute", ctx).Return(nil)
	mocks.resourcesWriter.On("Execute", ctx).Return("", nil)
	mocks.dragonDrop.On("InformComplete", ctx).Return(nil)

	err := job.Run(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, identifyCloudActorsErr, errors.Unwrap(err))

	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 1)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.driftDetector.AssertNumberOfCalls(t, "Execute", 1)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 1)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 0)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 0)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 0)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 0)
}

func TestRunJob_CannotCostEstimate(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	costEstimationErr := errors.New("cannot cost estimate")

	// When
	mocks.dragonDrop.On("InformCloudActorIdentification", ctx).Return(nil)
	mocks.dragonDrop.On("InformCostEstimation", ctx).Return(nil)
	mocks.dragonDrop.On("InformSecurityScan", ctx).Return(nil)

	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(true, nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute", ctx).Return(costEstimationErr)
	mocks.resourcesWriter.On("Execute", ctx).Return("", nil)
	mocks.dragonDrop.On("InformComplete", ctx).Return(nil)

	err := job.Run(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, costEstimationErr, errors.Unwrap(err))

	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 1)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.driftDetector.AssertNumberOfCalls(t, "Execute", 1)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 1)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 0)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 0)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 0)
}

func TestRunJob_CannotSecurityScan(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	securityScanErr := errors.New("cannot security scan")

	// When
	mocks.dragonDrop.On("InformCloudActorIdentification", ctx).Return(nil)
	mocks.dragonDrop.On("InformCostEstimation", ctx).Return(nil)
	mocks.dragonDrop.On("InformSecurityScan", ctx).Return(nil)

	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(true, nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute", ctx).Return(nil)
	mocks.terraformSecurity.On("ExecuteScan", ctx).Return(securityScanErr)
	mocks.resourcesWriter.On("Execute", ctx).Return("", nil)
	mocks.dragonDrop.On("InformComplete", ctx).Return(nil)

	err := job.Run(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, securityScanErr, errors.Unwrap(err))

	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 1)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.driftDetector.AssertNumberOfCalls(t, "Execute", 1)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 1)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 1)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 0)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 0)
}

func TestRunJob_CannotWriteResourcesOnVCS(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	writeResourcesErr := errors.New("cannot write resources on vcs")

	// When
	mocks.dragonDrop.On("InformCloudActorIdentification", ctx).Return(nil)
	mocks.dragonDrop.On("InformCostEstimation", ctx).Return(nil)
	mocks.dragonDrop.On("InformSecurityScan", ctx).Return(nil)

	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(true, nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute", ctx).Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", writeResourcesErr)
	mocks.dragonDrop.On("InformComplete", ctx).Return(nil)
	mocks.driftDetector.On("Execute", ctx).Return(true, nil)
	mocks.terraformSecurity.On("ExecuteScan", ctx).Return(nil)

	err := job.Run(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, writeResourcesErr, errors.Unwrap(err))

	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 1)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.driftDetector.AssertNumberOfCalls(t, "Execute", 1)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 1)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 1)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 0)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 1)
}

func TestRunJob_CannotInformCompleteStatus(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	informCompleteErr := errors.New("cannot inform incomplete status")

	// When
	mocks.dragonDrop.On("InformCloudActorIdentification", ctx).Return(nil)
	mocks.dragonDrop.On("InformCostEstimation", ctx).Return(nil)
	mocks.dragonDrop.On("InformSecurityScan", ctx).Return(nil)

	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute", ctx).Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)
	mocks.dragonDrop.On("PutJobPullRequestURL", ctx, "").Return(nil)
	mocks.dragonDrop.On("InformComplete", ctx).Return(informCompleteErr)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(true, nil)
	mocks.terraformSecurity.On("ExecuteScan", ctx).Return(nil)

	err := job.Run(ctx)

	// Then
	assert.NotNil(t, err)
	assert.ErrorIs(t, informCompleteErr, errors.Unwrap(err))

	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 1)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 1)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 1)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 1)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 1)
}

func TestRunJob_NotFoundNewResources_ButFoundManagedDriftedResources(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	innerResources := fmt.Errorf("[calculate_resource_to_workspace][error identifying new resources]%w", resourcesCalculator.ErrNoNewResources)
	calculateResourcesErr := fmt.Errorf("[resources_calculator][error calculating resources to workspace]%w", innerResources)

	// When
	mocks.dragonDrop.On("InformCloudActorIdentification", ctx).Return(nil)
	mocks.dragonDrop.On("InformCostEstimation", ctx).Return(nil)
	mocks.dragonDrop.On("InformSecurityScan", ctx).Return(nil)

	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(calculateResourcesErr)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute", ctx).Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)
	mocks.dragonDrop.On("PutJobPullRequestURL", ctx, "").Return(nil)
	mocks.dragonDrop.On("InformComplete", ctx).Return(nil)
	mocks.dragonDrop.On("InformRepositoryCloned", ctx).Return(nil)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(true, nil)
	mocks.terraformSecurity.On("ExecuteScan", ctx).Return(nil)

	err := job.Run(ctx)

	// Then
	assert.Nil(t, err)
	mocks.vcs.AssertNumberOfCalls(t, "Clone", 1)
	mocks.terraformWorkspace.AssertNumberOfCalls(t, "DownloadWorkspaceState", 1)
	mocks.terraformerExecutor.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformImportMigrationGenerator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesCalculator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.driftDetector.AssertNumberOfCalls(t, "Execute", 1)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 1)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 1)
	mocks.dragonDrop.AssertNumberOfCalls(t, "InformComplete", 1)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 1)
}

func Test_getProviderByCredential(t *testing.T) {
	type args struct {
		credential terraformValueObjects.Credential
	}
	tests := []struct {
		name    string
		args    args
		want    terraformValueObjects.Provider
		wantErr bool
	}{
		{
			name: "azurerm provider",
			args: args{
				credential: terraformValueObjects.Credential(
					`{"client_id": "123", "client_secret": "secret", "tenant_id": "tenant", "subscription_id": "subscription1"}`,
				),
			},
			want:    "azurerm",
			wantErr: false,
		},
		{
			name: "aws provider",
			args: args{
				credential: terraformValueObjects.Credential(
					`{"awsAccessKeyID": "AWS123", "awsSecretAccessKey": "DUGFVGBHAJ213"}`,
				),
			},
			want:    "aws",
			wantErr: false,
		},
		{
			name: "google provider",
			args: args{
				credential: terraformValueObjects.Credential(
					`{  "type": "service_account", "project_id": "project", "private_key_id": "123", "private_key": "key", 
						"client_email": "example@dragondrop.cloud", "client_id": "123456", "auth_uri": "https://accounts.google.com/o/oauth2/auth",
						"token_uri": "https://oauth2.googleapis.com/token", "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
						"client_x509_cert_url": "https://localhost.com"
					}`,
				),
			},
			want:    "google",
			wantErr: false,
		},
		{
			name: "error inferring provider",
			args: args{
				credential: terraformValueObjects.Credential(""),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty aws credentials",
			args: args{
				credential: terraformValueObjects.Credential(
					`{"awsAccessKeyID": "", "awsSecretAccessKey": ""}`,
				),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty azurerm credentials",
			args: args{
				credential: terraformValueObjects.Credential(
					`{"client_id": "", "client_secret": "", "tenant_id": "", "subscription_id": ""}`,
				),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "google provider",
			args: args{
				credential: terraformValueObjects.Credential(
					`{  "type": "", "project_id": "", "private_key_id": "", "private_key": "", 
						"client_email": "", "client_id": "", "auth_uri": "",
						"token_uri": "", "auth_provider_x509_cert_url": "",
						"client_x509_cert_url": ""
					}`,
				),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "json parsing error",
			args: args{
				credential: terraformValueObjects.Credential(
					`{error]`,
				),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getProviderByCredential(tt.args.credential)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equalf(t, tt.want, got, "getProviderByCredential(%v)", tt.args.credential)
		})
	}
}

func Test_getInferredData(t *testing.T) {
	type args struct {
		config JobConfig
	}
	tests := []struct {
		name    string
		args    args
		want    InferredData
		wantErr bool
	}{
		{
			name: "single provider aws",
			args: args{config: JobConfig{
				IsManagedDriftOnly: false,
				DivisionCloudCredentials: map[terraformValueObjects.Division]terraformValueObjects.Credential{
					terraformValueObjects.Division("division-1"): terraformValueObjects.Credential(
						`{"awsAccessKeyID": "AWS123", "awsSecretAccessKey": "DUGFVGBHAJ213"}`,
					),
				},
			}},
			want: InferredData{
				DivisionToProvider: map[terraformValueObjects.Division]terraformValueObjects.Provider{
					"division-1": "aws",
				},
			},
			wantErr: false,
		},
		{
			name: "two providers aws and azurerm",
			args: args{config: JobConfig{
				IsManagedDriftOnly: false,
				DivisionCloudCredentials: map[terraformValueObjects.Division]terraformValueObjects.Credential{
					terraformValueObjects.Division("division-1"): terraformValueObjects.Credential(
						`{"awsAccessKeyID": "AWS123", "awsSecretAccessKey": "DUGFVGBHAJ213"}`,
					),
					terraformValueObjects.Division("division-2"): terraformValueObjects.Credential(
						`{"client_id": "123", "client_secret": "secret", "tenant_id": "tenant", "subscription_id": "subscription1"}`,
					),
				},
			}},
			want: InferredData{
				DivisionToProvider: map[terraformValueObjects.Division]terraformValueObjects.Provider{
					"division-1": "aws",
					"division-2": "azurerm",
				},
			},
			wantErr: false,
		},
		{
			name: "three providers aws, azurerm and google",
			args: args{config: JobConfig{
				IsManagedDriftOnly: false,
				DivisionCloudCredentials: map[terraformValueObjects.Division]terraformValueObjects.Credential{
					terraformValueObjects.Division("division-1"): terraformValueObjects.Credential(
						`{"awsAccessKeyID": "AWS123", "awsSecretAccessKey": "DUGFVGBHAJ213"}`,
					),
					terraformValueObjects.Division("division-2"): terraformValueObjects.Credential(
						`{"client_id": "123", "client_secret": "secret", "tenant_id": "tenant", "subscription_id": "subscription1"}`,
					),
					terraformValueObjects.Division("division-3"): terraformValueObjects.Credential(
						`{  "type": "service_account", "project_id": "project", "private_key_id": "123", "private_key": "key", 
							"client_email": "example@dragondrop.cloud", "client_id": "123456", "auth_uri": "https://accounts.google.com/o/oauth2/auth",
							"token_uri": "https://oauth2.googleapis.com/token", "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
							"client_x509_cert_url": "https://localhost.com"
						}`,
					),
				},
			}},
			want: InferredData{
				DivisionToProvider: map[terraformValueObjects.Division]terraformValueObjects.Provider{
					"division-1": "aws",
					"division-2": "azurerm",
					"division-3": "google",
				},
			},
			wantErr: false,
		},
		{
			name: "two repeated divisions to provider aws",
			args: args{config: JobConfig{
				IsManagedDriftOnly: false,
				DivisionCloudCredentials: map[terraformValueObjects.Division]terraformValueObjects.Credential{
					terraformValueObjects.Division("division-1"): terraformValueObjects.Credential(
						`{"awsAccessKeyID": "AWS123", "awsSecretAccessKey": "DUGFVGBHAJ213"}`,
					),
					terraformValueObjects.Division("division-2"): terraformValueObjects.Credential(
						`{"awsAccessKeyID": "AWS123", "awsSecretAccessKey": "DUGFVGBHAJ213"}`,
					),
				},
			}},
			want: InferredData{
				DivisionToProvider: map[terraformValueObjects.Division]terraformValueObjects.Provider{
					"division-1": "aws",
					"division-2": "aws",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getInferredData(tt.args.config)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equalf(t, tt.want, got, "getInferredData(%v)", tt.args.config)
		})
	}
}
