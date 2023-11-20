// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

// Package ntag offers specialized support for interacting with NFC Forum Type 2 Tags,
// specifically the NXP NTAG series. This package simplifies operations such as reading from
// and writing to NTAG memory, handling tag authentication, and managing tag-specific features.
// It aims to provide a straightforward and efficient interface for working with NTAGs,
// ensuring compatibility and ease of use in NFC applications.
package ntag

import (
	"errors"
	"fmt"

	"github.com/happy-sdk/nfcsdk"
	"github.com/happy-sdk/nfcsdk/pcsc"
)

var (
	Error         = errors.New("ntag")
	ErrUnknownCMD = fmt.Errorf("%w: unknown command", Error)
)

// Commands
// | Command[1]            | ISO/IEC 14443      | NFC FORUM          | Command code (hexadecimal) |
// |-----------------------|--------------------|--------------------|----------------------------|
// | Request               | REQA               | SENS_REQ           | 26h (7 bit)                |
// | Wake-up               | WUPA               | ALL_REQ            | 52h (7 bit)                |
// | Anticollision CL1     | Anticollision CL1  | SDD_REQ CL1        | 93h 20h                    |
// | Select CL1            | Select CL1         | SEL_REQ CL1        | 93h 70h                    |
// | Anticollision CL2     | Anticollision CL2  | SDD_REQ CL2        | 95h 20h                    |
// | Select CL2            | Select CL2         | SEL_REQ CL2        | 95h 70h                    |
// | Halt                  | HLTA               | SLP_REQ            | 50h 00h                    |
// | GET_VERSION[2]        | -                  | -                  | 60h                        |
// | READ                  | -                  | READ               | 30h                        |
// | FAST_READ[2]          | -                  | -                  | 3Ah                        |
// | WRITE                 | -                  | WRITE              | A2h                        |
// | COMP_WRITE            | -                  | -                  | A0h                        |
// | READ_CNT[2]           | -                  | -                  | 39h                        |
// | PWD_AUTH[2]           | -                  | -                  | 1Bh                        |
// | READ_SIG[2]           | -                  | -                  | 3Ch                        |

type CMD int

const (
	CmdRAW            CMD = iota // Raw constructed command
	CmdRequest                   // REQA, SENS_REQ, 26h (7 bit)
	CmdWakeup                    // WUPA, ALL_REQ, 52h (7 bit)
	CmdAnticollision1            // Anticollision CL1, SDD_REQ CL1, 93h 20h
	CmdAnticollision2            // Anticollision CL2, SDD_REQ CL2, 95h 20h
	CmdSelect1                   // Select CL1, SEL_REQ CL1, 93h 70h
	CmdSelect2                   // Select CL2, SEL_REQ CL2, 95h 70h
	CmdHalt                      // HLTA, SLP_REQ, 50h 00h
	CmdGetVersion                // GET_VERSION, 60h
	CmdRead                      // READ, 30h
	CmdFastRead                  // FAST_READ, 3Ah
	CmdWrite                     // WRITE, A2h
	CmdCompWrite                 // OMP_WRITE, -A0h
	CmdReadCnt                   // READ_CNT, -39h
	CmdPwdAuth                   // PWD_AUTH, -1Bh
	CmdReadSig                   // READ_SIG, -3Ch

	// other commands
	CmdGetUID
)

var cmdNames = map[CMD]string{
	CmdRAW:            "RAW",
	CmdRequest:        "REQA, SENS_REQ",
	CmdWakeup:         "WUPA, ALL_REQ",
	CmdAnticollision1: "Anticollision CL1, SDD_REQ CL1",
	CmdAnticollision2: "Anticollision CL2, SDD_REQ CL2",
	CmdSelect1:        "Select CL1, SEL_REQ CL1",
	CmdSelect2:        "Select CL2, SEL_REQ CL2",
	CmdHalt:           "HLTA, SLP_REQ",
	CmdGetVersion:     "GET_VERSION",
	CmdRead:           "READ",
	CmdFastRead:       "FAST_READ",
	CmdWrite:          "WRITE",
	CmdCompWrite:      "COMP_WRITE",
	CmdReadCnt:        "READ_CNT",
	CmdPwdAuth:        "PWD_AUTH",
	CmdReadSig:        "READ_SIG",
	CmdGetUID:         "GET_UID",
}

var cmdBytes = map[CMD][2]byte{
	CmdRAW:            {0x00, pcsc.ZeroByte},
	CmdRequest:        {0x26, pcsc.ZeroByte},
	CmdWakeup:         {0x52, pcsc.ZeroByte},
	CmdAnticollision1: {0x93, 0x20},
	CmdAnticollision2: {0x95, 0x20},
	CmdSelect1:        {0x93, 0x70},
	CmdSelect2:        {0x95, 0x70},
	CmdHalt:           {0x50, pcsc.ZeroByte},
	CmdGetVersion:     {0x60, pcsc.ZeroByte},
	CmdRead:           {0x30, pcsc.ZeroByte},
	CmdFastRead:       {0x3A, pcsc.ZeroByte},
	CmdWrite:          {0xA2, pcsc.ZeroByte},
	CmdCompWrite:      {0xA0, pcsc.ZeroByte},
	CmdReadCnt:        {0x39, pcsc.ZeroByte},
	CmdPwdAuth:        {0x1B, pcsc.ZeroByte},
	CmdReadSig:        {0x3C, pcsc.ZeroByte},
	CmdGetUID:         {0xFF, 0xCA}, // Class and Instruction byte for the Get UID command
}

type Command struct {
	cmd     CMD
	name    string
	cla     byte   // Instruction class
	ins     byte   // Instruction code
	p1      byte   // Instruction parameter 1 for the command
	p2      byte   // Instruction parameter 2 for the command
	lc      []byte // Encodes the number (Nc) of bytes of command data to follow
	le      []byte // Encodes the maximum number (Ne) of response bytes expected
	payload []byte // payload
}

// Cmd creates command byte slice with optional payload
// The function can accept a nil or empty payload and
// it will handle such cases gracefully.
func Cmd(c CMD) (*Command, error) {
	cb, ok := cmdBytes[c]
	if !ok {
		return nil, fmt.Errorf("%w: %d", ErrUnknownCMD, c)
	}
	cmd, err := RawCmd(cb[0], cb[1], pcsc.ZeroByte, pcsc.ZeroByte)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func RawCmd(cla, ins, p1, p2 byte) (*Command, error) {
	cmd := &Command{
		name: cmdNames[CmdRAW],
		cmd:  CmdRAW,
		cla:  cla,
		ins:  ins,
		p1:   p1,
		p2:   p2,
	}
	for c, b := range cmdBytes {
		if cmd.cla == b[0] && cmd.ins == b[1] {
			cmd.cmd = c
			if name, ok := cmdNames[c]; ok {
				cmd.name = name
			}
		}
	}
	return cmd, nil
}

// Name returns command name
func (c *Command) Name() string {
	return c.name
}
func (c *Command) SetLE(le []byte) {
	c.le = le
}

// Bytes returns command byte slice with optional payload if present
func (c *Command) Bytes() []byte {
	// Start with the CLA, INS, P1, P2
	cmd := []byte{c.cla, c.ins, c.p1, c.p2}

	// If lc is not nil and payload is present, append the length and the payload
	if c.lc != nil && len(c.payload) > 0 {
		cmd = append(cmd, c.lc...)
		cmd = append(cmd, c.payload...)
	}
	// If le is not nil, append it
	if c.le != nil {
		cmd = append(cmd, c.le...)
	} else {
		cmd = append(cmd, pcsc.ZeroByte)
	}
	return cmd
}

// String returns string representation of currect command
func (c *Command) String() string {
	str := c.name + " ["
	str += nfcsdk.FormatByteSlice(c.Bytes())
	str += "]"
	return str
}

// Transmit command to card and retrun the response
func (c *Command) Transmit(card nfcsdk.Card) (pcsc.CardResponse, error) {
	if card == nil {
		return pcsc.CardResponse{}, fmt.Errorf("%w: no card provided", Error)
	}
	return card.Transmit(c.Bytes())
}
