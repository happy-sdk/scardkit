// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package ntag

import (
	"errors"
	"fmt"

	"github.com/happy-sdk/nfcsdk/pcsc"
)

// NewGetUIDCmd creates and returns a new pcsc.Command configured to retrieve the Unique Identifier (UID) of an NFC tag.
// The command is initialized with specific instruction bytes (0xFF, 0xCA) and sets both parameters to pcsc.ZeroByte,
// which is a standard request for retrieving the UID. The command expects a response length of 7 bytes, corresponding
// to the typical length of an NFC tag UID. It is assigned the name "GET_UID".
//
// Note: This function is intended for use with NFC tags supporting standard UID request commands as per NFC specifications.
// Separate handling of the UID response is required.
func NewGetUIDCmd() *pcsc.Command {
	cmd := pcsc.NewCmd(0xFF, 0xCA, pcsc.ZeroByte, pcsc.ZeroByte)
	cmd.SetLe([]byte{7})
	cmd.SetName("GET_UID")
	return cmd
}

// NewGetVersionCmd creates a new Command for the GET_VERSION NFC command, which is used to retrieve the
// version information of an NFC tag. It constructs a PC/SC command with the appropriate instruction byte.
// The command expects a response of 10 bytes:
// 8 bytes for product version information
// 2 bytes for CRC.
func NewGetVersionCmd() *pcsc.Command {
	cmd := pcsc.NewRawCmd([]byte{0x60})
	cmd.SetLe([]byte{10})
	cmd.SetName("GET_VERSION")
	return cmd
}

// NewCapabilityReadCmd creates and returns a pcsc.Command to read the Capability Container (CC) from page 3 of NFC Forum Type 2 Tags,
// such as NTAG213, NTAG215, and NTAG216. The command, named "GET_COMPATIBILITY", is set with a predefined instruction byte array,
// expecting an 18-byte response, which includes data and a CRC check. The post-processing function validates the response for
// length and CRC correctness, returning the first 4 bytes (CC) upon success. It returns an error for short responses or CRC failures.
//
// Note: This function targets NFC Forum Type 2 Tag compliance and requires separate handling of the response.
func NewCCReadCmd() *pcsc.Command {
	cmd := pcsc.NewRawCmd([]byte{0x30, 0x03})
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

// NewSignatureCmd creates a new command for reading a signature from a PC/SC device.
// It returns a pointer to a pcsc.Command.
func NewSignatureCmd() *pcsc.Command {
	// Create a new PC/SC command with the initial byte sequence 0x3C, 0x00.
	cmd := pcsc.NewRawCmd([]byte{0x3C, 0x00})
	cmd.SetName("READ_SIG")

	// Set the expected response length to 34 bytes (32 bytes for signature + 2 for CRC).
	cmd.SetLe([]byte{34})

	// Set a post-processing function that will be called after receiving the response data.
	cmd.SetPostProcessFunc(func(data []byte) ([]byte, error) {
		if len(data) < 34 {
			return nil, fmt.Errorf("expected 32 bytes response with CRC, got %d bytes", len(data))
		}

		// Extract the signature (first 32 bytes) and CRC (last 2 bytes) from the response.
		signatureb := data[:32]
		crc := data[32:]

		// Verify the CRC of the signature.
		if !VerifyCRCA(signatureb, crc) {
			return nil, errors.New("Signature CRC check failed")
		}

		return signatureb, nil
	})

	return cmd
}
