package resourcesCalculator

import (
	"reflect"
	"testing"

	"github.com/Jeffail/gabs/v2"
	driftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector/drift_detector"
)

func TestCreateDivisionToNewResourceData(t *testing.T) {
	// Given
	c := TerraformResourcesCalculator{}
	inputGabsContainer, err := gabs.ParseJSON([]byte(`{
"aws-dragondrop-dev.aws_lb_listener.tfer--number_1":"placeholder"
}`))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	inputTerraformerStateFile := driftDetector.TerraformerStateFile{
		Resources: []*driftDetector.TerraformerResource{
			{
				Mode:     "managed",
				Type:     "aws_lb_listener",
				Name:     "tfer--number_1",
				Provider: "aws",
				Instances: []driftDetector.TerraformerInstance{
					{
						AttributesFlat: map[string]string{
							"id": "arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-load-balancer/50dc6c495c0c9188/30dc6c495c0c9189",
						},
					},
				},
			},
		},
	}

	expectedOutput := map[ResourceID]NewResourceData{
		"arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-load-balancer/50dc6c495c0c9188/30dc6c495c0c9189": {
			ResourceType:            "aws_lb_listener",
			ResourceTerraformerName: "tfer--number_1",
			Region:                  "us-east-1",
		},
	}

	// When
	output, err := c.createNewResourceData(
		inputGabsContainer,
		inputTerraformerStateFile,
	)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Then
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected output to be:\n%v\ngot:\n%v\n", expectedOutput, output)
	}
}
