package terraformerCLI

import (
	"testing"

	"github.com/stretchr/testify/assert"

	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
)

func Test_terraformerCLI_getGroupListByResourceNames(t *testing.T) {
	type fields struct {
		config Config
	}
	type args struct {
		list []terraformValueObjects.ResourceName
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "Test with multiple resources",
			fields: fields{
				config: Config{},
			},
			args: args{
				list: []terraformValueObjects.ResourceName{
					"google_compute_address",
					"google_compute_autoscaler",
					"google_bigquery_dataset",
				},
			},
			want: []string{
				"addresses",
				"autoscalers",
				"bigQuery",
			},
		},
		{
			name: "Test with single resource",
			fields: fields{
				config: Config{},
			},
			args: args{
				list: []terraformValueObjects.ResourceName{
					"google_compute_disk",
				},
			},
			want: []string{
				"disks",
			},
		},
		{
			name: "Test with empty list",
			fields: fields{
				config: Config{},
			},
			args: args{
				list: []terraformValueObjects.ResourceName{},
			},
			want: []string{},
		},
		{
			name: "Test with multi-provider resource list",
			fields: fields{
				config: Config{},
			},
			args: args{
				list: []terraformValueObjects.ResourceName{
					"aws_accessanalyzer_analyzer",
					"google_compute_region_disk",
					"aws_elb",
					"google_cloudbuild_trigger",
				},
			},
			want: []string{
				"accessanalyzer",
				"regionDisks",
				"elb",
				"cloudbuild",
			},
		},
		{
			name: "Test only with aws resource list",
			fields: fields{
				config: Config{},
			},
			args: args{
				list: []terraformValueObjects.ResourceName{
					"aws_lambda_event_source_mapping",
					"aws_opsworks_static_web_layer",
					"aws_route53_zone",
				},
			},
			want: []string{
				"lambda",
				"opsworks",
				"route53",
			},
		},
		{
			name: "Test with azure resources",
			fields: fields{
				config: Config{},
			},
			args: args{
				list: []terraformValueObjects.ResourceName{
					"azurerm_resource_group",
					"azurerm_ssh_public_key",
					"azurerm_postgresql_database",
				},
			},
			want: []string{
				"resource_group",
				"virtual_machine",
				"database",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tfrCLI := &terraformerCLI{
				config: tt.fields.config,
			}
			assert.Equalf(t, tt.want, tfrCLI.getGroupListByResourceNames(tt.args.list), "getGroupListByResourceNames(%v)", tt.args.list)
		})
	}
}
