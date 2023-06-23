package documentize

import "testing"

func TestStringToWords(t *testing.T) {
	// case 1
	input := "module.google-backend-api_broski"

	expectedOutput := "google backend api broski"

	output := stringToWords(input)

	if output != expectedOutput {
		t.Errorf("got %v, expected %v", output, expectedOutput)
	}

	// case 2
	input = "module.google-iam-secrets.module.allowed_cors_origin"

	expectedOutput = "google iam secrets allowed cors origin"

	output = stringToWords(input)

	if output != expectedOutput {
		t.Errorf("got %v, expected %v", output, expectedOutput)
	}

	// case 3
	input = "projects/153130598315/secrets/allowed_cors_origin"

	expectedOutput = "allowed cors origin"

	output = stringToWords(input)

	if output != expectedOutput {
		t.Errorf("got %v, expected %v", output, expectedOutput)
	}
}
