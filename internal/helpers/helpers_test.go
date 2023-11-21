// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package helpers

import (
	"fmt"
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

func TestHumanizeBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1048576, "1.00 MB"},
		{1572864, "1.50 MB"},
		{1073741824, "1.00 GB"},
		{1610612736, "1.50 GB"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%dBytes", test.bytes), func(t *testing.T) {
			result := HumanizeBytes(test.bytes)
			if result != test.expected {
				t.Errorf("HumanizeBytes(%d) = %s; want %s", test.bytes, result, test.expected)
			}
		})
	}
}
