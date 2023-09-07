package dragondrop

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	terraformWorkspace "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_workspace"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// HTTPDragonDropClientConfig is configuration for the HTTPDragonDropClient struct that conforms
// to envconfig's format expectations.
type HTTPDragonDropClientConfig struct {
	// APIPath is the dragondrop api path to which requests are sent.
	APIPath string

	// JobID is the unique identification string for the current job run.
	JobID string

	// OrgToken is the token that authorizes access to the dragondrop API.
	OrgToken string

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

// newRequest creates a new http request.
func (c *HTTPDragonDropClient) newRequest(ctx context.Context, requestName string, method string, requestPath string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, method, requestPath, body)
	if err != nil {
		return nil, fmt.Errorf("[new_request][error in http request instantiation with name: %s, err: %v]", requestName, err)
	}

	request.Header = http.Header{
		"Authorization": {"Bearer " + c.config.OrgToken},
		"Content-Type":  {"application/json"},
	}

	return request, nil
}

func (c *HTTPDragonDropClient) getDeletedResourcesList() []interface{} {
	_, err := os.Stat("outputs/drift-resources-deleted.json")
	if os.IsNotExist(err) {
		return []interface{}{}
	}

	deletedResources, err := readOutputFileAsSlice("outputs/drift-resources-deleted.json")
	if err != nil {
		log.Errorf("[error reading drift-resources-deleted.json]%v", err)
		return []interface{}{}
	}

	return deletedResources
}
