package terraformerCLI

import (
	"testing"

	"github.com/stretchr/testify/assert"

	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
)

func Test_subsetMapOfDivisionToCredentials(t *testing.T) {
	type args struct {
		divisionCloudCredentials terraformValueObjects.DivisionCloudCredentialDecoder
		divisionToProvider       map[terraformValueObjects.Division]terraformValueObjects.Provider
		provider                 terraformValueObjects.Provider
	}
	tests := []struct {
		name string
		args args
		want map[terraformValueObjects.Division]terraformValueObjects.Credential
	}{
		{
			name: "Two divisions",
			args: args{
				divisionCloudCredentials: terraformValueObjects.DivisionCloudCredentialDecoder{
					terraformValueObjects.Division("division-1"): terraformValueObjects.Credential(`{"field1":"value1", "field2": "value2"}`),
					terraformValueObjects.Division("division-2"): terraformValueObjects.Credential(`{"field1":"value1", "field2": "value2"}`),
				},
				divisionToProvider: map[terraformValueObjects.Division]terraformValueObjects.Provider{
					terraformValueObjects.Division("division-1"): terraformValueObjects.Provider("azurerm"),
					terraformValueObjects.Division("division-2"): terraformValueObjects.Provider("azurerm"),
				},
				provider: "azurerm",
			},
			want: map[terraformValueObjects.Division]terraformValueObjects.Credential{
				terraformValueObjects.Division("division-1"): terraformValueObjects.Credential(`{"field1":"value1", "field2": "value2"}`),
				terraformValueObjects.Division("division-2"): terraformValueObjects.Credential(`{"field1":"value1", "field2": "value2"}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, subsetMapOfDivisionToCredentials(tt.args.divisionCloudCredentials, tt.args.divisionToProvider, tt.args.provider), "subsetMapOfDivisionToCredentials(%v, %v, %v)", tt.args.divisionCloudCredentials, tt.args.divisionToProvider, tt.args.provider)
		})
	}
}
