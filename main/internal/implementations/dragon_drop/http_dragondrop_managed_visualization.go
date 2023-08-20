package dragonDrop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// SendCloudPerchData sends CloudPerchData to DragonDrop.
func (c *HTTPDragonDropClient) SendCloudPerchData(ctx context.Context) error {
	if c.config.JobID == "empty" || c.config.JobID == "" {
		return nil
	}
	resourceInventoryData, newResources, err := c.getResourceInventoryData(ctx)
	if err != nil {
		return fmt.Errorf("[error getting ResourceInventoryData]%w", err)
	}

	cloudCostsData, err := c.getCloudCostsData(ctx, formatResources(newResources))
	if err != nil {
		return fmt.Errorf("[error getting CloudCostsData]%w", err)
	}

	cloudSecurityData, err := c.getCloudSecurityData(ctx)
	if err != nil {
		return fmt.Errorf("[error getting CloudSecurityData]%w", err)
	}

	terraformFootprintData, err := c.getTerraformFootprintData(ctx)
	if err != nil {
		return fmt.Errorf("[error getting TerraformFootprintData]%w", err)
	}

	// Only sending highly anonymized data to the DragonDrop API for managed cloud-concierge instances
	cloudPerchData := &CloudPerchData{
		JobRunID:               c.config.JobID,
		ResourceInventoryData:  resourceInventoryData,
		CloudCostsData:         cloudCostsData,
		CloudSecurityData:      cloudSecurityData,
		TerraformFootprintData: terraformFootprintData,
	}
	fmt.Printf("[cloud perch data] %+v\n", cloudPerchData)

	return c.sendRequest(ctx, cloudPerchData)
}

// sendRequest sends a request to the Dragon Drop API
func (c *HTTPDragonDropClient) sendRequest(ctx context.Context, cloudPerchData *CloudPerchData) error {
	jsonBody, err := json.Marshal(cloudPerchData)
	if err != nil {
		return fmt.Errorf("[send_cloud_perch_data][error in json.Marshal]%w", err)
	}

	request, err := c.newRequest(
		ctx,
		"SendCloudPerchData",
		"POST",
		fmt.Sprintf("%v/cloud-perch/", c.config.APIPath),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("[send_cloud_perch_data][error in newRequest]%w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("[send_cloud_perch_data][error in http POST request]%w", err)
	}

	defer response.Body.Close()
	if response.StatusCode != 201 {
		return fmt.Errorf("[send_cloud_perch_data][was unsuccessful, with the server returning: %v]", response.StatusCode)
	}

	return nil
}
