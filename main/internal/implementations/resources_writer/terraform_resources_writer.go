package resourcesWriter

import (
	"context"
	"fmt"
	"os"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/hclcreate"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/pyscriptexec"
)

// TerraformResourceWriter is a struct that implements the ResourcesWriter interface for usage within
// "live" dragondrop Jobs.
type TerraformResourceWriter struct {
	// hclCreate needed to manage terraform resources
	hclCreate hclcreate.HCLCreate

	// vcs instance needed to manage pull request changes
	vcs interfaces.VCS

	// pyScriptExec is the python script executor
	pyScriptExec pyscriptexec.PyScriptExec

	// jobName is the name of the current job
	jobName string

	// dragonDrop is an implementation of the DragonDrop interface
	dragonDrop interfaces.DragonDrop
}

// NewTerraformResourceWriter instantiates and returns a new instance of the TerraformResourceWriter.
func NewTerraformResourceWriter(hclCreate hclcreate.HCLCreate, vcs interfaces.VCS, pyScriptExec pyscriptexec.PyScriptExec, dragonDrop interfaces.DragonDrop) interfaces.ResourcesWriter {
	return &TerraformResourceWriter{hclCreate: hclCreate, vcs: vcs, pyScriptExec: pyScriptExec, dragonDrop: dragonDrop}
}

// Execute writes new resources to the relevant version control system,
// and returns a pull request url corresponding to the new changes.
func (w *TerraformResourceWriter) Execute(ctx context.Context, jobName string, createDummyFile bool, workspaceToDirectory map[string]string) (string, error) {
	w.jobName = jobName

	err := w.checkoutNewBranch(ctx)
	if err != nil {
		return "", fmt.Errorf("[terraform_resource_writer]%w", err)
	}

	err = w.writeNewResourcesAndMigrationStatements(ctx, createDummyFile, workspaceToDirectory)
	if err != nil {
		return "", fmt.Errorf("[terraform_resource_writer]%w", err)
	}

	err = w.writeNewMarkdownAnalysis(ctx)
	if err != nil {
		return "", fmt.Errorf("[terraform_resource_writer]%w", err)
	}

	prURL, err := w.commitChangesOpenPullRequest(ctx)
	if err != nil {
		return "", fmt.Errorf("[terraform_resource_writer]%w", err)
	}

	w.dragonDrop.PostLogAlert(ctx, fmt.Sprintf("Job is complete, pull request opened at URL: %v", prURL))
	return prURL, nil
}

// commitChangesOpenPullRequest adds new files to the VCS, commits the changes,
// and opens a pull request for the branch.
func (w *TerraformResourceWriter) commitChangesOpenPullRequest(ctx context.Context) (string, error) {
	w.dragonDrop.PostLog(ctx, "Beginning to add, commit, push and open a pull request for changes made.")

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

	w.dragonDrop.PostLog(ctx, "Done opening a pull request for changes made.")
	return prURL, nil
}

// writeNewMarkdownAnalysis writes out the markdown analysis of the identified resources which are currently outside
// of Terraform control.
func (w *TerraformResourceWriter) writeNewMarkdownAnalysis(ctx context.Context) error {
	w.dragonDrop.PostLog(ctx, "Beginning to generate and save markdown analysis.")

	id, err := w.vcs.GetID()
	if err != nil {
		return fmt.Errorf("[write_new_resources_and_migration_statements][error getting the vcs id]%w", err)
	}

	err = w.pyScriptExec.RunStateOfCloudReport(id, w.jobName)
	if err != nil {
		return fmt.Errorf("[write_new_resources_and_migration_statements][error in pse.RunStateOfCloudReport]%w", err)
	}

	w.dragonDrop.PostLog(ctx, "Done generating and saving markdown analysis.")
	return nil
}

// writeNewResourcesAndMigrationStatements writes new resources and tfmigrate migration configuration to
// the customer's current code branch.
func (w *TerraformResourceWriter) writeNewResourcesAndMigrationStatements(ctx context.Context, createDummyFile bool, workspaceToDirectory map[string]string) error {
	if createDummyFile {
		w.dragonDrop.PostLog(ctx, "Beginning to write dummy file because there are no new resources.")

		err := w.writeDummyFile(ctx, workspaceToDirectory)
		if err != nil {
			return fmt.Errorf("[write_new_resources_and_migration_statements][error in write_dummy_file]%w", err)
		}

		return nil
	}

	w.dragonDrop.PostLog(ctx, "Beginning to write new resources and migration statements.")

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

	w.dragonDrop.PostLog(ctx, "Done writing new resources and migration statements.")
	return nil
}

// checkoutNewBranch checks out a new branch within the version control system
func (w *TerraformResourceWriter) checkoutNewBranch(ctx context.Context) error {
	w.dragonDrop.PostLog(ctx, "Beginning to checkout new branch.")

	err := w.vcs.Checkout(w.jobName)

	if err != nil {
		return fmt.Errorf("[checkout_new_branch][error in checkout to new branch with vcs]%w", err)
	}

	w.dragonDrop.PostLog(ctx, "Done checkin-out a new branch.")

	return nil
}

func (w *TerraformResourceWriter) writeDummyFile(ctx context.Context, workspaceToDirectory map[string]string) error {
	for _, directory := range workspaceToDirectory {
		err := os.MkdirAll(fmt.Sprintf("repo%vdragondrop/placeholder", directory), 0400)
		if err != nil {
			return fmt.Errorf("error creating placeholder folder %v: %v", directory, err)
		}

		newFilePath := fmt.Sprintf("repo%vdragondrop/placeholder/dragondrop_placeholder.txt", directory)

		err = os.WriteFile(newFilePath, []byte("Placeholder file for opening a PR"), 0400)
		if err != nil {
			return fmt.Errorf("error writing the placeholder file %v", err)
		}

		err = os.WriteFile("mappings/new-resources-to-documents.json", []byte("{}"), 0400)
		if err != nil {
			return fmt.Errorf("error writing new resources empty JSON file: %v", err)
		}
		break
	}

	return nil
}
