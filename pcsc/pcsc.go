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
