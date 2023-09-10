package dragondrop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// AuthorizeJob Checks with DragonDropAPI for valid auth of the current job, for an oss job
func (c *HTTPDragonDropClient) AuthorizeJob(ctx context.Context) error {
	request, err := c.newRequest(
		ctx,
		"GetJobAuthorization",
		"GET",
		fmt.Sprintf("%v/job/authorize/oss/", c.config.APIPath),
		nil,
	)

	if err != nil {
		return fmt.Errorf("[authorize_job][error in newRequest]%w", err)
	}

	response, err := c.httpClient.Do(request)

	if err != nil {
		return fmt.Errorf("[authorize_job] error in http GET request]%w", err)
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("[authorize_job][was unsuccessful, with the server returning: %v]", response.StatusCode)
	}

	return nil
}

type NLPEnginePostBody struct {
	NewResourceToDoc string `json:"new_resource_docs"`
	WorkspaceToDoc   string `json:"workspace_docs"`
	Token            string `json:"token"`
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
		Token:            c.config.OrgToken,
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
