// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package ntag

import (
	"bytes"
	"testing"
)

func TestCalculateCRCA(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{
			name:  "Test 1",
			input: []byte{0x00, 0x04, 0x04, 0x02, 0x01, 0x00, 0x0F, 0x03},
			want:  []byte{0x80, 0x91},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateCRCA(tt.input); !bytes.Equal(got, tt.want) {
				t.Errorf("CalculateCRCA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyCRCA(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		receivedCRC []byte
		want        bool
	}{
		{
			name:        "Valid CRC",
			data:        []byte{0x00, 0x04, 0x04, 0x02, 0x01, 0x00, 0x0F, 0x03},
			receivedCRC: []byte{0x80, 0x91},
			want:        true,
		},
		{
			name:        "Invalid CRC",
			data:        []byte{0x00, 0x04, 0x04, 0x02, 0x01, 0x00, 0x0F, 0x03},
			receivedCRC: []byte{0x00, 0x00},
			want:        false,
		},
		// Additional test cases can be added here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyCRCA(tt.data, tt.receivedCRC); got != tt.want {
				t.Errorf("VerifyCRCA() = %v, want %v", got, tt.want)
			}
		})
	}
}
