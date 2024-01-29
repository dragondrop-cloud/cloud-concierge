package dragondrop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	terraformWorkspace "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_workspace"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// HTTPDragonDropClientConfig is configuration for the HTTPDragonDropClient struct that conforms
// to envconfig's format expectations.
type HTTPDragonDropClientConfig struct {
	// VCSRepo is the full path of the repo containing a customer's infrastructure specification.
	// At the moment, must be a valid GitHub repository URL.
	VCSRepo string `required:"true"`

	// NLPEndpoint is the endpoint for the NLP service.
	NLPEndpoint string

	// WorkspaceDirectories is a slice of directories that contains terraform workspaces within the user repo.
	WorkspaceDirectories terraformWorkspace.WorkspaceDirectoriesDecoder `required:"true"`
}

// HTTPDragonDropClient is a struct that implements the DragonDrop interface and makes
// HTTP calls to the dragondrop API.
type HTTPDragonDropClient struct {
	// httpClient is a http client shared across all http requests within this package.
	httpClient http.Client

	// Configuration parameters
	config HTTPDragonDropClientConfig
}

// NewHTTPDragonDropClient creates a new instance of HTTPDragonDropClient, which implements the DragonDrop interface.
func NewHTTPDragonDropClient(httpDragonDropClientConfig HTTPDragonDropClientConfig) interfaces.DragonDrop {
	return &HTTPDragonDropClient{config: httpDragonDropClientConfig}
}

type NLPEnginePostBody struct {
	NewResourceToDoc string `json:"new_resource_docs"`
	WorkspaceToDoc   string `json:"workspace_docs"`
}

// PostNLPEngine posts a correctly formatted request to the NLP engine endpoint, receiving and then saving out
// data on the mapping of new resources to state files.
func (c *HTTPDragonDropClient) PostNLPEngine(ctx context.Context) error {
	newResourceToDocBytes, err := os.ReadFile("outputs/new-resources-to-documents.json")
	if err != nil {
		return fmt.Errorf("[post_nlp_engine][error reading new-resources-to-documents.json]%v", err)
	}
	workspaceToDocBytes, err := os.ReadFile("outputs/workspace-to-documents.json")
	if err != nil {
		return fmt.Errorf("[post_nlp_engine][error reading workspace-to-documents.json]%v", err)
	}

	jsonBody, err := json.Marshal(&NLPEnginePostBody{
		NewResourceToDoc: string(newResourceToDocBytes),
		WorkspaceToDoc:   string(workspaceToDocBytes),
	})

	if err != nil {
		return fmt.Errorf("[post_nlp_engine][error in json marshal]%v", err)
	}

	request, err := c.newRequest(
		ctx,
		"PostNLPEngine",
		"POST",
		fmt.Sprintf("%v", c.config.NLPEndpoint),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("[post_nlp_engine][error in newRequest]%w", err)
	}

	log.Info("Sending request to NLP engine...")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("[post_nlp_engine] error in http POST request]%w", err)
	}

	defer response.Body.Close()
	if response.StatusCode != 201 {
		return fmt.Errorf("[post_nlp_engine][was unsuccessful, with the server returning: %v]", response.StatusCode)
	}
	log.Info("NLP engine completed successfully.")

	// Read response body into a string
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("[error reading response body]%v", err)
	}

	err = os.WriteFile("outputs/new-resources-to-workspace.json", body, 0400)
	if err != nil {
		return fmt.Errorf("[error writing new-resources-to-workspace.json]%v", err)
	}

	return nil
}
