package dragondrop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

// AuthorizeManagedJobRequestBody is a struct for sending a request to the dragondrop API to
// authorize a managed job.
type AuthorizeManagedJobRequestBody struct {
	JobID string
}

// PutPRURLRequest is a struct for the data in a request to the dragondrop API to
// update the PRURL for a given JobID.
type PutPRURLRequest struct {
	JobID string
	PRURL string
}

// JobStatusPostBody is a struct for sending a job status update to the dragondrop API.
type JobStatusPostBody struct {
	JobID  string
	Status string
}

// PostLogRequestBody is a struct for sending a log update to the dragondrop API.
type PostLogRequestBody struct {
	Log     string
	IsAlert bool
}

// CheckLoggerAndToken Check that logging is successful, and implicitly check that the OrgToken env variable is valid.
func (c *HTTPDragonDropClient) CheckLoggerAndToken(ctx context.Context) error {
	return c.postLog(ctx, "Successfully started job", true)
}

// InformStarted Informs to DragonDropAPI when job is Starting
func (c *HTTPDragonDropClient) InformStarted(ctx context.Context) error {
	return c.postJobStatus(ctx, "Started Job")
}

// InformComplete Informs to DragonDropAPI when job is Complete
func (c *HTTPDragonDropClient) InformComplete(ctx context.Context) error {
	return c.postJobStatus(ctx, "Complete")
}

// InformRepositoryCloned Informs to DragonDropAPI when job cloned the repository
func (c *HTTPDragonDropClient) InformRepositoryCloned(ctx context.Context) error {
	return c.postJobStatus(ctx, "Pulled Repository from VCS")
}

// InformCloudEnvironmentScanned Informs to DragonDropAPI when job scanned the cloud environment
func (c *HTTPDragonDropClient) InformCloudEnvironmentScanned(ctx context.Context) error {
	return c.postJobStatus(ctx, "Scanned cloud environment")
}

// InformCloudActorIdentification Informs to DragonDropAPI when job is identifying the cloud actors
func (c *HTTPDragonDropClient) InformCloudActorIdentification(ctx context.Context) error {
	return c.postJobStatus(ctx, "Identifying Cloud Actors")
}

// InformCostEstimation Informs to DragonDropAPI when job is estimating costs
func (c *HTTPDragonDropClient) InformCostEstimation(ctx context.Context) error {
	return c.postJobStatus(ctx, "Estimating Costs")
}

// InformSecurityScan Informs to DragonDropAPI when job is assessing cloud security
func (c *HTTPDragonDropClient) InformSecurityScan(ctx context.Context) error {
	return c.postJobStatus(ctx, "Assessing Cloud Security")
}

// InformCloudResourcesMappedToStateFile Informs to DragonDropAPI when job has mapped the resources to state file
func (c *HTTPDragonDropClient) InformCloudResourcesMappedToStateFile(ctx context.Context) error {
	return c.postJobStatus(ctx, "Map of cloud resources to state file")
}

// InformNoResourcesFound Informs to DragonDropAPI when job is InformNoResourcesFound
func (c *HTTPDragonDropClient) InformNoResourcesFound(ctx context.Context) error {
	return c.postJobStatus(ctx, "No Resources Found")
}

// postJobStatus sends a job status to the dragondrop API.
func (c *HTTPDragonDropClient) postJobStatus(ctx context.Context, status string) error {
	if c.config.JobID == "empty" || c.config.JobID == "" {
		return nil
	}

	logrus.Debugf("Post job status: %s", status)

	// Building Post Request Body
	jsonBody, err := json.Marshal(&JobStatusPostBody{
		JobID:  c.config.JobID,
		Status: status,
	})

	if err != nil {
		return fmt.Errorf("[post_job_status][error in json marshal]%v", err)
	}

	request, err := c.newRequest(
		ctx,
		"PostJobStatus",
		"POST",
		fmt.Sprintf("%v/job/status/", c.config.APIPath),
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return fmt.Errorf("[post_job_status][error in newRequest]%w", err)
	}

	response, err := c.httpClient.Do(request)

	if err != nil {
		return fmt.Errorf("[post_job_status][error in http POST request] %w", err)
	}

	defer response.Body.Close()
	if response.StatusCode != 201 {
		return fmt.Errorf("[post_job_status][was unsuccessful, with the server returning: %d]", response.StatusCode)
	}

	return nil
}

// PostLogAlert sends alert log to the dragondrop API.
func (c *HTTPDragonDropClient) PostLogAlert(ctx context.Context, log string) {
	logrus.Debugf("Post log alert: %s", log)

	err := c.postLog(ctx, log, true)
	if err != nil {
		logrus.Errorf("[dragon_drop][cannot send log alert]%s", err.Error())
		return
	}
}

// PostLog sends log to the dragondrop API.
func (c *HTTPDragonDropClient) PostLog(ctx context.Context, log string) {
	logrus.Debugf("Post log: %s", log)

	err := c.postLog(ctx, log, false)
	if err != nil {
		logrus.Errorf("[dragon_drop][cannot send log]%s", err.Error())
		return
	}
}

// postLog is a private method that sends a log to the dragondrop API.
func (c *HTTPDragonDropClient) postLog(ctx context.Context, log string, isAlert bool) error {
	if c.config.JobID == "empty" || c.config.JobID == "" {
		return nil
	}

	// Building Post Log Request body
	jsonBody, err := json.Marshal(
		PostLogRequestBody{
			Log:     log,
			IsAlert: isAlert,
		})

	if err != nil {
		return fmt.Errorf("[post_log][error in json marshal]%w", err)
	}

	request, err := c.newRequest(
		ctx,
		"PostLog",
		"POST",
		fmt.Sprintf("%v/log/", c.config.APIPath),
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return fmt.Errorf("[PostLog][error in newRequest]%w", err)
	}

	response, err := c.httpClient.Do(request)

	if err != nil {
		return fmt.Errorf("[post_log][error in http POST request]%w", err)
	}

	defer response.Body.Close()
	if response.StatusCode != 201 {
		return fmt.Errorf("[post_log][was unsuccessful, with the server returning: %d]", response.StatusCode)
	}

	return nil
}

// PutJobPullRequestURL sends the job url to the dragondrop API
func (c *HTTPDragonDropClient) PutJobPullRequestURL(ctx context.Context, prURL string) error {
	if c.config.JobID == "empty" || c.config.JobID == "" {
		return nil
	}

	// Building Post Log Request body
	jsonBody, err := json.Marshal(
		PutPRURLRequest{
			JobID: c.config.JobID,
			PRURL: prURL,
		})

	if err != nil {
		return fmt.Errorf("[PostLog] error in json.Marshal: %v", err)
	}

	request, err := c.newRequest(
		ctx,
		"PutURL",
		"PUT",
		fmt.Sprintf("%v/job/pr-link/", c.config.APIPath),
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return fmt.Errorf("[PutJobPullRequestURL][error in newRequest: %v]", err)
	}

	response, err := c.httpClient.Do(request)

	if err != nil {
		return fmt.Errorf("[PutJobPullRequestURL][error in http PUT request: %v]", err)
	}

	defer response.Body.Close()
	if response.StatusCode != 201 {
		return fmt.Errorf("[PutJobPullRequestURL][ was unsuccessful, with the server returning: %v]", response.StatusCode)
	}
	return nil
}

// AuthorizeManagedJob check with DragonDropAPI for valid auth of the current job, for a job managed
// by dragondrop.
func (c *HTTPDragonDropClient) AuthorizeManagedJob(ctx context.Context) (string, string, error) {
	// Building authorization request body
	jsonBody, err := json.Marshal(
		AuthorizeManagedJobRequestBody{
			JobID: c.config.JobID,
		})

	if err != nil {
		return "", "", fmt.Errorf("[authorize_managed_job][error in json marshal]%v", err)
	}

	request, err := c.newRequest(
		ctx,
		"GetJobAuthorization",
		"PUT",
		fmt.Sprintf("%v/job/authorize/", c.config.APIPath),
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return "", "", fmt.Errorf("[authorize_managed_job][error in newRequest]%w", err)
	}

	response, err := c.httpClient.Do(request)

	if err != nil {
		return "", "", fmt.Errorf("[authorize_managed_job][error in http GET request]%w", err)
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", "", fmt.Errorf("[authorize_managed_job][was unsuccessful, with the server returning: %v]", response.StatusCode)
	}

	// Read in response body to bytes array.
	outputBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", fmt.Errorf("[authorize_job][error in reading response into bytes array]%w", err)
	}

	var jobAuthResponse struct {
		JobName           string
		InfracostAPIToken string
	}
	err = json.Unmarshal(outputBytes, &jobAuthResponse)
	if err != nil {
		return "", "", fmt.Errorf("[authorize_job][unable to unmarshal %v]", string(outputBytes))
	}

	return jobAuthResponse.InfracostAPIToken, jobAuthResponse.JobName, nil
}
