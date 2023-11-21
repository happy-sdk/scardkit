// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package pcsc

import (
	"fmt"

	"github.com/happy-sdk/nfcsdk/internal/helpers"
)

// Command struct represents a structure for composing commands that are sent to smart cards.
// This struct encapsulates various elements of a typical smart card command, including the
// instruction class and code, parameters, data payload, and expected response size.
type Command struct {
	// postProcessFunc is an optional function applied to the response from the card.
	// It takes a byte slice as input and returns a processed byte slice and an error if any.
	postProcessFunc func([]byte) ([]byte, error)

	// raw indicates whether the command is a predefined standard command or a custom one.
	// Standard commands follow established smart card protocols, while custom commands
	// might be specific to a particular card or application.
	raw bool

	// name is the descriptive name of the command. This can either be a predefined
	// name from a specification or standard (e.g., "READ", "WRITE", "GET_VERSION" from
	// the NTAG specification) or a custom name assigned by the application.
	name string

	// cla represents the Instruction class of the command, defining the type of command being sent.
	cla byte

	// ins is the Instruction code that determines the specific command or operation being performed.
	ins byte

	// p1 is the first parameter of the instruction.
	p1 byte

	// p2 is the second parameter of the instruction.
	p2 byte

	// lc encodes the number (Nc) of bytes of command data to follow.
	// It specifies the length of the command payload.
	lc []byte

	// le encodes the maximum number (Ne) of response bytes expected from the card.
	// In standard APDU commands, 'le' should not exceed 3 bytes; exceeding this limit
	// results in an error when the Command is used. The SetLe function sets this field
	// and ensures compliance with the APDU command size limits. For custom commands created
	// with e.g. NewRawCmd, 'le' determines the response buffer size in PC/SC but is not included
	// in the .Bytes() output of the command.
	le []byte

	// payload contains the actual data sent to the card as part of the command.
	payload []byte

	// err captures any errors that occur during the construction of the command.
	// It allows error handling to be deferred until a command's execution.
	err error
}

// NewRawCmd creates a new Command struct with raw mode enabled.
// It initializes the command with the given byte slice as payload and sets the name to "RAW".
func NewRawCmd(cmd []byte) *Command {
	return &Command{
		raw:     true,
		payload: cmd,
		name:    "RAW",
	}
}

// NewCmd creates a new Command struct for a known command type.
// It sets the command fields CLA, INS, P1, and P2 with the provided byte values and
// determines the command's name using the lookupKnownCommand function based on CLA and INS.
func NewCmd(cla, ins, p1, p2 byte) *Command {
	cmd := &Command{
		cla: cla,
		ins: ins,
		p1:  p1,
		p2:  p2,
	}
	cmd.resolveAPDUCommandName()
	return cmd
}

// SetLe sets the expected response length (Le) for the Command using 'le'. Le, represented by a
// byte slice, should not exceed 3 bytes for standard APDU commands. If 'le' is longer than 3 bytes,
// the method adds an error to the Command and exits early when Command is used.
// For custom commands created with NewRawCmd, 'le' is used to determine the response buffer
// size in pcsc, but it's not included in the .Bytes() output.
//
// le []byte: Byte slice for expected response length.
// In custom commands, it defines buffer size, not APDU Le.
func (c *Command) SetLe(le []byte) {
	if l := len(le); l > 3 {
		c.addErr(fmt.Errorf("(Le) invalid length %d [%s...]", l, helpers.FormatByteSlice(le[0:3])))
		return
	}
	c.le = le
}

// resolveAPDUCommandName searches the name of the APDU command and sets c.name is found
func (c *Command) resolveAPDUCommandName() {
	// https://web.archive.org/web/20090630004017/http://cheef.ru/docs/HowTo/APDU.info
	// https://gist.github.com/hemantvallabh/d24d71a933e7319727cd3daa50ad9f2c
	c.name = "UNKNOWN"
}
