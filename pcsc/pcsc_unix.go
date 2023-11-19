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
