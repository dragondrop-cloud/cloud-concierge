package terraformWorkspace

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerraformCloud_extractWorkspaceName(t *testing.T) {
	type args struct {
		ctx          context.Context
		versionsFile []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy path",
			args: args{
				ctx: context.Background(),
				versionsFile: []byte(`
					terraform {
					  required_version = "~> 1.2.6"
					
					  required_providers {
						google = {
						  source  = "hashicorp/google"
						  version = "~>4.31.0"
						}
					
						tfe = {
						  source  = "hashicorp/tfe"
						  version = "~>0.33.0"
						}
					  }
					
					  cloud {
						organization = "dragondrop-cloud"
					
						workspaces {
						  name = "google-backend-api-dev"
						}
					  }
					}
				`),
			},
			want:    "google-backend-api-dev",
			wantErr: assert.NoError,
		},
		{
			name: "no workspace name",
			args: args{
				ctx: context.Background(),
				versionsFile: []byte(`
					terraform {
					  required_version = "~> 1.2.6"
					
					  required_providers {
						google = {
						  source  = "hashicorp/google"
						  version = "~>4.31.0"
						}
					
						tfe = {
						  source  = "hashicorp/tfe"
						  version = "~>0.33.0"
						}
					  }
					}
				`),
			},
			want:    "",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &TerraformCloud{}
			got, err := c.extractWorkspaceName(tt.args.ctx, tt.args.versionsFile)
			if !tt.wantErr(t, err, fmt.Sprintf("extractWorkspaceName(%v, %v)", tt.args.ctx, tt.args.versionsFile)) {
				return
			}
			assert.Equalf(t, tt.want, got, "extractWorkspaceName(%v, %v)", tt.args.ctx, tt.args.versionsFile)
		})
	}
}

func TestTerraformCloud_cleanDirectoryName(t *testing.T) {
	type args struct {
		directory string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "happy path",
			args: args{
				directory: "google-backend-api-dev",
			},
			want: "google-backend-api-dev",
		},
		{
			name: "happy path with trailing slash",
			args: args{
				directory: "google-backend-api-dev/",
			},
			want: "google-backend-api-dev",
		},
		{
			name: "happy path with trailing slash and whitespace",
			args: args{
				directory: "google-backend-api-dev/ ",
			},
			want: "google-backend-api-dev",
		},
		{
			name: "happy path with trailing at beginning and end",
			args: args{
				directory: "/google-backend-api-dev/",
			},
			want: "google-backend-api-dev",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &TerraformCloud{}
			assert.Equalf(t, tt.want, c.cleanDirectoryName(tt.args.directory), "cleanDirectoryName(%v)", tt.args.directory)
		})
	}
}

func TestWorkspaceDirectoriesDecoder_Decode(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name        string
		d           WorkspaceDirectoriesDecoder
		args        args
		wantErr     assert.ErrorAssertionFunc
		valueWanted []string
	}{
		{
			name: "happy path",
			d:    WorkspaceDirectoriesDecoder{},
			args: args{
				value: "[\"google-backend-api-dev/\"]",
			},
			wantErr: assert.NoError,
			valueWanted: []string{
				"google-backend-api-dev/",
			},
		},
		{
			name: "2 directories",
			d:    WorkspaceDirectoriesDecoder{},
			args: args{
				value: "[\"google-backend-api-dev/\",\"google-backend-api-prod/\"]",
			},
			wantErr: assert.NoError,
			valueWanted: []string{
				"google-backend-api-dev/",
				"google-backend-api-prod/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, tt.d.Decode(tt.args.value), fmt.Sprintf("Decode(%v)", tt.args.value))
			assert.Equalf(t, tt.valueWanted, []string(tt.d), "Decode(%v)", tt.args.value)
		})
	}
}
