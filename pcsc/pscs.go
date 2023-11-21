// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

// Package pcsc provides a Go wrapper for libpcsclite, facilitating communication
// with smart card readers using the PC/SC standard. It offers essential functions
// to connect, communicate, and interact with smart cards through readers.
package pscs

const (
	// Constants defining card and reader states, error codes, etc.
	CardAbsent  = "CardAbsent"
	CardPresent = "CardPresent"
	// ...
)

// ListReaders lists available PC/SC readers connected to the system.
func ListReaders() ([]ReaderInfo, error)

// ConnectToCard establishes a connection with a card in the specified reader.
func ConnectToCard(readerName string) (*Card, error)

// ReaderStatus checks the status of a specified PC/SC reader.
func ReaderStatus(readerName string) (string, error)

// ReaderInfo represents information about a PC/SC reader.
type ReaderInfo struct {
	// Fields like reader name, connection status, etc.
}

// Card represents a smart card in a PC/SC reader.
type Card struct {
	// Fields representing card properties and status.
}

// Transmit sends an APDU command to the card and receives a response.
func (c *Card) Transmit(apduCommand []byte) ([]byte, error)

// Disconnect releases the connection with the card.
func (c *Card) Disconnect() error

// CardStatus retrieves the current status of a smart card.
func (c *Card) CardStatus() (string, error)
