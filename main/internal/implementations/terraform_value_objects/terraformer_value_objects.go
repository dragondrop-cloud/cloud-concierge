package terraformvalueobjects

import (
	"fmt"
	"strings"
)

var AwsRegions = map[string]bool{
	"us-east-2":      true,
	"us-east-1":      true,
	"us-west-1":      true,
	"us-west-2":      true,
	"af-south-1":     true,
	"ap-east-1":      true,
	"ap-south-2":     true,
	"ap-southeast-3": true,
	"ap-southeast-4": true,
	"ap-south-1":     true,
	"ap-northeast-3": true,
	"ap-northeast-2": true,
	"ap-southeast-1": true,
	"ap-southeast-2": true,
	"ap-northeast-1": true,
	"ca-central-1":   true,
	"eu-central-1":   true,
	"eu-west-1":      true,
	"eu-west-2":      true,
	"eu-south-1":     true,
	"eu-west-3":      true,
	"eu-south-2":     true,
	"eu-north-1":     true,
	"eu-central-2":   true,
	"me-south-1":     true,
	"me-central-1":   true,
	"sa-east-1":      true,
	"us-gov-east-1":  true,
	"us-gov-west-":   true,
}

var AzureRegions = map[string]bool{
	"australiacentral":   true,
	"australiacentral2":  true,
	"australiaeast":      true,
	"australiasoutheast": true,
	"brazilsouth":        true,
	"brazilsoutheast":    true,
	"canadacentral":      true,
	"canadaeast":         true,
	"centralindia":       true,
	"centralus":          true,
	"centraluseuap":      true,
	"eastasia":           true,
	"eastus":             true,
	"eastus2":            true,
	"eastus2euap":        true,
	"francecentral":      true,
	"francesouth":        true,
	"germanycentral":     true,
	"germanynorth":       true,
	"germanywestcentral": true,
	"japaneast":          true,
	"japanwest":          true,
	"koreacentral":       true,
	"koreasouth":         true,
	"northcentralus":     true,
	"northeurope":        true,
	"southafricanorth":   true,
	"southafricawest":    true,
	"southcentralus":     true,
	"southeastasia":      true,
	"southindia":         true,
	"switzerlandnorth":   true,
	"switzerlandwest":    true,
	"uksouth":            true,
	"ukwest":             true,
	"westcentralus":      true,
	"westeurope":         true,
	"westindia":          true,
	"westus":             true,
	"westus2":            true,
	"westus3":            true,
}

var GoogleRegions = map[string]bool{
	"us-west1":                true,
	"asia-south1":             true,
	"asia-south2":             true,
	"asia-east1":              true,
	"asia-east2":              true,
	"asia-northeast1":         true,
	"asia-northeast2":         true,
	"asia-northeast3":         true,
	"asia-southeast1":         true,
	"australia-southeast1":    true,
	"australia-southeast2":    true,
	"europe-central2":         true,
	"europe-north2":           true,
	"europe-southwest1":       true,
	"europe-west1":            true,
	"europe-west2":            true,
	"europe-west3":            true,
	"europe-west4":            true,
	"europe-west6":            true,
	"europe-west8":            true,
	"europe-west9":            true,
	"northamerica-northeast1": true,
	"northamerica-northeast2": true,
	"southamerica-east1":      true,
	"us-central1":             true,
	"us-east1":                true,
	"us-east4":                true,
	"us-west2":                true,
	"us-west3":                true,
	"us-west4":                true,
}

type CloudRegion string

type CloudRegionsDecoder []CloudRegion

func (d *CloudRegionsDecoder) Decode(value string) error {
	azureAlreadySet := false
	awsAlreadySet := false
	googleAlreadySet := false
	regions := make([]CloudRegion, 0)
	if strings.Trim(value, " ") == "" || strings.Trim(value, " ") == "[]" {
		return nil
	}

	if value[0] != '[' || value[len(value)-1] != ']' {
		return fmt.Errorf("the value %v is not a list", value)
	}

	value = value[1 : len(value)-1]
	regionsList := strings.Split(value, ",")

	for _, region := range regionsList {
		region = strings.Trim(region, " ")
		if region[0] != '"' || region[len(region)-1] != '"' {
			return fmt.Errorf("the value %v is not a string", value)
		}
		region = region[1 : len(region)-1]

		if strings.Contains(region, `"`) {
			return fmt.Errorf("the value %v contains quotation marks in the middle", region)
		}

		if AzureRegions[region] {
			err := addProviderRegionIfValid(region, &regions, "azure", &azureAlreadySet)
			if err != nil {
				return err
			}
		} else if AwsRegions[region] {
			err := addProviderRegionIfValid(region, &regions, "aws", &awsAlreadySet)
			if err != nil {
				return err
			}
		} else if GoogleRegions[region] {
			err := addProviderRegionIfValid(region, &regions, "google", &googleAlreadySet)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("the region %v is not valid", region)
		}
	}

	*d = regions
	return nil
}

func addProviderRegionIfValid(region string, regions *[]CloudRegion, provider string, providerRegionAlreadySet *bool) error {
	if *providerRegionAlreadySet {
		return fmt.Errorf("the region is already set for the provider %v, only one region allowed per provider", provider)
	}
	*providerRegionAlreadySet = true
	*regions = append(*regions, CloudRegion(region))
	return nil
}
