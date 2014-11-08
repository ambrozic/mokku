package tests

import (
	"testing"
)

func TestTesting(t *testing.T) {

	if "a" != "a" {
		t.Error("'a' should be 'a'")
	}
	if "a" == "b" {
		t.Error("'a' should not be 'b'")
	}

	//	assert.Equal(t, "a", "a", "'a' should be 'a'");
	//	assert.NotEqual(t, "a", "b", "'a' should not be 'b'");
	//
	//	assert.Equal(t, 1, 1, "1 should be 1");
	//	assert.NotEqual(t, 1, 2, "1 should not be 2");
	//
	//	assert.NotNil(t, obj, "obj shouldn't be nil");
}
