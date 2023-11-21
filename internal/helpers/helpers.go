// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package helpers

import (
	"fmt"
	"strings"
)

// FormatByteSlice converts a byte slice to a human-readable hex string.
// Each byte is represented in two hexadecimal characters, separated by colons.
func FormatByteSlice(slice []byte) string {
	const f = "%02X" // Define the format specifier
	var b strings.Builder

	// Write each byte in hex format to the builder
	for i, v := range slice {
		if i > 0 {
			b.WriteString(":") // Add a colon before each byte except the first one
		}
		b.WriteString(fmt.Sprintf(f, v))
	}

	return b.String()
}

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
