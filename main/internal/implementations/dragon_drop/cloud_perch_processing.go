package dragondrop

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// ModulesVersions is a map of Terraform modules sources and their versions.
type ModulesVersions map[string]map[string]int

// CloudPerchData is a struct that contains all the data that is sent to DragonDrop.
type CloudPerchData struct {
	JobRunID string `json:"job_run_id"`
	CloudActorData
	CloudCostsData
	CloudSecurityData
	ResourceInventoryData
	TerraformFootprintData
}

// ResourceInventoryData is a struct that contains the number of resources that are and are not managed by Terraform.
type ResourceInventoryData struct {
	DriftedResources                 int `json:"drifted_resources"`
	DeletedResources                 int `json:"deleted_resources"`
	ResourcesOutsideTerraformControl int `json:"resources_outside_terraform_control"`
}

// CloudActorData is a struct that contains the number of resources modified and created outside of Terraform control
// aggregated by cloud actor but in a string format, using marshalled ActorData list.
type CloudActorData struct {
	ActorsData string `json:"actors_data"`
}

// ActorData is a struct that contains the number of resources modified and created outside of Terraform control
// for a given cloud actor.
type ActorData struct {
	Actor    string `json:"actor_name"`
	Modified int    `json:"modified"`
	Created  int    `json:"created"`
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
func (c *HTTPDragonDropClient) getResourceInventoryData(newResources map[string]interface{}, driftedResources []interface{}, deletedResources []interface{}) (ResourceInventoryData, map[string]interface{}, error) {
	return ResourceInventoryData{
		DriftedResources:                 getUniqueDriftedResourceCount(driftedResources),
		DeletedResources:                 len(deletedResources),
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

// getCloudActorData returns the number of resources modified and created outside of Terraform control aggregated by cloud actor.
func (c *HTTPDragonDropClient) getCloudActorData(_ context.Context, cloudActorBytes []byte) (CloudActorData, error) {
	resourceToActions := &terraformValueObjects.ResourceActionMap{}
	err := json.Unmarshal(cloudActorBytes, resourceToActions)
	if err != nil {
		return CloudActorData{}, fmt.Errorf("failed to unmarshal cloud actor bytes: %w", err)
	}

	if (&terraformValueObjects.ResourceActionMap{}) == resourceToActions {
		return CloudActorData{}, nil
	}

	// Capturing and building data structure within a "helper" map before converting to a slice of ActorData
	// to enable ~O(1) lookup for each new resource action.
	actorToActorData := map[string]*ActorData{}
	for _, actions := range *resourceToActions {
		if actions.Creator != nil {
			actor := string(actions.Creator.Actor)
			if _, ok := actorToActorData[actor]; !ok {
				actorToActorData[actor] = &ActorData{
					Actor:    actor,
					Modified: 0,
					Created:  1,
				}
			} else {
				actorToActorData[actor].Created++
			}
		}
		if actions.Modifier != nil {
			actor := string(actions.Modifier.Actor)
			if _, ok := actorToActorData[actor]; !ok {
				actorToActorData[actor] = &ActorData{
					Actor:    actor,
					Modified: 1,
					Created:  0,
				}
			} else {
				actorToActorData[actor].Modified++
			}
		}
	}

	var actorsData []ActorData
	for _, actorData := range actorToActorData {
		actorsData = append(actorsData, *actorData)
	}

	marshalledActorsData, err := json.Marshal(actorsData)
	if err != nil {
		return CloudActorData{}, err
	}
	cloudActorData := CloudActorData{
		ActorsData: string(marshalledActorsData),
	}

	return cloudActorData, nil
}

// getCloudCostsData returns the costs of the resources outside of Terraform control and the costs of the resources
// already controlled by Terraform.
func (c *HTTPDragonDropClient) getCloudCostsData(_ context.Context, newResources map[string]interface{}, costEstimation []interface{}) (CloudCostsData, error) {
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
func (c *HTTPDragonDropClient) getCloudSecurityData(_ context.Context, securityScan map[string]interface{}) (CloudSecurityData, error) {
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

	var loadedFiles [][]byte
	for _, filename := range files {
		mainFileContent, err := readFile(filename)
		if err != nil {
			continue
		}
		loadedFiles = append(loadedFiles, mainFileContent)
	}

	terraformFootprintData, err := c.parseFootprintDataFromBytes(loadedFiles)
	if err != nil {
		return TerraformFootprintData{}, fmt.Errorf("error parsing footprint data from bytes: %w", err)
	}

	return terraformFootprintData, nil
}

// parseFootprintDataFromBytes parses Terraform footprint data from loaded of Terraform files
func (c *HTTPDragonDropClient) parseFootprintDataFromBytes(loadedFiles [][]byte) (TerraformFootprintData, error) {
	terraformFootprintData := TerraformFootprintData{}
	versionsTFModules := ModulesVersions{}

	for _, content := range loadedFiles {
		inputHCLFile, hclDiag := hclwrite.ParseConfig(
			content,
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

	return terraformFootprintData, err
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
