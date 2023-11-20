// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package helpers

import (
	"fmt"
	"strings"
)

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
