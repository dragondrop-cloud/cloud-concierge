package terraformerCLI

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func TestAWSScanner_configureEnvironment(t *testing.T) {
	// Given
	scanner := AWSScanner{}
	credentials := terraformValueObjects.Credential("{\n\"awsAccessKeyID\": \"123456ASD\",\n\"awsSecretAccessKey\": \"987654MNB\"\n}")

	// When
	err := scanner.configureEnvironment(credentials)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "123456ASD", os.Getenv("AWS_ACCESS_KEY_ID"))
	assert.Equal(t, "987654MNB", os.Getenv("AWS_SECRET_ACCESS_KEY"))
}
