package driftdetector

import "testing"

func TestParseIDFromField(t *testing.T) {
	if "" != parseIDFromField("") {
		t.Error("Expected empty string")
	}

	if "projects/test" != parseIDFromField("asd/asdasd/projects/test") {
		t.Error("Expected projects/test")
	}

	if "namespaces/test" != parseIDFromField("namespaces/test") {
		t.Error("Expected projects/test")
	}

	if "" != parseIDFromField("projects") {
		t.Error("Expected ''")
	}
}

func TestResourceIDCalculator_SelfLink(t *testing.T) {
	// Given
	attributesFlat := map[string]string{
		"self_link": "https://google.com/projects/test",
		"id":        "example-id",
	}

	// When
	output, _ := ResourceIDCalculator(attributesFlat, "google", "placeholder")

	// Then
	if output != "projects/test" {
		t.Errorf("got:\n%s\nexpected:\n%s\n", output, "projects/test")
	}
}

func TestResourceIDCalculator_ID(t *testing.T) {
	// Given
	attributesFlat := map[string]string{
		"id": "namespaces/example-id",
	}

	// When
	output, _ := ResourceIDCalculator(attributesFlat, "google", "placeholder")

	// Then
	if output != "namespaces/example-id" {
		t.Errorf("got:\n%s\nexpected:\n%s\n", output, "namespaces/example-id")
	}
}

func TestResourceIDCalculator_GoogleEdgeCase(t *testing.T) {
	// Given
	attributesFlat := map[string]string{
		"self_link": "test",
		"id":        "bucket_example",
	}

	// When
	output, _ := ResourceIDCalculator(attributesFlat, "google", "google_storage_bucket")

	// Then
	if output != "projects/_/buckets/bucket_example" {
		t.Errorf("got:\n%s\nexpected:\n%s\n", output, "projects/_/buckets/bucket_example")
	}
}

func TestResourceIDCalculator_NonGoogle(t *testing.T) {
	// Given
	attributesFlat := map[string]string{
		"self_link": "https://google.com/projects/test",
		"id":        "namespaces/example-id",
	}

	// When
	output, _ := ResourceIDCalculator(attributesFlat, "not-google", "placeholder")

	// Then
	if output != "namespaces/example-id" {
		t.Errorf("got:\n%s\nexpected:\n%s\n", output, "namespaces/example-id")
	}
}
