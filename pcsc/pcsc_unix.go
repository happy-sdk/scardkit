//go:build !windows && !darwin
// +build !windows,!darwin

// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package pcsc

// #cgo pkg-config: libpcsclite
// #include <stdlib.h>
// #include <winscard.h>
import "C"

import (
	"strings"
	"time"
	"unsafe"
)

func sCardEstablishContext(scope ScardScope, reserved1, reserved2 lpcvoid) (hContext, returnValue) {
	var hctx C.SCARDCONTEXT
	rv := C.SCardEstablishContext(C.DWORD(scope), C.LPCVOID(reserved1), C.LPCVOID(reserved2), &hctx)
	return hContext(hctx), returnValue(rv)
}

func sCardReleaseContext(hctx hContext) returnValue {
	rv := C.SCardReleaseContext(C.SCARDCONTEXT(hctx))
	return returnValue(rv)
}

func sCardIsValidContext(hctx hContext) returnValue {
	rv := C.SCardIsValidContext(C.SCARDCONTEXT(hctx))
	return returnValue(rv)
}

// sCardListReaders gets the list of smart card readers available and returns them as a slice of strings.
func sCardListReaders(hctx hContext) ([]string, returnValue) {
	var readersLen C.DWORD
	// First, get the length of the buffer needed to store the readers
	res := C.SCardListReaders(C.SCARDCONTEXT(hctx), nil, nil, &readersLen)
	if res != C.SCARD_S_SUCCESS {
		return nil, returnValue(res)
	}

	// Allocate a buffer of the required size
	readers := make([]byte, readersLen)
	res = C.SCardListReaders(C.SCARDCONTEXT(hctx), nil, (*C.char)(unsafe.Pointer(&readers[0])), &readersLen)
	if res != C.SCARD_S_SUCCESS {
		return nil, returnValue(res)
	}

	// Convert the buffer into a Go string, then split it into individual reader names
	readerString := string(readers[:readersLen-1]) // Exclude the last null character
	readerString = strings.TrimRight(readerString, "\x00")
	readerList := strings.Split(readerString, "\x00")
	return readerList, returnValue(res)
}

// Convert Go ReaderState to C SCARD_READERSTATE
func toCReaderState(readerState *ReaderState) (C.SCARD_READERSTATE, func()) {
	cReader := C.CString(readerState.Reader)
	cUserData := readerState.userData

	cleanup := func() {
		C.free(unsafe.Pointer(cReader))
	}

	return C.SCARD_READERSTATE{
		szReader:       cReader,
		pvUserData:     cUserData,
		dwCurrentState: C.DWORD(readerState.CurrentState),
		dwEventState:   C.DWORD(readerState.EventState),
		rgbAtr:         *(*[MaxAtrSize]C.uchar)(unsafe.Pointer(&readerState.Atr)),
	}, cleanup
}

func sCardGetStatusChange(hctx hContext, states []ReaderState, timeout time.Duration) returnValue {
	cReaderStates := make([]C.SCARD_READERSTATE, len(states))
	cleanups := make([]func(), len(states))

	for i, state := range states {
		cState, cleanup := toCReaderState(&state)
		cReaderStates[i] = cState
		cleanups[i] = cleanup
	}

	t := timeout.Milliseconds()
	if timeout < 0 {
		t = Infinite
	}

	res := C.SCardGetStatusChange(C.SCARDCONTEXT(hctx), C.DWORD(t), &cReaderStates[0], C.DWORD(len(states)))

	// Cleanup C strings
	for _, cleanup := range cleanups {
		cleanup()
	}

	if res != C.SCARD_S_SUCCESS {
		return returnValue(res)
	}

	// Update the Go slice with the results
	for i := range states {
		states[i].CurrentState = ScardState(cReaderStates[i].dwCurrentState)
		states[i].EventState = ScardState(cReaderStates[i].dwEventState)
		copy(states[i].Atr[:], C.GoBytes(unsafe.Pointer(&cReaderStates[i].rgbAtr), MaxAtrSize))
	}

	return returnValue(res)
}

func sCardCancel(hctx hContext) returnValue {
	rv := C.SCardCancel(C.SCARDCONTEXT(hctx))
	return returnValue(rv)
}
