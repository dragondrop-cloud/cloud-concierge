package resourceswriter

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/hclcreate"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/markdowncreation"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// TerraformResourceWriter is a struct that implements the ResourcesWriter interface for usage within
// Jobs.
type TerraformResourceWriter struct {
	// hclCreate needed to manage terraform resources
	hclCreate hclcreate.HCLCreate

	// vcs instance needed to manage pull request changes
	vcs interfaces.VCS

	// markdownCreator is an instance of the MarkdownCreator struct that is used to create markdown
	markdownCreator *markdowncreation.MarkdownCreator

	// jobName is the name of the current job
	jobName string
}

// NewTerraformResourceWriter instantiates and returns a new instance of the TerraformResourceWriter.
func NewTerraformResourceWriter(hclCreate hclcreate.HCLCreate, vcs interfaces.VCS, markdownCreator *markdowncreation.MarkdownCreator, jobName string) interfaces.ResourcesWriter {
	return &TerraformResourceWriter{hclCreate: hclCreate, vcs: vcs, jobName: jobName, markdownCreator: markdownCreator}
}

// Execute writes new resources to the relevant version control system,
// and returns a pull request url corresponding to the new changes.
func (w *TerraformResourceWriter) Execute(ctx context.Context, createDummyFile bool, workspaceToDirectory map[string]string) (string, error) {
	logrus.Debugf("[terraform_resource_writer] Executing with jobName: %v, createDummyFile: %v, workspaceToDirectory: %v", w.jobName, createDummyFile, workspaceToDirectory)

	err := w.checkoutNewBranch(ctx)
	if err != nil {
		return "", fmt.Errorf("[terraform_resource_writer]%w", err)
	}

	err = w.writeNewResourcesAndMigrationStatements(ctx, createDummyFile, workspaceToDirectory)
	if err != nil {
		return "", fmt.Errorf("[terraform_resource_writer]%w", err)
	}

	err = w.writeNewMarkdownAnalysis()
	if err != nil {
		return "", fmt.Errorf("[terraform_resource_writer]%w", err)
	}

	prURL, err := w.commitChangesOpenPullRequest()
	if err != nil {
		return "", fmt.Errorf("[terraform_resource_writer]%w", err)
	}
	return prURL, nil
}

// commitChangesOpenPullRequest adds new files to the VCS, commits the changes,
// and opens a pull request for the branch.
func (w *TerraformResourceWriter) commitChangesOpenPullRequest() (string, error) {
	logrus.Debugf("[commit_changes_open_pull_request] Executing with jobName: %v", w.jobName)

	err := w.vcs.AddChanges()
	if err != nil {
		return "", fmt.Errorf("[commit_changes_open_pull_request][error in vcs.AddChanges]%w", err)
	}

	err = w.vcs.Commit()
	if err != nil {
		return "", fmt.Errorf("[commit_changes_open_pull_request][error in vcs.Commit]%w", err)
	}

	err = w.vcs.Push()
	if err != nil {
		return "", fmt.Errorf("[commit_changes_open_pull_request][error in vcs.Push]%w", err)
	}

	prURL, err := w.vcs.OpenPullRequest(w.jobName)
	if err != nil {
		return "", fmt.Errorf("[commit_changes_open_pull_request][error in vcs.OpenPullRequest]%w", err)
	}

	logrus.Debugf("[commit_changes_open_pull_request] prURL: %v", prURL)
	return prURL, nil
}

// writeNewMarkdownAnalysis writes out the markdown analysis of the identified resources which are currently outside
// of Terraform control.
func (w *TerraformResourceWriter) writeNewMarkdownAnalysis() error {
	err := w.markdownCreator.CreateMarkdownFile(w.jobName)
	if err != nil {
		return fmt.Errorf("[write_new_resources_and_migration_statements][error in pse.RunStateOfCloudReport]%w", err)
	}
	return nil
}

// writeNewResourcesAndMigrationStatements writes new resources and tfmigrate migration configuration to
// the customer's current code branch.
func (w *TerraformResourceWriter) writeNewResourcesAndMigrationStatements(ctx context.Context, createDummyFile bool, workspaceToDirectory map[string]string) error {
	logrus.Debugf("[write_new_resources_and_migration_statements] createDummyFile: %v, workspaceToDirectory: %v", createDummyFile, workspaceToDirectory)
	if createDummyFile {

		err := w.writeDummyFile(ctx, workspaceToDirectory)
		if err != nil {
			return fmt.Errorf("[write_new_resources_and_migration_statements][error in write_dummy_file]%w", err)
		}

		return nil
	}

	err := w.hclCreate.ExtractResourceDefinitions(createDummyFile, workspaceToDirectory)
	if err != nil {
		return fmt.Errorf("[write_new_resources_and_migration_statements][error in hclc.ExtractResourceDefinitions()]%w", err)
	}

	id, err := w.vcs.GetID()
	if err != nil {
		return fmt.Errorf("[write_new_resources_and_migration_statements][error getting the vcs id]%w", err)
	}

	err = w.hclCreate.CreateImports(id, workspaceToDirectory)
	if err != nil {
		return fmt.Errorf("[write_new_resources_and_migration_statements][error in hclc.CreateImports]%w", err)
	}

	return nil
}

// checkoutNewBranch checks out a new branch within the version control system
func (w *TerraformResourceWriter) checkoutNewBranch(ctx context.Context) error {
	logrus.Debugf("[terraform_resource_writer] Executing checkoutNewBranch with jobName: %v", w.jobName)

	err := w.vcs.Checkout(w.jobName)
	if err != nil {
		return fmt.Errorf("[checkout_new_branch][error in checkout to new branch with vcs]%w", err)
	}

	return nil
}

func (w *TerraformResourceWriter) writeDummyFile(_ context.Context, workspaceToDirectory map[string]string) error {
	for _, directory := range workspaceToDirectory {
		err := os.MkdirAll(fmt.Sprintf("repo%vcloud-concierge/placeholder", directory), 0o400)
		if err != nil {
			return fmt.Errorf("error creating placeholder folder %v: %v", directory, err)
		}

		newFilePath := fmt.Sprintf("repo%vcloud-concierge/placeholder/dragondrop_placeholder.txt", directory)

		err = os.WriteFile(newFilePath, []byte("Placeholder file for opening a PR"), 0o400)
		if err != nil {
			return fmt.Errorf("error writing the placeholder file %v", err)
		}

		err = os.WriteFile("outputs/new-resources-to-documents.json", []byte("{}"), 0o400)
		if err != nil {
			return fmt.Errorf("error writing new resources empty JSON file: %v", err)
		}
		break
	}

	return nil
}
