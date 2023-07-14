package terraformValueObjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloudRegionsDecoder_Decode(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name          string
		d             CloudRegionsDecoder
		args          args
		expectedValue []CloudRegion
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name:          "no brackets error",
			d:             CloudRegionsDecoder{},
			args:          args{value: `"us-east-1"`},
			expectedValue: []CloudRegion{},
			wantErr:       assert.Error,
		},
		{
			name:          "only first bracket error",
			d:             CloudRegionsDecoder{},
			args:          args{value: `["us-east-1"`},
			expectedValue: []CloudRegion{},
			wantErr:       assert.Error,
		},
		{
			name:          "only last bracket error",
			d:             CloudRegionsDecoder{},
			args:          args{value: `"us-east-1"]`},
			expectedValue: []CloudRegion{},
			wantErr:       assert.Error,
		},
		{
			name:          "no quotation marks error",
			d:             CloudRegionsDecoder{},
			args:          args{value: `[us-east-1]`},
			expectedValue: []CloudRegion{},
			wantErr:       assert.Error,
		},
		{
			name:          "no commas separator error",
			d:             CloudRegionsDecoder{},
			args:          args{value: `["us-east-1" "us-east-2"]`},
			expectedValue: []CloudRegion{},
			wantErr:       assert.Error,
		},
		{
			name:          "no value and return empty cloud regions list",
			d:             CloudRegionsDecoder{},
			args:          args{value: ``},
			expectedValue: []CloudRegion{},
			wantErr:       nil,
		},
		{
			name:          "empty list string and return empty cloud regions list",
			d:             CloudRegionsDecoder{},
			args:          args{value: `[]`},
			expectedValue: []CloudRegion{},
			wantErr:       nil,
		},
		{
			name:          "happy path",
			d:             CloudRegionsDecoder{},
			args:          args{value: `["us-east-1","us-west1", "westus2"]`},
			expectedValue: []CloudRegion{"us-east-1", "us-west1", "westus2"},
			wantErr:       nil,
		},
		{
			name:          "only one region per provider error aws",
			d:             CloudRegionsDecoder{},
			args:          args{value: `["us-east-1","us-east-2", "us-west1", "westus2"]`},
			expectedValue: []CloudRegion{},
			wantErr:       assert.Error,
		},
		{
			name:          "only one region per provider error azure",
			d:             CloudRegionsDecoder{},
			args:          args{value: `["us-east-1","us-west1","westus2","centraluseuap"]`},
			expectedValue: []CloudRegion{},
			wantErr:       assert.Error,
		},
		{
			name:          "only one region per provider error gcp",
			d:             CloudRegionsDecoder{},
			args:          args{value: `["us-east-1","us-west1","europe-west3","westus2"]`},
			expectedValue: []CloudRegion{},
			wantErr:       assert.Error,
		},
		{
			name:          "invalid zones error",
			d:             CloudRegionsDecoder{},
			args:          args{value: `["region123", "westus2"]`},
			expectedValue: []CloudRegion{},
			wantErr:       assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.d.Decode(tt.args.value)

			if tt.wantErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedValue, []CloudRegion(tt.d))
		})
	}
}
