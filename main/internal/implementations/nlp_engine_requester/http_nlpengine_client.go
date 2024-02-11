package nlpenginerequestor

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
	log "github.com/sirupsen/logrus"
)

// HTTPNLPEngineClientConfig is configuration for the HTTPNLPEngineClient struct that conforms
// to envconfig's format expectations.
type HTTPNLPEngineClientConfig struct {
	// VCSRepo is the full path of the repo containing a customer's infrastructure specification.
	// At the moment, must be a valid GitHub repository URL.
	VCSRepo string `required:"true"`

	// NLPEndpoint is the endpoint for the NLP service.
	NLPEndpoint string

	// WorkspaceDirectories is a slice of directories that contains terraform workspaces within the user repo.
	WorkspaceDirectories terraformWorkspace.WorkspaceDirectoriesDecoder `required:"true"`
}

// HTTPNLPEngineClient is a struct that implements the NLPEngine interface and makes
// HTTP calls to the dragondrop API.
type HTTPNLPEngineClient struct {
	// httpClient is a http client shared across all http requests within this package.
	httpClient http.Client

	// Configuration parameters
	config HTTPNLPEngineClientConfig
}

// NewHTTPNLPEngineClient creates a new instance of HTTPNLPEngineClient, which implements the NLPEngine interface.
func NewHTTPNLPEngineClient(httpNLPEngineClientConfig HTTPNLPEngineClientConfig) interfaces.NLPEngine {
	return &HTTPNLPEngineClient{config: httpNLPEngineClientConfig}
}

type NLPEnginePostBody struct {
	NewResourceToDoc string `json:"new_resource_docs"`
	WorkspaceToDoc   string `json:"workspace_docs"`
}

// PostNLPEngine posts a correctly formatted request to the NLP engine endpoint, receiving and then saving out
// data on the mapping of new resources to state files.
func (c *HTTPNLPEngineClient) PostNLPEngine(ctx context.Context) error {
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

	err = os.WriteFile("outputs/new-resources-to-workspace.json", body, 0o400)
	if err != nil {
		return fmt.Errorf("[error writing new-resources-to-workspace.json]%v", err)
	}

	return nil
}

// newRequest creates a new http request with the given context, request type, request path, and body.
func (c *HTTPNLPEngineClient) newRequest(ctx context.Context, requestType string, requestPath string, body *bytes.Buffer) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, requestType, requestPath, body)
	if err != nil {
		return nil, fmt.Errorf("[http.NewRequestWithContext]%v", err)
	}

	request.Header = http.Header{
		"Content-Type": {"application/json"},
	}

	return request, nil
}
