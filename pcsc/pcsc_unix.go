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

func sCardFreeMemory(hctx hContext, pvMem unsafe.Pointer) returnValue {
	rv := C.SCardFreeMemory(C.SCARDCONTEXT(hctx), C.LPCVOID(pvMem))
	return returnValue(rv)
}

func sCardConnect(hctx hContext, reader string, mode ScardSharedMode, protocol ScardProtocol) (uintptr, ScardProtocol, returnValue) {
	var handle C.SCARDHANDLE
	var aprotocol C.DWORD
	creader := C.CString(reader)

	rv := C.SCardConnect(C.SCARDCONTEXT(hctx), C.LPCSTR(creader), C.DWORD(mode), C.DWORD(protocol), &handle, &aprotocol)
	return uintptr(handle), ScardProtocol(aprotocol), returnValue(rv)
}

func sCardDisconnect(handle uintptr, d ScardDisposition) returnValue {
	rv := C.SCardDisconnect(C.SCARDHANDLE(handle), C.DWORD(d))
	return returnValue(rv)
}

// SCardStatus wraps the SCardStatus function
func sCardStatus(cardHandle uintptr) (string, ScardCardState, ScardProtocol, []byte, returnValue) {
	var readerName [MaxReadername]C.char
	var readerLen C.ulong = C.DWORD(MaxReadername)
	var state, protocol C.DWORD
	var atr [MaxAtrSize]C.uchar
	var atrLen C.ulong = C.DWORD(MaxAtrSize)

	rv := C.SCardStatus(
		C.SCARDHANDLE(cardHandle),
		&readerName[0],
		&readerLen,
		&state,
		&protocol,
		&atr[0],
		&atrLen,
	)

	rReaderName := strings.TrimRight(C.GoStringN(&readerName[0], C.int(readerLen)), "\x00")
	rState := ScardCardState(state)
	rProtocol := ScardProtocol(protocol)
	rAtr := C.GoBytes(unsafe.Pointer(&atr[0]), C.int(atrLen))
	return rReaderName, rState, rProtocol, rAtr, returnValue(rv)
}

func sCardBeginTransaction(cardHandle uintptr) returnValue {
	rv := C.SCardBeginTransaction(C.SCARDHANDLE(cardHandle))
	return returnValue(rv)
}

func sCardEndTransaction(cardHandle uintptr, d ScardDisposition) returnValue {
	rv := C.SCardEndTransaction(C.SCARDHANDLE(cardHandle), C.DWORD(d))
	return returnValue(rv)
}

func sCardTransmit(cardHandle uintptr, proto ScardProtocol, cmd []byte, rsp []byte) (int, returnValue) {
	var sendPci *C.SCARD_IO_REQUEST
	switch proto {
	case ScardProtocolT0:
		sendPci = C.SCARD_PCI_T0
	case ScardProtocolT1:
		sendPci = C.SCARD_PCI_T1
	case ScardProtocolRaw:
		sendPci = C.SCARD_PCI_RAW
	// Add other protocols if needed
	default:
		return 0, returnValue(C.SCARD_E_PROTO_MISMATCH)
	}

	var recvPci C.SCARD_IO_REQUEST
	var recvLen C.DWORD = C.DWORD(len(rsp))

	rv := C.SCardTransmit(
		C.SCARDHANDLE(cardHandle),
		sendPci,
		(*C.uchar)(unsafe.Pointer(&cmd[0])),
		C.DWORD(len(cmd)),
		&recvPci,
		(*C.uchar)(unsafe.Pointer(&rsp[0])),
		&recvLen,
	)

	return int(recvLen), returnValue(rv)
}
