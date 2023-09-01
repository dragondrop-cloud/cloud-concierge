package dragondrop

import (
	"context"
	"fmt"
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
