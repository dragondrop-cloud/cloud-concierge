package terraformValueObjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
