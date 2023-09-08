package main

import (
	"testing"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getProviderByCredential(t *testing.T) {
	type args struct {
		provider map[terraformValueObjects.Provider]string
	}
	tests := []struct {
		name    string
		args    args
		want    terraformValueObjects.Provider
		wantErr bool
	}{
		{
			name: "azurerm",
			args: args{
				provider: map[terraformValueObjects.Provider]string{"azurerm": ""},
			},
			want:    "azurerm",
			wantErr: false,
		},
		{
			name: "aws",
			args: args{
				provider: map[terraformValueObjects.Provider]string{"aws": ""},
			},
			want:    "aws",
			wantErr: false,
		},
		{
			name: "google",
			args: args{
				provider: map[terraformValueObjects.Provider]string{"google": ""},
			},
			want:    "google",
			wantErr: false,
		},
		{
			name: "error inferring provider",
			args: args{
				provider: map[terraformValueObjects.Provider]string{"google": "", "aws": ""},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getProviderFromProviderVersion(tt.args.provider)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equalf(t, tt.want, got, "getProviderFromProviderVersion(%v)", tt.args.provider)
		})
	}
}

func Test_getInferredData(t *testing.T) {
	type args struct {
		config JobConfig
	}
	tests := []struct {
		name    string
		args    args
		want    InferredData
		wantErr bool
	}{
		{
			name: "single provider aws",
			args: args{config: JobConfig{
				IsManagedDriftOnly: false,
				Provider:           map[terraformValueObjects.Provider]string{"aws": ""},
				VCSRepo:            "https://github.com/test-org/test-repo.git",
				JobID:              "test-pull",
			},
			},
			want: InferredData{
				Provider:  terraformValueObjects.Provider("aws"),
				VCSSystem: "github",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getInferredData(tt.args.config)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equalf(t, tt.want, got, "getInferredData(%v)", tt.args.config)
		})
	}
}

func Test_parseAWSCredentialValues(t *testing.T) {
	// Given
	inputCredentialBytes := []byte(
		`[default]
aws_access_key_id = AWS123
aws_secret_access_key = Secret123
`)
	expectedOutput := []byte(`{"awsAccessKeyID":"AWS123","awsSecretAccessKey":"Secret123"}`)

	// When
	credential, err := parseAWSCredentialValues(inputCredentialBytes)
	if err != nil {
		t.Errorf("parseAWSCredentialValues() unexpected error:%v", err)
	}

	// Then
	assert.Equal(t, terraformValueObjects.Credential(expectedOutput), credential)
}

func Test_searchAwsAccess(t *testing.T) {
	type args struct {
		credentials []byte
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "aws credentials with spaces",
			args: args{
				credentials: []byte(
					`[default]
aws_access_key_id = AWS123
aws_secret_access_key = Secret123
`),
			},
			want: []string{"\naws_access_key_id = AWS123\naws_secret_access_key = Secret123", "AWS123", "Secret123"},
		},
		{
			name: "aws credentials without spaces",
			args: args{
				credentials: []byte(
					`[default]
aws_access_key_id=AWS123
aws_secret_access_key=Secret123
`),
			},
			want: []string{"\naws_access_key_id=AWS123\naws_secret_access_key=Secret123", "AWS123", "Secret123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, searchAwsAccess(tt.args.credentials), "searchAwsAccess(%v)", tt.args.credentials)
		})
	}
}
