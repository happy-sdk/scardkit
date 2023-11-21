// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

// Package sckutils offers a suite of utilities and common functionalities
// designed to support various aspects of smart card communication and processing.
// It includes helper functions, common data structures, and utility routines
// that facilitate development across the scardkit's various smart card protocols
// and operations.
package sckutils

import (
	"fmt"
	"strings"
)

// ConvertToHexString converts a byte slice to a hex string.
func ConvertToHexString(data []byte) string {
	const f = "%02X" // Define the format specifier
	var b strings.Builder

	// Write each byte in hex format to the builder
	for i, v := range data {
		if i > 0 {
			b.WriteString(":") // Add a colon before each byte except the first one
		}
		b.WriteString(fmt.Sprintf(f, v))
	}

	return b.String()
}

// ParseHexString parses a hex string into a byte slice.
func ParseHexString(hexStr string) ([]byte, error) { return nil, nil }

// CheckStatusCode evaluates a status code to determine if an operation was successful.
func CheckStatusCode(code byte) bool { return false }

// NewDataBuffer creates a new DataBuffer with the provided data.
func NewDataBuffer(data []byte) *DataBuffer { return nil }

// UnmarshalDataBuffer sets the DataBuffer fields from a byte slice.
func UnmarshalDataBuffer(data []byte) (*DataBuffer, error) { return nil, nil }

// HumanizeBytes converts a size in bytes to a human-readable string in KB, MB, GB, etc.
func HumanizeBytes(bytes int64) string {
	const (
		kB int64 = 1 << 10 // 1024
		mB int64 = 1 << 20 // 1048576
		gB int64 = 1 << 30 // 1073741824
	)

	format := "%.2f %s"
	switch {
	case bytes < kB:
		return fmt.Sprintf("%d B", bytes)
	case bytes < mB:
		return fmt.Sprintf(format, float64(bytes)/float64(kB), "KB")
	case bytes < gB:
		return fmt.Sprintf(format, float64(bytes)/float64(mB), "MB")
	default:
		// When file size is larger than 1 MB
		return fmt.Sprintf(format, float64(bytes)/float64(gB), "GB")
	}
}

// DataBuffer represents a generic data buffer for smart card operations.
type DataBuffer struct {
	Data []byte
}

// Marshal serializes a DataBuffer into a byte slice.
func (db *DataBuffer) Marshal() ([]byte, error) { return nil, nil }
