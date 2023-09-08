package terraformvalueobjects

import (
	"reflect"
	"testing"
)

func TestResourceNameList_Decode(t *testing.T) {
	// Simple case with expected success
	envVar := ResourceNameList{}
	input := `["google_storage_bucket","aws_acm_certificate"]`

	err := envVar.Decode(input)
	if err != nil {
		t.Errorf("unexpected error in envVar.Decode: %v", err)
	}

	expectedValue := ResourceNameList{
		"google_storage_bucket",
		"aws_acm_certificate",
	}

	if !reflect.DeepEqual(expectedValue, envVar) {
		t.Errorf("got:\n%v\nexpected:\n%v\n", envVar, expectedValue)
	}

	// Simple case with None passed in
	envVar = ResourceNameList{}
	input = "None"

	err = envVar.Decode(input)
	if err != nil {
		t.Errorf("unexpected error in envVar.Decode: %v", err)
	}

	expectedValue = nil

	if !reflect.DeepEqual(expectedValue, envVar) {
		t.Errorf("got:\n%v\nexpected:\n%v\n", envVar, expectedValue)
	}

	// Less simple case with expected error
	envVar = ResourceNameList{}
	input = `"google_storage_bucket","aws_acm_certificate"`

	err = envVar.Decode(input)
	if err == nil {
		t.Errorf("expected error in envVar.Decode: %v", err)
	}

	// Less simple case with expected error
	envVar = ResourceNameList{}
	input = "[]"

	err = envVar.Decode(input)
	if err != nil {
		t.Errorf("unexpected error in envVar.Decode: %v", err)
	}

	if len(envVar) == 0 {
		t.Errorf("got length:\n%v\nexpected length of 0.", len(envVar))
	}
}
