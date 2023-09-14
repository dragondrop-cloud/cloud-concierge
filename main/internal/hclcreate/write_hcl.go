package hclcreate

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

// CreateMainTF outputs a bytes slice which defines a baseline main.tf file.
func (h *hclCreate) CreateMainTF(providers map[string]string) ([]byte, error) {
	logrus.Debugf("[hclcreate][write_hcl] CreateMainTF: %v", providers)

	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	terraformBlock := rootBody.AppendNewBlock("terraform", nil)
	terraformBody := terraformBlock.Body()
	terraformBody.SetAttributeValue("required_version", cty.StringVal(h.config.TerraformVersion))
	terraformBody.AppendNewline()

	requiredProvidersBlock := terraformBody.AppendNewBlock("required_providers", nil)
	requiredProvidersBody := requiredProvidersBlock.Body()

	for provider, version := range providers {
		err := requiredProviderSubBlock(requiredProvidersBody, provider, version)
		if err != nil {
			return nil, err
		}
	}

	if providers["azurerm"] != "" {
		providerBlock := rootBody.AppendNewBlock("provider", []string{"azurerm"})
		providerBody := providerBlock.Body()
		providerBody.AppendNewBlock("features", nil)
	}

	return f.Bytes(), nil
}

// requiredProviderSubBlock creates a sub-chunk of hcl within the passed body for a required provider
// and version.
func requiredProviderSubBlock(body *hclwrite.Body, provider string, version string) error {
	body.SetAttributeValue(string(provider), cty.ObjectVal(map[string]cty.Value{
		"source":  cty.StringVal(fmt.Sprintf("hashicorp/%v", string(provider))),
		"version": cty.StringVal(string(version)),
	}))
	body.AppendNewline()

	return nil
}
