package resourcesCalculator

import (
	"reflect"
	"testing"

	driftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector/drift_detector"
)

func TestCreateDivisionToNewResourceData(t *testing.T) {
	// Given
	c := TerraformResourcesCalculator{}
	inputBytesJSON := []byte(`{
"aws_lb_listener.tfer--number_1":"placeholder"
}`)

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
		inputBytesJSON,
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
