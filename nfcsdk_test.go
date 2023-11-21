// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package nfcsdk

import (
	"testing"
)

func TestFormatByteSlice(t *testing.T) {
	// Test case 1: An empty byte slice should return an empty string.
	emptySlice := []byte{}
	result := FormatByteSlice(emptySlice)
	expected := ""
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}

	// Test case 2: A byte slice with some data should be formatted as expected.
	dataSlice := []byte{72, 101, 108, 108, 111} // ASCII values for "Hello"
	result = FormatByteSlice(dataSlice)
	expected = "48:65:6C:6C:6F" // Hexadecimal representation of ASCII values
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}
