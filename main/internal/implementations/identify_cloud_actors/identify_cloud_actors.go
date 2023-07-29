package identifyCloudActors

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/Jeffail/gabs/v2"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// Config is a collection of query_param_data that parameterizes a IdentifyCloudActors instance.
type Config struct {
	// CloudCredential is a cloud credential with read-only access to a cloud division and, if applicable, access to read Terraform state files.
	CloudCredential terraformValueObjects.Credential `required:"true"`
}

// IdentifyCloudActors implements the interfaces.IdentifyCloudActors interface.
type IdentifyCloudActors struct {
	// Config is a collection of query_param_data that parameterizes a IdentifyCloudActors instance.
	config Config

	// dragonDrop is a client for interacting with the dragondrop API.
	dragonDrop interfaces.DragonDrop

	// providerToLogQuerier is a map between each provider and an instantiation of that provider's logQuerier.
	providerToLogQuerier map[terraformValueObjects.Provider]LogQuerier

	// DivisionToProvider is a map between the string representing a division and the corresponding
	// cloud provider (aws, azurerm, google, etc.).
	// For AWS, an account is the division, for GCP a project name is the division,
	// and for azurerm a resource group is a division.
	provider terraformValueObjects.Provider `required:"true"`
}

// NewIdentifyCloudActors returns a new instance of IdentifyCloudActors.
func NewIdentifyCloudActors(config Config, dragonDrop interfaces.DragonDrop, provider terraformValueObjects.Provider) (interfaces.IdentifyCloudActors, error) {
	providerToLogQuerier, err := NewProviderToLogQuerierMap(config, divisionToProvider)
	if err != nil {
		return nil, fmt.Errorf("[NewProviderToLogQuerierMap]%w", err)
	}

	return &IdentifyCloudActors{
		config:               config,
		providerToLogQuerier: providerToLogQuerier,
		dragonDrop:           dragonDrop,
		divisionToProvider:   divisionToProvider,
	}, nil
}

// Execute creates structured query_param_data mapping new or drifted resources to the cloud actor (service principal or user)
// responsible for the latest changes for that resource.
func (ica *IdentifyCloudActors) Execute(ctx context.Context) error {
	providerResourceActions := terraformValueObjects.ProviderResourceActions{}
	for provider, logQuerier := range ica.providerToLogQuerier {
		fmt.Printf("Beginning to pull cloud actors for all %v divisions\n", provider)
		currentProviderResourceActions, err := logQuerier.QueryForAllResources(ctx)
		if err != nil {
			return fmt.Errorf("[%v logQuerier.QueryForAllResources]%v", provider, err)
		}

		providerResourceActions[provider] = currentProviderResourceActions
	}

	jsonBytes, err := ica.convertProviderResourceActionsToJSON(providerResourceActions)
	if err != nil {
		return fmt.Errorf("[ica.convertProviderResourceActionsToJSON]%v", err)
	}
	err = os.WriteFile("mappings/resources-to-cloud-actions.json", jsonBytes, 0400)
	if err != nil {
		return fmt.Errorf("[os.WriteFile mappings/resources-to-cloud-actions.json]%v", err)
	}
	return nil
}

// convertProviderResourceActionsToJSON takes as input an object of type terraformValueObjects.ProviderResourceActions
// and outputs a formatted JSON equivalent of the struct.
func (ica *IdentifyCloudActors) convertProviderResourceActionsToJSON(actions terraformValueObjects.ProviderResourceActions) ([]byte, error) {
	jsonObj := gabs.New()

	for provider, divisionToResourceAction := range actions {
		for division, resourceNameToResourceAction := range divisionToResourceAction {
			for resourceName, resourceActions := range resourceNameToResourceAction {
				if resourceActions.Creator.Actor != "" {
					_, err := jsonObj.Set(resourceActions.Creator.Actor, string(provider), string(division), string(resourceName), "creation", "actor")
					if err != nil {
						return nil, fmt.Errorf("[jsonObj.Set(resourceActions.Creator.Actor] %v", err)
					}
					_, err = jsonObj.Set(resourceActions.Creator.Timestamp, string(provider), string(division), string(resourceName), "creation", "timestamp")
					if err != nil {
						return nil, fmt.Errorf("[jsonObj.Set(resourceActions.Creator.Timestamp] %v", err)
					}
				}
				if resourceActions.Modifier.Actor != "" {
					_, err := jsonObj.Set(resourceActions.Modifier.Actor, string(provider), string(division), string(resourceName), "modified", "actor")
					if err != nil {
						return nil, fmt.Errorf("[jsonObj.Set(resourceActions.Modifier.Actor] %v", err)
					}
					_, err = jsonObj.Set(resourceActions.Modifier.Timestamp, string(provider), string(division), string(resourceName), "modified", "timestamp")
					if err != nil {
						return nil, fmt.Errorf("[jsonObj.Set(resourceActions.Modifier.Timestamp] %v", err)
					}
				}
			}
		}
	}
	return jsonObj.Bytes(), nil
}

// executeCommand wraps os.exec.Command with capturing of std output and errors.
func executeCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	// Setting up logging objects
	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("[error executing command: %s, %s]%w", stderr.String(), out.String(), err)
	}
	fmt.Printf("\n%s Output:\n\n%v\n", command, out.String())
	return out.String(), nil
}
