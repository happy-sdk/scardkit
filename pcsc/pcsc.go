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

func rvToError(rv returnValue) error {
	if err, known := cerrors[rv]; known {
		if errors.Is(err, ErrScardSSuccess) {
			return nil
		}
		return err
	}
	return fmt.Errorf("%w: unknown return value %X", ErrScardFUnknownError, rv)
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
