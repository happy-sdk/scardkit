// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

// Package pcsc offers a Go interface to PCSC lite, enabling direct interaction with smart cards.
// This package leverages CGO for card reader connectivity and provides a streamlined approach
// for handling smart card communication and data transmission in NFC applications.
package pcsc

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/happy-sdk/nfcsdk/internal/helpers"
)

// Error maps PCSC Lite C errorCodes
type Error struct {
	code     returnValue
	name     string
	messages string
	err      error
}

func (e Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s: %s", e.name, e.messages, e.err.Error())
	}
	return fmt.Sprintf("%s: %s", e.name, e.messages)
}

// Unwrap provides the ability to retrieve the underlying error, if any
func (e Error) Unwrap() error {
	return e.err
}

type HContext struct {
	hContext hContext
}

func (hctx *HContext) Release() error {
	return SCardReleaseContext(hctx)
}

func (hctx *HContext) IsValid() error {
	return SCardIsValidContext(hctx)
}

func (hctx *HContext) ListReaders() ([]string, error) {
	return SCardListReaders(hctx)
}

func (hctx *HContext) GetStatusChange(states []ReaderState, timeout time.Duration) error {
	return SCardGetStatusChange(hctx, states, timeout)
}

func (hctx *HContext) Cancel() error {
	return SCardCancel(hctx)
}

func (hctx *HContext) Connect(reader string, mode ScardSharedMode, protocol ScardProtocol) (*Card, error) {
	return SCardConnect(hctx, reader, mode, protocol)
}

func SCardEstablishContext(scope ScardScope) (*HContext, error) {
	hctx, rv := sCardEstablishContext(scope, nil, nil)
	if err := rvToError(rv); err != nil {
		return nil, err
	}
	return &HContext{hContext: hctx}, nil
}

func SCardReleaseContext(hctx *HContext) error {
	return rvToError(sCardReleaseContext(hctx.hContext))
}

func SCardIsValidContext(hctx *HContext) error {
	return rvToError(sCardIsValidContext(hctx.hContext))
}

func SCardListReaders(hctx *HContext) ([]string, error) {
	readers, rv := sCardListReaders(hctx.hContext)
	if err := rvToError(rv); err != nil {
		return nil, err
	}
	return readers, nil
}

func SCardGetStatusChange(hctx *HContext, states []ReaderState, timeout time.Duration) error {
	rv := sCardGetStatusChange(hctx.hContext, states, timeout)
	if err := rvToError(rv); err != nil {
		return err
	}
	return nil
}

func SCardCancel(hctx *HContext) error {
	return rvToError(sCardCancel(hctx.hContext))
}

func SCardConnect(hctx *HContext, reader string, mode ScardSharedMode, protocol ScardProtocol) (*Card, error) {
	handle, aprotocol, rv := sCardConnect(hctx.hContext, reader, mode, protocol)
	if err := rvToError(rv); err != nil {
		return nil, err
	}

	return &Card{handle: handle, protocol: aprotocol}, nil
}

func SCardStatus(cardHandle uintptr) (CardStatus, error) {
	reader, state, protocol, atr, rv := sCardStatus(cardHandle)
	if err := rvToError(rv); err != nil {
		return CardStatus{}, err
	}
	return CardStatus{
		Reader:   reader,
		State:    state,
		Protocol: protocol,
		Atr:      atr,
	}, nil
}

func rvToError(rv returnValue) error {
	if err, known := cerrors[rv]; known {
		if errors.Is(err, ErrScardSSuccess) {
			return nil
		}
		return err
	}
	return fmt.Errorf("%w: unknown return value %X", ErrScardFUnknownError, rv)
}

type ReaderState struct {
	userData     unsafe.Pointer // ignore for now
	Reader       string
	CurrentState ScardState
	EventState   ScardState
	Atr          [MaxAtrSize]byte
}

func (s ScardState) String() string {
	var states []string

	if s&ScardStateUnaware != ScardStateUnaware {
		states = append(states, "Unaware")
	}
	if s&ScardStateIgnore != 0 {
		states = append(states, "Ignore")
	}
	if s&ScardStateChanged != 0 {
		states = append(states, "Changed")
	}
	if s&ScardStateUnknown != 0 {
		states = append(states, "Unknown")
	}
	if s&ScardStateUnavailable != 0 {
		states = append(states, "Unavailable")
	}
	if s&ScardStateEmpty != 0 {
		states = append(states, "Empty")
	}
	if s&ScardStatePresent != 0 {
		states = append(states, "Present")
	}
	if s&ScardStateAtrMatch != 0 {
		states = append(states, "AtrMatch")
	}
	if s&ScardStateExclusive != 0 {
		states = append(states, "Exclusive")
	}
	if s&ScardStateInUse != 0 {
		states = append(states, "InUse")
	}
	if s&ScardStateMute != 0 {
		states = append(states, "Mute")
	}
	if s&ScardStateUnpowered != 0 {
		states = append(states, "Unpowered")
	}

	return strings.Join(states, ", ")
}

func (p ScardProtocol) String() string {
	protocols := []string{}

	if p&ScardProtocolUndefined != 0 {
		protocols = append(protocols, "Undefined")
	}
	if p&ScardProtocolT0 != 0 {
		protocols = append(protocols, "T0")
	}
	if p&ScardProtocolT1 != 0 {
		protocols = append(protocols, "T1")
	}
	if p&ScardProtocolRaw != 0 {
		protocols = append(protocols, "Raw")
	}
	if p&ScardProtocolT15 != 0 {
		protocols = append(protocols, "T15")
	}
	if p&ScardProtocolAny != 0 {
		protocols = append(protocols, "Any")
	}

	if len(protocols) == 0 {
		return "No Protocol"
	}

	return strings.Join(protocols, ", ")
}

func (cs ScardCardState) String() string {
	var parts []string

	if cs&StateUnknown != 0 {
		parts = append(parts, "Unknown")
	}
	if cs&StateAbsent != 0 {
		parts = append(parts, "Absent")
	}
	if cs&StatePresent != 0 {
		parts = append(parts, "Present")
	}
	if cs&StateSwallowed != 0 {
		parts = append(parts, "Swallowed")
	}
	if cs&StatePowered != 0 {
		parts = append(parts, "Powered")
	}
	if cs&StateNegotiable != 0 {
		parts = append(parts, "Negotiable")
	}
	if cs&StateSpecific != 0 {
		parts = append(parts, "Specific")
	}

	if len(parts) == 0 {
		return fmt.Sprintf("Unknown CardState (0x%X)", uint64(cs))
	}
	return strings.Join(parts, ", ")
}

type Card struct {
	mu       sync.Mutex
	handle   uintptr
	protocol ScardProtocol
	status   CardStatus
}

func (c *Card) Protocol() ScardProtocol {
	return c.protocol
}

func (c *Card) Disconnect(d ScardDisposition) error {
	rv := sCardDisconnect(c.handle, d)
	if err := rvToError(rv); err != nil {
		return err
	}
	return nil
}

func (c *Card) CurrentStatus() CardStatus {
	c.mu.Lock()
	defer c.mu.Unlock()
	status := c.status
	return status
}

func (c *Card) RefreshStatus() error {
	status, err := SCardStatus(c.handle)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status = status
	return nil
}

func (card *Card) Transmit(cmd *Command) (CardResponse, error) {
	if cmd == nil {
		return CardResponse{}, errors.New("no command provided to card.Transmit")
	}
	if cmd.err != nil {
		return CardResponse{}, fmt.Errorf("%s: %w", cmd.name, cmd.err)
	}
	// Define a response buffer with a reasonable initial size for provided command.
	rsp, err := cmd.newRespBuffer()
	if err != nil {
		return CardResponse{}, err
	}
	// rsp := make([]byte, MaxBufferSizeExtended)

	// Call sCardTransmit with the card's handle, protocol, command, and response buffer.
	recvLen, rv := sCardTransmit(card.handle, card.protocol, cmd.Bytes(), rsp)
	if err := rvToError(rv); err != nil {
		return CardResponse{}, fmt.Errorf("%w: command %s", err, cmd.Name())
	}

	// Trim the response slice to the actual length of the response.
	rsp = rsp[:recvLen]
	if len(rsp) < 2 {
		return CardResponse{}, fmt.Errorf("%w: no sw1 and sw2 returned for %s", ErrResponse, cmd.Name())
	}
	pllen := len(rsp) - 2 // payload length

	status := rsp[pllen:]
	response := CardResponse{}
	response.status = SW1SW2(status)
	var payload []byte
	if cmd.postProcessFunc == nil {
		payload = rsp[:pllen]
	} else {
		payload, err = cmd.postProcessFunc(rsp[:pllen])
		if err != nil {
			return CardResponse{}, fmt.Errorf("%s: %w", cmd.Name(), err)
		}
	}
	response.payload = payload
	if !response.Status().Success() {
		return CardResponse{}, fmt.Errorf("%s: %w", cmd.Name(), response.Status().error())
	}

	return response, nil
}

func (card *Card) BeginTransaction() error {
	rv := sCardBeginTransaction(card.handle)
	if err := rvToError(rv); err != nil {
		return err
	}
	return nil
}

func (card *Card) EndTransaction(d ScardDisposition) error {
	rv := sCardEndTransaction(card.handle, d)
	if err := rvToError(rv); err != nil {
		return err
	}
	return nil
}

type CardStatus struct {
	Reader   string
	State    ScardCardState
	Protocol ScardProtocol
	Atr      []byte
}

const ZeroByte = 0x00

var success = SW1SW2{0x90, 0x00}

var (
	ErrResponse = errors.New("error")
)

var responses = map[SW1SW2]string{
	success:          "success",
	{0x62, ZeroByte}: "no information given",
	{0x62, 0x81}:     "returned data may be corrupted",
	{0x62, 0x82}:     "the end of the file has been reached before the end of reading",
	{0x62, 0x83}:     "invalid DF",
	{0x62, 0x84}:     "selected file is not valid - file descriptor error",
	{0x63, ZeroByte}: "authentification failed - invalid secret code or forbidden value",
	{0x63, 0x81}:     "file filled up by the last write",
	{0x65, 0x03}:     "memory failure: EEPROM read/write or hardware problem",
	{0x65, 0x81}:     "write problem / memory failure / unknown mode",
	{0x67, ZeroByte}: "incorrect length or address range",
	{0x68, ZeroByte}: "the request function is not supported by the card",
	{0x68, 0x81}:     "logical channel not supported",
	{0x68, 0x82}:     "secure messaging not supported ",
	{0x69, ZeroByte}: "no successful transaction executed during session",
	{0x69, 0x81}:     "cannot select indicated file, command not compatible with file organization",
	{0x69, 0x82}:     "access conditions not fulfilled",
	{0x69, 0x83}:     "secret code locked",
	{0x69, 0x84}:     "referenced data invalidated",
	{0x69, 0x85}:     "no currently selected EF, no command to monitor / no Transaction Manager File",
	{0x69, 0x86}:     "command not allowed (no current EF)",
	{0x69, 0x87}:     "expected SM data objects missing",
	{0x69, 0x88}:     "SM data objects incorrect",
	{0x6A, ZeroByte}: "bytes P1 and/or P2 are incorrect.",
	{0x6A, 0x80}:     "the parameters in the data field are incorrect",
	{0x6A, 0x81}:     "card is blocked or command not supported",
	{0x6A, 0x82}:     "file not found",
	{0x6A, 0x83}:     "record not found",
	{0x6A, 0x84}:     "there is insufficient memory space in record or file",
	{0x6A, 0x85}:     "Lc inconsistent with TLV structure",
	{0x6A, 0x86}:     "incorrect parameters P1-P2",
	{0x6A, 0x87}:     "the P3 value is not consistent with the P1 and P2 values",
	{0x6A, 0x88}:     "referenced data not found",
	{0x6B, ZeroByte}: "incorrect reference; illegal address; invalid P1 or P2 parameter",
	{0x6D, ZeroByte}: "command not allowed. invalid instruction byte (INS)",
	{0x6E, ZeroByte}: "incorrect application (CLA parameter of a command)",
	{0x6F, ZeroByte}: "checking error",
	{0x91, ZeroByte}: "purse balance error cannot perform transaction",
	{0x91, 0x02}:     "purse balance error",
	{0x92, 0x02}:     "write problem / memory failure",
	{0x92, 0x40}:     "error, memory problem",
	{0x94, 0x04}:     "purse selection error or invalid purse",
	{0x94, 0x06}:     "invalid purse detected during the replacement debit step",
	{0x94, 0x08}:     "key file selection error",
	{0x94, ZeroByte}: "security warning",
	{0x94, 0x04}:     "access authorization not fulfilled",
	{0x94, 0x06}:     "access authorization in Debit not fulfilled for the replacement debit step",
	{0x94, 0x20}:     "no temporary transaction key established",
	{0x94, 0x34}:     "update SSD order sequence not respected",
}

type CardResponse struct {
	status  SW1SW2
	payload []byte
}

func (r CardResponse) Payload() []byte {
	return r.payload
}

func (r CardResponse) Status() SW1SW2 {
	return r.status
}
func (r CardResponse) SW1() byte {
	return r.status.SW1()
}
func (r CardResponse) SW2() byte {
	return r.status.SW2()
}

func (r CardResponse) LogAttr() slog.Attr {
	return slog.Group("",
		slog.String("status", r.status.String()),
		slog.String("sw1", fmt.Sprintf("%02X", r.SW1())),
		slog.String("sw2", fmt.Sprintf("%02X", r.SW2())),
		slog.Int("payload.len", len(r.payload)),
	)
}

func (r CardResponse) String() string {
	return helpers.FormatByteSlice(r.payload)
}

type SW1SW2 [2]byte

func (rs SW1SW2) error() error {
	if rs.Success() {
		return nil
	}
	return fmt.Errorf("%w: %s", ErrResponse, rs.String())
}

func (rs SW1SW2) SW1() byte {
	return rs[0]
}

func (rs SW1SW2) SW2() byte {
	return rs[1]
}

func (rs SW1SW2) Success() bool {
	return rs == success || rs[0] == 0x9F
}

func (rs SW1SW2) String() string {
	if msg, ok := responses[rs]; ok {
		return msg
	}
	// handle other cases
	switch rs[0] {
	case 0x61:
		return fmt.Sprintf("(%d) response bytes still available", rs[1])
	case 0x63:
		return fmt.Sprintf("command response code(%d)", rs[1])
	case 0x67:
		return fmt.Sprintf("incorrect parameter P3 (ISO code) %d", rs[1])
	case 0x6C:
		return fmt.Sprintf("incorrect P3 length (%d) or response buffer to small", rs[1])
	case 0x92:
		return fmt.Sprintf("memory error %d", rs[1])
	case 0x94:
		return fmt.Sprintf("file error %d", rs[1])
	case 0x98:
		return fmt.Sprintf("security error %d", rs[1])
	case 0x9F:
		return fmt.Sprintf("success %d bytes of data available to be read via Get_Response", rs[1])
	default:
		return fmt.Sprintf("unknown response status: %X %X", rs[0], rs[1])
	}
}

type Command struct {
	postProcessFunc func([]byte) ([]byte, error)
	custom          bool
	name            string
	cla             byte   // Instruction class
	ins             byte   // Instruction code
	p1              byte   // Instruction parameter 1 for the command
	p2              byte   // Instruction parameter 2 for the command
	lc              []byte // Encodes the number (Nc) of bytes of command data to follow
	le              []byte // Encodes the maximum number (Ne) of response bytes expected
	payload         []byte // payload
	err             error
}

func NewCustomCmd(payload []byte) *Command {
	return &Command{
		custom:  true,
		payload: payload,
		name:    "CUSTOM",
	}
}

func NewCmd(cla, ins, p1, p2 byte) *Command {
	cmd := &Command{
		name: "RAW",
		cla:  cla,
		ins:  ins,
		p1:   p1,
		p2:   p2,
	}

	return cmd
}

// Name returns command name
func (c *Command) Name() string {
	return c.name
}
func (c *Command) SetName(name string) {
	c.name = name
}

func (c *Command) SetLe(le []byte) {
	if l := len(le); l > 3 {
		c.addErr(fmt.Errorf("(Le) invalid length %d [%s...]", l, helpers.FormatByteSlice(le[0:3])))
		return
	}
	c.le = le
}

func (c *Command) addErr(err error) {
	c.err = errors.Join(c.err, err)
}

// Bytes returns command byte slice with optional payload if present
func (c *Command) Bytes() []byte {
	if c.custom {
		return c.payload
	}
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
		cmd = append(cmd, ZeroByte)
	}
	return cmd
}

// String returns string representation of currect command
func (c *Command) String() string {
	str := c.name + " ["
	cmd := c.Bytes()
	if len(cmd) < 10 {
		str += helpers.FormatByteSlice(cmd)
	} else {
		str += helpers.FormatByteSlice(cmd[:10])
		str += "..."
	}
	str += "]"
	return str
}

func (c *Command) SetPostProcessFunc(f func([]byte) ([]byte, error)) {
	if c.postProcessFunc != nil {
		c.addErr(errors.New("SetPostProcessFunc can only used once for command"))
		return
	}
	c.postProcessFunc = f
}

func (c *Command) newRespBuffer() ([]byte, error) {
	var buf []byte
	if c.le == nil {
		buf = make([]byte, MaxBufferSizeExtended)
		return buf, nil
	}

	switch len(c.le) {
	case 0:
		// No Le field, but still need room for SW1 and SW2
		buf = make([]byte, 2)

	case 1:
		// Short Le field
		le := int(c.le[0])
		if le == 0 {
			le = 256
		}
		buf = make([]byte, le+2) // Add 2 for SW1 and SW2

	case 2:
		// Extended Le field
		le := int(c.le[0])<<8 + int(c.le[1])
		if le == 0 {
			le = 65536
		}
		buf = make([]byte, le+2) // Add 2 for SW1 and SW2

	case 3:
		// Special case, first byte should be 0, next two bytes are Le
		if c.le[0] != 0 {
			return nil, fmt.Errorf("invalid Le field")
		}
		le := int(c.le[1])<<8 + int(c.le[2])
		buf = make([]byte, le+2) // Add 2 for SW1 and SW2

	default:
		return nil, fmt.Errorf("invalid Le field length")
	}

	return buf, nil
}
