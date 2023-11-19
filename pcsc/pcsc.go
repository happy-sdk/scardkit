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
	"strings"
	"sync"
	"time"
	"unsafe"
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
		return fmt.Sprintf("Unknown CardState (0x%X)", cs)
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

type CardStatus struct {
	Reader   string
	State    ScardCardState
	Protocol ScardProtocol
	Atr      []byte
}
