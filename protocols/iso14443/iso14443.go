// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

// Package iso14443 provides tools and utilities for handling ISO/IEC 14443 standard,
// which is primarily used in proximity cards and NFC. It includes support for card
// detection, data exchange, and protocol-specific functionalities.
package iso14443

const (
	// Constants for ISO 14443 specific parameters
	BaudRate106 = "106 kbps"
	// ...
)

// DetectCard performs a proximity card detection process.
func DetectCard(readerConfig ReaderConfig) (*Card, error) { return nil, nil }

// ReadData reads data from an ISO 14443 card.
func ReadData(card *Card, blockAddress byte) ([]byte, error) { return nil, nil }

// WriteData writes data to an ISO 14443 card.
func WriteData(card *Card, blockAddress byte, data []byte) error { return nil }

// UnmarshalCard sets the Card fields from a byte slice.
func UnmarshalCard(data []byte) (*Card, error) { return nil, nil }

// CheckCardCompatibility checks if a card is compatible with ISO 14443 standards.
func CheckCardCompatibility(card *Card) bool { return false }

// Card represents an ISO 14443 card with specific attributes.
type Card struct {
	// Fields like UID, ATS (Answer to Select), etc.
}

// MarshalCard serializes a Card into a byte slice.
func (c *Card) Marshal() ([]byte, error) { return nil, nil }

// ReaderConfig represents configuration settings for an ISO 14443 reader.
type ReaderConfig struct {
	// Configuration fields like baud rate, modulation, etc.
}
