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

// TODO: Refactor of these unit tests needed
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
				CloudCredential:    terraformValueObjects.Credential(`{"awsAccessKeyID": "AWS123", "awsSecretAccessKey": "DUGFVGBHAJ213"}`),
				VCSRepo:            "https://github.com/test-org/test-repo.git",
			},
			},
			want: InferredData{
				Provider:  terraformValueObjects.Provider("aws"),
				VCSSystem: "github",
			},
			wantErr: false,
		},
		{
			name: "three providers azurerm",
			args: args{config: JobConfig{
				IsManagedDriftOnly: false,
				CloudCredential:    terraformValueObjects.Credential(`{"client_id": "123", "client_secret": "secret", "tenant_id": "tenant", "subscription_id": "subscription1"}`),
				VCSRepo:            "https://github.com/test-org/test-repo.git",
			},
			},
			want: InferredData{
				Provider:  terraformValueObjects.Provider("azurerm"),
				VCSSystem: "github",
			},
			wantErr: false,
		},
		{
			name: "three providers aws, azurerm and google",
			args: args{config: JobConfig{
				IsManagedDriftOnly: false,
				CloudCredential: terraformValueObjects.Credential(
					`{  "type": "service_account", "project_id": "project", "private_key_id": "123", "private_key": "key", 
							"client_email": "example@dragondrop.cloud", "client_id": "123456", "auth_uri": "https://accounts.google.com/o/oauth2/auth",
							"token_uri": "https://oauth2.googleapis.com/token", "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
							"client_x509_cert_url": "https://localhost.com"
						}`),
				VCSRepo: "https://github.com/test-org/test-repo.git",
			},
			},
			want: InferredData{
				Provider:  terraformValueObjects.Provider("google"),
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
	expectedOutput := []byte(
		`{"awsAccessKeyID": "AWS123",
      "awsSecretAccessKey": "Secret123"
}`)

	// When
	credential, err := parseAWSCredentialValues(inputCredentialBytes)
	if err != nil {
		t.Errorf("parseAWSCredentialValues() unexpected error:%v", err)
	}

	// Then
	assert.Equal(t, terraformValueObjects.Credential(expectedOutput), credential)
}
