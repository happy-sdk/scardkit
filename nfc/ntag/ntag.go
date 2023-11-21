// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

// Package ntag specializes in the NTAG family of NFC tags (NTAG213, NTAG215, NTAG216).
// It offers tailored functionalities for reading, writing, and interacting with these
// NFC tags, based on the specifications provided by NXP.
package ntag

const (
	// Constants defining page sizes, memory layout, etc.
	PageSize = 4 // bytes
	// ...
)

// NewTag initializes a new NTAG NFC tag representation.
func NewTag(uid []byte) *Tag

// Tag represents an NTAG NFC tag with its attributes.
type Tag struct {
	// Fields like UID, memory layout, etc.
}

// ReadPage reads a specific page from the NFC tag.
func (t *Tag) ReadPage(pageNumber int) (*Page, error) { return nil, nil }

// WritePage writes data to a specific page on the NFC tag.
func (t *Tag) WritePage(pageNumber int, data []byte) error { return nil }

// CalculateTagCapacity returns the total memory capacity of the tag.
func (t *Tag) CalculateTagCapacity() int { return 0 }

// VerifyIntegrity checks the integrity of the NFC tag data.
func (t *Tag) VerifyIntegrity() bool { return false }

// Page represents a data page in the NTAG NFC tag.
type Page struct {
	// Fields representing page data and other attributes.
}

// MarshalPage serializes a Page into a byte slice.
func (p *Page) Marshal() ([]byte, error)

// UnmarshalPage sets the Page fields from a byte slice.
func (p *Page) Unmarshal(data []byte) error
