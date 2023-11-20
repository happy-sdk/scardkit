// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package ntag

import "bytes"

// CalculateCRCA calculates the CRC-A checksum for a given slice of data.
// This function implements the CRC-A algorithm used in ISO/IEC 14443 Type A
// standard, suitable for NFC communication. The function returns a 2-byte slice
// containing the CRC in little-endian format.
func CalculateCRCA(data []byte) []byte {
	poly := uint16(0x1021)
	crc := uint16(0xC6C6)

	for _, b := range data {
		curByte := reflectByte(uint16(b))
		crc ^= curByte << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc = crc << 1
			}
		}
	}

	crc = reflectUint16(crc)
	return []byte{byte(crc & 0xFF), byte(crc >> 8)}
}

// VerifyCRCA compares the calculated CRC-A checksum of the provided data
// against a received CRC checksum. This is typically used to validate
// the integrity of data received from NFC tags. The function returns true
// if the calculated CRC matches the received CRC, indicating that the data
// is likely intact and uncorrupted.
func VerifyCRCA(data []byte, crc []byte) bool {
	calculatedCRC := CalculateCRCA(data)
	return bytes.Equal(calculatedCRC, crc)
}

// reflectByte reverses the order of the last 8 bits in x.
func reflectByte(x uint16) uint16 {
	reflection := uint16(0)
	for i := uint16(0); i < 8; i++ {
		if x&(1<<i) != 0 {
			reflection |= 1 << (7 - i)
		}
	}
	return reflection
}

// reflectUint16 reverses the order of the last 16 bits in x.
func reflectUint16(x uint16) uint16 {
	reflection := uint16(0)
	for i := uint16(0); i < 16; i++ {
		if x&(1<<i) != 0 {
			reflection |= 1 << (15 - i)
		}
	}
	return reflection
}
