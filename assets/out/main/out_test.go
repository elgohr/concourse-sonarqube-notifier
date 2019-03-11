package main

import (
	"bytes"
	"testing"
)

func TestDoesNothing(t *testing.T) {
	stdout := &bytes.Buffer{}
	status := run(stdout)

	if status != 0 {
		t.Errorf("Expected status 0, but got %v", status)
	}
	if stdout.String() != `[]` {
		t.Errorf("Expected empty array, but got %v", stdout.String())
	}
}
