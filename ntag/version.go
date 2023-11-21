// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package ntag

import (
	"fmt"

	"github.com/happy-sdk/nfcsdk/internal/helpers"
	"github.com/happy-sdk/nfcsdk/pcsc"
)

type Version struct {
	VendorID       Byte
	ProductType    Byte
	ProductSubtype Byte
	MajorVersion   uint8
	MinorVersion   uint8
	StorageSize    Size
	ProtocolType   ProtocolType
	Valid          bool
}

// NewGetVersionCmd creates a new Command for the GET_VERSION NFC command, which is used to retrieve the
// version information of an NFC tag. It constructs a PC/SC command with the appropriate instruction byte.
// The command expects a response of 10 bytes:
// 8 bytes for product version information
// 2 bytes for CRC.
func NewGetVersionCmd() *pcsc.Command {
	cmd := pcsc.NewCustomCmd([]byte{0x60})
	cmd.SetLe([]byte{10})
	cmd.SetName("GET_VERSION")
	return cmd
}

func (v *Version) Unmarshal(payload []byte) error {
	if l := len(payload); l < 7 {
		return fmt.Errorf("invalid data length: expected at least 7 bytes, got %d", l)
	}
	var (
		data = payload[0:8] // First 8 bytes are data
		crc  []byte         // Last 2 bytes are CRC
	)

	if len(payload) == 10 {
		crc = payload[8:]
		v.Valid = VerifyCRCA(data, crc)
	}
	v.VendorID = vendorID(data[1])
	v.ProductType = productType(data[2])
	v.ProductSubtype = productSubtype(data[3])

	v.MajorVersion = data[4]
	v.MinorVersion = data[5]
	v.StorageSize = parseVersionStorageSize(data[6])
	if data[7] == 0x03 {
		v.ProtocolType = ProtocolType(data[7])
	}
	return nil
}

type ProtocolType uint8

const (
	ProtocolTypeISO14443A = 0x01
	ProtocolTypeISO14443B = 0x02
	ProtocolTypeISO144433 = 0x03
)

var protocolTypes = map[ProtocolType]string{
	ProtocolTypeISO14443A: "ISO/IEC 14443 Type A",
	ProtocolTypeISO14443B: "ISO/IEC 14443 Type B",
	ProtocolTypeISO144433: "ISO/IEC 14443-3 compliant",
	// Add more protocol types as required by the supported NFC tags.
}

func (p ProtocolType) String() string {
	typstr, ok := protocolTypes[p]
	if !ok {
		// If the protocol type is unknown or not supported, you can return a default message or an error.
		return fmt.Sprintf("unknown protocol type (0x%02X)", uint8(p))
	}
	return typstr
}

type vendorID byte

func (v vendorID) Byte() byte { return byte(v) }
func (v vendorID) String() string {
	if v == 0x04 {
		return "NXP Semiconductors"
	}
	return fmt.Sprintf("Unkown vendor id %02Xh", byte(v))
}

type productType byte

func (v productType) Byte() byte { return byte(v) }
func (v productType) String() string {
	if v == 0x04 {
		return "NXP NTAG"
	}
	return fmt.Sprintf("Unkown product type %02Xh", byte(v))
}

type productSubtype byte

func (v productSubtype) Byte() byte { return byte(v) }
func (v productSubtype) String() string {
	if v == 0x02 {
		return "50 pF"
	}
	return fmt.Sprintf("Unkown product subtype %02Xh", byte(v))
}

type storageSize struct {
	b byte
	v int
}

func (s storageSize) Byte() byte     { return s.b }
func (s storageSize) Size() int      { return s.v }
func (s storageSize) String() string { return helpers.HumanizeBytes(int64(s.v)) }

func parseVersionStorageSize(sizeb byte) (size storageSize) {
	size.b = sizeb
	// Extract the most significant 7 bits and calculate 2^n
	n := int(sizeb >> 1)
	base := 1 << n // 2^n
	// Check if it's one of the known models
	if sizeb&1 == 1 {
		// Precise sizes for known models
		switch sizeb {
		case 0x0F: // NTAG213
			size.v = 144
		case 0x11: // NTAG215
			size.v = 504
		case 0x13: // NTAG216
			size.v = 888
		default:
			// For unknown models, return an estimate (e.g., the average of the range)
			size.v = (base + (base << 1)) / 2
		}
	} else {
		// If the least significant bit is 0, the size is exactly 2^n
		size.v = base
	}
	return size
}
