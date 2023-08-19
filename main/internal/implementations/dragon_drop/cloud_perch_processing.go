package dragonDrop

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// ModulesVersions is a map of Terraform modules sources and their versions.
type ModulesVersions map[string]map[string]int

// CloudPerchData is a struct that contains all the data that is sent to DragonDrop.
type CloudPerchData struct {
	JobRunID string `json:"job_run_id"`
	ResourceInventoryData
	CloudCostsData
	TerraformFootprintData
	CloudSecurityData
}

// ResourceInventoryData is a struct that contains the number of resources that are and are not managed by Terraform.
type ResourceInventoryData struct {
	DriftedResources                 int `json:"drifted_resources"`
	ResourcesOutsideTerraformControl int `json:"resources_outside_terraform_control"`
}

// CloudCostsData is a struct that contains the costs of resources that are and are not managed by Terraform.
type CloudCostsData struct {
	CostsTerraformControlled float64 `json:"costs_terraform_controlled"`
	CostsOutsideOfTerraform  float64 `json:"costs_outside_of_terraform"`
}

// TerraformFootprintData is a struct that contains the versions of Terraform, Terraform Providers, and Terraform Modules.
type TerraformFootprintData struct {
	VersionsTFModules   string `json:"versions_tfmodules"`
	VersionsTFProviders string `json:"versions_tfproviders"`
	VersionsTF          string `json:"versions_tf"`
}

// CloudSecurityData is a struct that contains the number of resources with a given security risk.
type CloudSecurityData struct {
	SecurityRiskCritical int `json:"security_risk_critical"`
	SecurityRiskHigh     int `json:"security_risk_high"`
	SecurityRiskMedium   int `json:"security_risk_medium"`
	SecurityRiskLow      int `json:"security_risk_low"`
}

// formatResources formats the resources to be in the format of "resourceType.resourceTerraformerName"
func formatResources(resources map[string]interface{}) map[string]interface{} {
	formattedResources := map[string]interface{}{}
	for _, value := range resources {
		resourceType := value.(map[string]interface{})["ResourceType"].(string)
		resourceTerraformerName := value.(map[string]interface{})["ResourceTerraformerName"].(string)

		formattedResources[fmt.Sprintf("%s.%s", resourceType, resourceTerraformerName)] = value
	}

	return formattedResources
}

// getResourceInventoryData returns the number of resources outside of terraform control and the number of drifted resources
func (c *HTTPDragonDropClient) getResourceInventoryData(ctx context.Context) (ResourceInventoryData, map[string]interface{}, error) {
	newResources, err := readOutputFileAsMap("new-resources.json")
	if err != nil {
		return ResourceInventoryData{}, map[string]interface{}{}, fmt.Errorf("[error getting new resources]%w", err)
	}

	driftedResources, err := readOutputFileAsSlice("drift-resources-differences.json")
	if err != nil {
		return ResourceInventoryData{}, map[string]interface{}{}, fmt.Errorf("[error getting drifted resources]%w", err)
	}

	return ResourceInventoryData{
		DriftedResources:                 getUniqueDriftedResourceCount(driftedResources),
		ResourcesOutsideTerraformControl: len(newResources),
	}, newResources, nil
}

// getUniqueDriftedResourceCount returns the number of unique drifted resources
func getUniqueDriftedResourceCount(jsonInput []interface{}) int {
	uniqueDriftedResources := map[string]bool{}
	for _, driftedResource := range jsonInput {
		driftedResourceMap := driftedResource.(map[string]interface{})
		uniqueDriftedResources[driftedResourceMap["InstanceID"].(string)] = true
	}
	return len(uniqueDriftedResources)
}

// getCloudCostsData returns the costs of the resources outside of terraform control and the costs of the resources
func (c *HTTPDragonDropClient) getCloudCostsData(ctx context.Context, newResources map[string]interface{}) (CloudCostsData, error) {
	costEstimation, err := readOutputFileAsSlice("cost-estimates.json")
	if err != nil {
		return CloudCostsData{}, err
	}

	if len(costEstimation) == 0 {
		return CloudCostsData{}, nil
	}

	cloudCostsData := CloudCostsData{}
	for _, costEstimate := range costEstimation {
		costEstimateMap := costEstimate.(map[string]interface{})
		_, ok := newResources[costEstimateMap["resource_name"].(string)]

		if ok {
			cloudCostsData.CostsOutsideOfTerraform = roundFloat(cloudCostsData.CostsOutsideOfTerraform + getMonthlyCost(costEstimateMap))
		} else {
			cloudCostsData.CostsTerraformControlled = roundFloat(cloudCostsData.CostsTerraformControlled + getMonthlyCost(costEstimateMap))
		}
	}

	return cloudCostsData, nil
}

// getMonthlyCost returns the monthly cost of a resource
func getMonthlyCost(costEstimateMap map[string]interface{}) float64 {
	monthlyCost := costEstimateMap["monthly_cost"].(string)
	if strings.Trim(monthlyCost, "") == "" || monthlyCost == "0" {
		price := costEstimateMap["price"].(string)
		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return 0
		}

		return priceFloat
	}

	monthlyCostFloat, err := strconv.ParseFloat(monthlyCost, 64)
	if err != nil {
		return 0
	}

	return monthlyCostFloat
}

// roundFloat rounds a float to 3 decimal places
func roundFloat(val float64) float64 {
	ratio := math.Pow(10, float64(3))
	return math.Round(val*ratio) / ratio
}

// getCloudSecurityData returns the number of security risks found in the security scan
func (c *HTTPDragonDropClient) getCloudSecurityData(ctx context.Context) (CloudSecurityData, error) {
	securityScan, err := readOutputFileAsMap("security-scan.json")
	if err != nil {
		return CloudSecurityData{}, err
	}

	cloudSecurityData := CloudSecurityData{}
	results := securityScan["results"].([]interface{})
	for _, result := range results {
		severity := result.(map[string]interface{})["severity"].(string)
		switch severity {
		case "CRITICAL":
			cloudSecurityData.SecurityRiskCritical++
		case "HIGH":
			cloudSecurityData.SecurityRiskHigh++
		case "MEDIUM":
			cloudSecurityData.SecurityRiskMedium++
		case "LOW":
			cloudSecurityData.SecurityRiskLow++
		}
	}

	return cloudSecurityData, nil
}

// getTerraformFootprintData returns terraform footprint data for all terraform files
func (c *HTTPDragonDropClient) getTerraformFootprintData(ctx context.Context) (TerraformFootprintData, error) {
	files := []string{"current_cloud/versions.tf", "current_cloud/main.tf"}
	files = append(files, getAllTFFiles(ctx, c.config.WorkspaceDirectories)...)
	terraformFootprintData := TerraformFootprintData{}
	versionsTFModules := ModulesVersions{}

	for _, filename := range files {
		mainFileContent, err := readFile(filename)
		if err != nil {
			continue
		}

		inputHCLFile, hclDiag := hclwrite.ParseConfig(
			mainFileContent,
			"placeholder.tf",
			hcl.Pos{Line: 0, Column: 0, Byte: 0},
		)
		if hclDiag.HasErrors() {
			continue
		}

		if terraformFootprintData.VersionsTF == "" {
			terraformVersion, err := getTerraformVersions(inputHCLFile)
			if err == nil {
				terraformFootprintData.VersionsTF = terraformVersion
			}
		}

		if terraformFootprintData.VersionsTFProviders == "" {
			providerVersions, err := getProviderVersions(inputHCLFile)
			if err == nil {
				terraformFootprintData.VersionsTFProviders = providerVersions
			}
		}

		modulesVersions, err := getModulesVersions(inputHCLFile)
		if err == nil {
			versionsTFModules = concatenateVersions(versionsTFModules, modulesVersions)
		}
	}

	versionsTFModulesJSON, err := json.Marshal(versionsTFModules)
	if err != nil {
		return TerraformFootprintData{}, err
	}

	terraformFootprintData.VersionsTFModules = string(versionsTFModulesJSON)
	return terraformFootprintData, nil
}

// concatenateVersions concatenates the versions of all modules used in the terraform files
func concatenateVersions(modules ModulesVersions, newVersions ModulesVersions) ModulesVersions {
	for source, versions := range newVersions {
		if _, ok := modules[source]; ok {
			for version, count := range newVersions[source] {
				modules[source][version] += count
			}
		} else {
			modules[source] = versions
		}
	}

	return modules
}

// getModulesVersions returns the versions of all modules used in the terraform files
func getModulesVersions(inputHCLFile *hclwrite.File) (ModulesVersions, error) {
	blocks := inputHCLFile.Body().Blocks()
	if len(blocks) == 0 {
		return ModulesVersions{}, errors.New("[error parsing module blocks]")
	}

	modules := make([]*hclwrite.Block, 0)
	for _, block := range blocks {
		if block.Type() == "module" {
			modules = append(modules, block)
		}
	}

	if len(modules) == 0 {
		return ModulesVersions{}, errors.New("[error parsing module blocks]")
	}

	versions := ModulesVersions{}
	for _, module := range modules {
		versionAttribute := module.Body().GetAttribute("version")
		if versionAttribute == nil {
			continue
		}

		versionBytes := versionAttribute.Expr().BuildTokens(nil).Bytes()
		version, err := getAttributeValue(versionBytes)
		if err != nil {
			continue
		}

		sourceAttribute := module.Body().GetAttribute("source")
		if sourceAttribute == nil {
			continue
		}

		sourceBytes := sourceAttribute.Expr().BuildTokens(nil).Bytes()
		source, err := getAttributeValue(sourceBytes)
		if err != nil {
			continue
		}

		if _, ok := versions[source]; ok {
			versions[source][version]++
		} else {
			versions[source] = map[string]int{version: 1}
		}
	}

	return versions, nil
}

// getAttributeValue returns the value of an attribute
func getAttributeValue(attribute []byte) (string, error) {
	re := regexp.MustCompile(`"(.*)"`)
	matches := re.FindStringSubmatch(string(attribute))
	if len(matches) < 2 {
		return "", errors.New("[error parsing attribute value]")
	}

	return matches[1], nil
}

// getProviderVersions gets the provider versions from the terraform block within the file passed as parameter
func getProviderVersions(inputHCLFile *hclwrite.File) (string, error) {
	terraform := inputHCLFile.Body().FirstMatchingBlock("terraform", nil)
	if terraform == nil {
		return "", errors.New("[error parsing terraform]")
	}
	requiredProviders := terraform.Body().FirstMatchingBlock("required_providers", nil)
	if requiredProviders == nil {
		return "", errors.New("[error parsing terraform required_providers]")
	}
	providerAttribute, err := getProviderAttribute(requiredProviders)
	if err != nil {
		return "", nil
	}

	providerVersionValue := getVersionFromProviderAttribute(providerAttribute)
	providerSourceValue := getSourceFromProviderAttribute(providerAttribute)
	return fmt.Sprintf(`{"%s":{"%s":1}}`, providerSourceValue, providerVersionValue), nil
}

// getTerraformVersions gets the terraform version from the terraform block within the file passed as parameter
func getTerraformVersions(inputHCLFile *hclwrite.File) (string, error) {
	terraform := inputHCLFile.Body().FirstMatchingBlock("terraform", nil)
	if terraform == nil {
		return "", errors.New("[error parsing terraform]")
	}
	terraformVersionAttribute := terraform.Body().GetAttribute("required_version")
	if terraformVersionAttribute == nil {
		return "", errors.New("[error parsing terraform required_version]")
	}
	terraformVersionValue := string(terraformVersionAttribute.BuildTokens(nil).Bytes())

	re := regexp.MustCompile(`"(.*)"`)
	matches := re.FindStringSubmatch(terraformVersionValue)
	if len(matches) < 2 {
		return "", errors.New("[error parsing terraform version value]")
	}

	return fmt.Sprintf(`{"%s":1}`, matches[1]), nil
}

// getProviderAttribute gets the provider attribute from the required_providers block searching for aws, azurerm or google
func getProviderAttribute(requiredProviders *hclwrite.Block) (*hclwrite.Attribute, error) {
	awsProvider := requiredProviders.Body().GetAttribute("aws")
	if awsProvider != nil {
		return awsProvider, nil
	}

	azurermProvider := requiredProviders.Body().GetAttribute("azurerm")
	if azurermProvider != nil {
		return azurermProvider, nil
	}

	googleProvider := requiredProviders.Body().GetAttribute("google")
	if googleProvider != nil {
		return googleProvider, nil
	}

	return nil, errors.New("[provider not supported]")
}

// getSourceFromProviderAttribute gets the source from the provider attribute passed as a parameter
func getSourceFromProviderAttribute(attribute *hclwrite.Attribute) string {
	if attribute == nil {
		return ""
	}

	tokens := attribute.BuildTokens(nil)
	sourceIndex := -1
	for i, token := range tokens {
		if string(token.Bytes) == "source" {
			sourceIndex = i
			break
		}
	}

	sourceValueBytes := make([]byte, 0)
	beginToGetToken := false
	for _, token := range tokens[sourceIndex:] {
		if beginToGetToken {
			sourceValueBytes = append(sourceValueBytes, token.Bytes...)
			break
		}
		if string(token.Bytes) == "\"" {
			beginToGetToken = true
		}
	}

	return string(sourceValueBytes)
}

// getVersionFromProviderAttribute gets the version from the provider attribute passed as a parameter
func getVersionFromProviderAttribute(attribute *hclwrite.Attribute) string {
	if attribute == nil {
		return ""
	}

	tokens := attribute.BuildTokens(nil)
	versionIndex := -1
	for i, token := range tokens {
		if string(token.Bytes) == "version" {
			versionIndex = i
			break
		}
	}

	versionValueBytes := make([]byte, 0)
	beginToGetToken := false
	for _, token := range tokens[versionIndex:] {
		if beginToGetToken {
			versionValueBytes = append(versionValueBytes, token.Bytes...)
			break
		}
		if string(token.Bytes) == "\"" {
			beginToGetToken = true
		}
	}

	versionValue := string(versionValueBytes)
	versionValue = strings.Trim(versionValue, "~>=")
	return versionValue
}
