package dragondrop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// SendCloudPerchData sends CloudPerchData to DragonDrop.
func (c *HTTPDragonDropClient) SendCloudPerchData(ctx context.Context) error {
	if c.config.JobID == "empty" || c.config.JobID == "" {
		return nil
	}

	newResources, err := readOutputFileAsMap("new-resources.json")
	if err != nil {
		return fmt.Errorf("[error reading new-resources.json]%w", err)
	}
	logrus.Debugf("[dragon_drop][cloud_perch_processing] new resources: %+v", newResources)

	driftedResources, err := readOutputFileAsSlice("drift-resources-differences.json")
	if err != nil {
		return fmt.Errorf("[error reading drift-resources-differences.json]%w", err)
	}
	logrus.Debugf("[dragon_drop][cloud_perch_processing] drifted resources: %+v", driftedResources)

	costData, err := readOutputFileAsSlice("cost-estimates.json")
	if err != nil {
		return fmt.Errorf("[error reading cost-estimates.json]%w", err)
	}
	logrus.Debugf("[dragon_drop][cloud_perch_processing] cost data: %+v", costData)

	securityData, err := readOutputFileAsMap("security-scan.json")
	if err != nil {
		return fmt.Errorf("[error reading security-scan.json]%w", err)
	}
	logrus.Debugf("[dragon_drop][cloud_perch_processing] security data: %+v", securityData)

	cloudActorBytes, err := os.ReadFile("outputs/resources-to-cloud-actions.json")
	if err != nil {
		return fmt.Errorf("[error reading resources-to-cloud-actions.json]%w", err)
	}
	logrus.Debugf("[dragon_drop][cloud_perch_processing] cloud actor data: %+v", cloudActorBytes)

	deletedResources := c.getDeletedResourcesList()
	logrus.Debugf("[dragon_drop][cloud_perch_processing] deleted resources: %+v", deletedResources)

	resourceInventoryData, newResources, err := c.getResourceInventoryData(newResources, driftedResources, deletedResources)
	if err != nil {
		return fmt.Errorf("[error getting ResourceInventoryData]%w", err)
	}

	cloudCostsData, err := c.getCloudCostsData(ctx, newResources, costData)
	if err != nil {
		return fmt.Errorf("[error getting CloudCostsData]%w", err)
	}

	cloudSecurityData, err := c.getCloudSecurityData(ctx, securityData)
	if err != nil {
		return fmt.Errorf("[error getting CloudSecurityData]%w", err)
	}

	terraformFootprintData, err := c.getTerraformFootprintData(ctx)
	if err != nil {
		return fmt.Errorf("[error getting TerraformFootprintData]%w", err)
	}

	cloudActorData, err := c.getCloudActorData(ctx, cloudActorBytes)
	if err != nil {
		return fmt.Errorf("[error getting CloudActorData]%w", err)
	}

	// Only sending highly anonymized data to the DragonDrop API for managed cloud-concierge instances
	cloudPerchData := &CloudPerchData{
		JobRunID:               c.config.JobID,
		ResourceInventoryData:  resourceInventoryData,
		CloudActorData:         cloudActorData,
		CloudCostsData:         cloudCostsData,
		CloudSecurityData:      cloudSecurityData,
		TerraformFootprintData: terraformFootprintData,
	}
	logrus.Debugf("[cloud perch data] %+v\n", cloudPerchData)

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
