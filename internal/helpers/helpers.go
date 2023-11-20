package helpers

import (
	"fmt"
	"strings"
)

func FormatByteSlice(slice []byte) string {
	const hexFormat = "%02X" // Define the format specifier as a constant
	var b strings.Builder

	// Write each byte in hex format to the builder
	// The loop uses range to ensure it works correctly with any slice length
	for i, v := range slice {
		if i > 0 {
			b.WriteString(":") // Add a colon before each byte except the first one
		}
		b.WriteString(fmt.Sprintf(hexFormat, v)) // Use the constant format specifier
	}

	return b.String()
}
