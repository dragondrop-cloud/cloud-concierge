package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionChecker_getCloudConciergeVersionMessage(t *testing.T) {
	type args struct {
		currentVersion      string
		latestStableVersion string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "should return the correct message for the same version",
			args: args{
				currentVersion:      "v0.1.5",
				latestStableVersion: "v0.1.5",
			},
			want:  "You are currently running version v0.1.5 of cloud-concierge, this is the latest version.",
			want1: "green",
		},
		{
			name: "should return the warning message for different versions",
			args: args{
				currentVersion:      "v0.1.4",
				latestStableVersion: "v0.1.5",
			},
			want:  "You are currently running version v0.1.4 of cloud-concierge, the latest version is v0.1.5, run docker pull dragondrop/cloud-concierge:latest to update.",
			want1: "red",
		},
		{
			name: "should return the correct message for empty latest version",
			args: args{
				currentVersion:      "v0.1.4",
				latestStableVersion: "",
			},
			want:  "You are currently running version v0.1.4 of cloud-concierge, this is the latest version.",
			want1: "green",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VersionChecker{}
			got, got1 := v.getCloudConciergeVersionMessage(tt.args.currentVersion, tt.args.latestStableVersion)
			assert.Equalf(t, tt.want, got, "getCloudConciergeVersionMessage(%v)", tt.args.latestStableVersion)
			assert.Equalf(t, tt.want1, got1, "getCloudConciergeVersionMessage(%v)", tt.args.latestStableVersion)
		})
	}
}
