// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package ntag

import (
	"errors"
	"fmt"

	"github.com/happy-sdk/nfcsdk/internal/helpers"
	"github.com/happy-sdk/nfcsdk/pcsc"
)

type CapabilityContainer struct {
	TagType       Byte
	Version       Byte
	MemorySize    Size
	ReadOnly      bool
	accessControl byte
	raw           []byte
}

func NewCapabilityReadCmd() *pcsc.Command {
	cmd := pcsc.NewCustomCmd([]byte{0x30, 0x03})
	cmd.SetLe([]byte{18}) // 3 blocks + CRC
	cmd.SetName("GET_COMPATIBILITY")
	cmd.SetPostProcessFunc(func(data []byte) ([]byte, error) {
		if len(data) < 18 {
			return nil, fmt.Errorf("expected 3 block response with CRC got %d bytes", len(data))
		}
		blocks := data[:16]
		crc := data[16:]
		if !VerifyCRCA(blocks, crc) {
			return nil, errors.New("Capability Container (CC bytes) CRC check failed")
		}
		cc := data[:4]
		return cc, nil
	})
	return cmd
}

func (cc *CapabilityContainer) Unmarshal(data []byte) error {
	if len(data) < 4 {
		return fmt.Errorf("invalid data length: expected at least 4 bytes, got %d", len(data))
	}

	cc.TagType = tagType(data[0])
	cc.Version = tagVersion(data[1])
	cc.MemorySize = tagMemorySize(data[2])
	cc.accessControl = data[3]
	cc.ReadOnly = (cc.accessControl & ccReadOnly) != 0
	cc.raw = data
	return nil
}

func (cc *CapabilityContainer) String() string {
	return helpers.FormatByteSlice(cc.raw)
}

const ccReadOnly byte = 0x40

type tagType byte

func (t tagType) Byte() byte { return byte(t) }

// Interprets MagicNumber to return the tag type.
func (t tagType) String() string {
	switch t {
	case 0xE1:
		return "NFC Forum Type 2 Tag"
		// Add cases for other known magic numbers
	default:
		return "Unknown"
	}
}

type tagVersion byte

func (tv tagVersion) Byte() byte { return byte(tv) }

// String Converts Version to a human-readable format.
func (tv tagVersion) String() string {
	major := tv >> 4
	minor := tv & 0x0F
	return fmt.Sprintf("%d.%d", major, minor)
}

type tagMemorySize byte

func (ms tagMemorySize) Byte() byte { return byte(ms) }

// Interprets MemorySize to return total memory available for NDEF messages.
func (ms tagMemorySize) Size() int {
	switch ms {
	case 0x12:
		return 144 // NTAG213
	case 0x3E:
		return 496 // NTAG215
	case 0x6D:
		return 872 // NTAG216
	default:
		return 0 // Unknown or unsupported tag
	}
}

// Interprets MemorySize to return total memory available for NDEF messages.
func (ms tagMemorySize) String() string {
	return helpers.HumanizeBytes(int64(ms.Size()))
}
