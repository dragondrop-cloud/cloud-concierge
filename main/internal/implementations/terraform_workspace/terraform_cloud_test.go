package terraformworkspace

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
