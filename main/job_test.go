package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	resourcesCalculator "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/resources_calculator"
	. "github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

type JobDependenciesMock struct {
	vcs                               *VCSMock
	terraformWorkspace                *TerraformWorkspaceMock
	terraformerExecutor               *TerraformerExecutorMock
	terraformImportMigrationGenerator *TerraformImportMigrationGeneratorMock
	resourcesCalculator               *ResourcesCalculatorMock
	resourcesWriter                   *ResourcesWriterMock
	identifyCloudActors               *IdentifyCloudActorsMock
	costEstimator                     *CostEstimationMock
	driftDetector                     *TerraformManagedResourcesDriftDetectorMock
	terraformSecurity                 *TerraformSecurityMock
}

func createValidJob(_ *testing.T) (*JobDependenciesMock, *Job) {
	vcs := new(VCSMock)
	terraformWorkspace := new(TerraformWorkspaceMock)
	terraformerExecutor := new(TerraformerExecutorMock)
	terraformImportMigrationGenerator := new(TerraformImportMigrationGeneratorMock)
	resourcesCalculator := new(ResourcesCalculatorMock)
	resourcesWriter := new(ResourcesWriterMock)
	identifyCloudActors := new(IdentifyCloudActorsMock)
	costEstimator := new(CostEstimationMock)
	driftDetector := new(TerraformManagedResourcesDriftDetectorMock)
	tfSec := new(TerraformSecurityMock)

	costEstimator.On("SetInfracostAPIToken", "xyz").Return()
	vcs.On("SetToken").Return()

	job := &Job{
		costEstimator:                     costEstimator,
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

	return &JobDependenciesMock{
		costEstimator:                     costEstimator,
		vcs:                               vcs,
		terraformWorkspace:                terraformWorkspace,
		terraformerExecutor:               terraformerExecutor,
		terraformImportMigrationGenerator: terraformImportMigrationGenerator,
		resourcesCalculator:               resourcesCalculator,
		resourcesWriter:                   resourcesWriter,
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

	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute").Return(nil)
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
	mocks.costEstimator.On("Execute").Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)

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
	mocks.costEstimator.On("Execute").Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)
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
	mocks.costEstimator.On("Execute").Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)

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
	mocks.costEstimator.On("Execute").Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)

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
	mocks.costEstimator.On("Execute").Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", nil)

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
	mocks.costEstimator.On("Execute").Return(nil)
	mocks.resourcesWriter.On("Execute", ctx).Return("", nil)

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
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 0)
}

func TestRunJob_CannotIdentifyCloudActors(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	identifyCloudActorsErr := errors.New("cannot identify cloud actors")

	// When
	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(true, nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(identifyCloudActorsErr)
	mocks.costEstimator.On("Execute").Return(nil)
	mocks.resourcesWriter.On("Execute", ctx).Return("", nil)

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
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 0)
}

func TestRunJob_CannotCostEstimate(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	costEstimationErr := errors.New("cannot cost estimate")

	// When
	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(true, nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute").Return(costEstimationErr)
	mocks.resourcesWriter.On("Execute", ctx).Return("", nil)

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
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 0)
}

func TestRunJob_CannotSecurityScan(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	securityScanErr := errors.New("cannot security scan")

	// When
	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(true, nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute").Return(nil)
	mocks.terraformSecurity.On("ExecuteScan", ctx).Return(securityScanErr)
	mocks.resourcesWriter.On("Execute", ctx).Return("", nil)

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
}

func TestRunJob_CannotWriteResourcesOnVCS(t *testing.T) {
	// Given
	mocks, job := createValidJob(t)
	ctx := context.Background()
	divisionToProvider := make(map[string]string)

	writeResourcesErr := errors.New("cannot write resources on vcs")

	// When
	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(nil)
	mocks.driftDetector.On("Execute", ctx, divisionToProvider).Return(true, nil)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute").Return(nil)
	mocks.resourcesWriter.On("Execute").Return("", writeResourcesErr)
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
	mocks.vcs.On("Clone").Return(nil)
	mocks.terraformWorkspace.On("FindTerraformWorkspaces", ctx).Return(divisionToProvider, nil)
	mocks.terraformWorkspace.On("DownloadWorkspaceState").Return(nil)
	mocks.terraformerExecutor.On("Execute").Return(nil)
	mocks.terraformImportMigrationGenerator.On("Execute").Return(nil)
	mocks.resourcesCalculator.On("Execute").Return(calculateResourcesErr)
	mocks.identifyCloudActors.On("Execute", ctx).Return(nil)
	mocks.costEstimator.On("Execute").Return(nil)
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
	mocks.driftDetector.AssertNumberOfCalls(t, "Execute", 1)
	mocks.identifyCloudActors.AssertNumberOfCalls(t, "Execute", 1)
	mocks.costEstimator.AssertNumberOfCalls(t, "Execute", 1)
	mocks.resourcesWriter.AssertNumberOfCalls(t, "Execute", 1)
	mocks.terraformSecurity.AssertNumberOfCalls(t, "ExecuteScan", 1)
}
