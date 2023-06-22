package dragonDrop

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/dragondrop-cloud/driftmitigation/interfaces"
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
