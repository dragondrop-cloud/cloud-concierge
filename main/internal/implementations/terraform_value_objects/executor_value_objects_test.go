package terraformValueObjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDivisionCloudCredentialDecoder(t *testing.T) {
	tests := []struct {
		name      string
		envVarStr string
		expected  Credential
		division  string
	}{
		{
			name:      "google-dev",
			envVarStr: "google-dev:{'part_1': 'asd', 'part_2': 'xyz'}",
			expected:  Credential("{'part_1': 'asd', 'part_2': 'xyz'}"),
			division:  "google-dev",
		},
		{
			name:      "google-dev",
			envVarStr: "google-dev:{'part_1': 'asd', 'part_2': 'xyz'},google-prod:{'part_1': 'xyz', 'part_2': 'asd'}",
			expected:  Credential("{'part_1': 'asd', 'part_2': 'xyz'}"),
			division:  "google-dev",
		},
		{
			name:      "google-prod",
			envVarStr: "google-dev:{'part_1': 'asd', 'part_2': 'xyz'},google-prod:{'part_1': 'xyz', 'part_2': 'asd'}",
			expected:  Credential("{'part_1': 'xyz', 'part_2': 'asd'}"),
			division:  "google-prod",
		},
		{
			name:      "private-key",
			envVarStr: "google-dev:{'type': 'service_account','project_id': 'xyz-dev','private_key_id': '4f','private_key': '-----BEGIN PRIVATE KEY-----\nManJ\nzit0rTmOa\nlAmxD56XAFA==\n-----END PRIVATE KEY-----\n', 'client_x509_cert_url': 'https://xyz-dev.iam.gserviceaccount.com'}",
			expected:  Credential("{'type': 'service_account','project_id': 'xyz-dev','private_key_id': '4f','private_key': '-----BEGIN PRIVATE KEY-----\nManJ\nzit0rTmOa\nlAmxD56XAFA==\n-----END PRIVATE KEY-----\n', 'client_x509_cert_url': 'https://xyz-dev.iam.gserviceaccount.com'}"),
			division:  "google-dev",
		},
		{
			name:      "azure-credentials",
			envVarStr: "dragondrop-dev:{'client_id':'123','client_secret':'secret','tenant_id':'123456','subscription_id':'123-sub','resource_group': 'dragondrop-dev'}",
			expected:  Credential("{'client_id':'123','client_secret':'secret','tenant_id':'123456','subscription_id':'123-sub','resource_group': 'dragondrop-dev'}"),
			division:  "dragondrop-dev",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envVar := DivisionCloudCredentialDecoder{}
			err := envVar.Decode(tt.envVarStr)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, envVar[Division(tt.division)])
		})
	}
}

func TestVersionDecoder(t *testing.T) {
	tests := []struct {
		name          string
		envVarStr     string
		expected      Version
		expectedError bool
	}{
		{
			name:          "valid version",
			envVarStr:     "1.2.3",
			expected:      Version("1.2.3"),
			expectedError: false,
		},
		{
			name:          "invalid version",
			envVarStr:     "~>1.2.3",
			expected:      Version(""),
			expectedError: true,
		},
		{
			name:          "invalid version",
			envVarStr:     "~>1.2",
			expected:      Version(""),
			expectedError: true,
		},
		{
			name:          "invalid version",
			envVarStr:     "2.6.5",
			expected:      Version(""),
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var envVar Version
			err := envVar.Decode(tt.envVarStr)

			if tt.expectedError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, envVar)
			}
		})
	}
}
