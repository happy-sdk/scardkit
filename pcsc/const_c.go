// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

//go:build ignore
// +build ignore

package pcsc

/*
#cgo pkg-config: libpcsclite
#include <PCSC/winscard.h>
#include <PCSC/wintypes.h>
*/
import "C"

type (
	long         C.LONG
	scardContext C.SCARDCONTEXT
	dword        C.DWORD
	hContext     C.SCARDCONTEXT
	lpcvoid      C.LPCVOID
	lpcstr       C.LPCSTR
	lpstr        C.LPSTR
	lpdword      C.LPDWORD
	returnValue  long
)

// Error definitions
var errDef = map[returnValue]Error{
	C.SCARD_S_SUCCESS:             {C.SCARD_S_SUCCESS, "SCARD_S_SUCCESS", "no error encountered"},
	C.SCARD_F_INTERNAL_ERROR:      {C.SCARD_F_INTERNAL_ERROR, "SCARD_F_INTERNAL_ERROR", "internal consistency check failed"},
	C.SCARD_E_CANCELLED:           {C.SCARD_E_CANCELLED, "SCARD_E_CANCELLED", "action cancelled by SCardCancel request"},
	C.SCARD_E_INVALID_HANDLE:      {C.SCARD_E_INVALID_HANDLE, "SCARD_E_INVALID_HANDLE", "supplied handle was invalid"},
	C.SCARD_E_INVALID_PARAMETER:   {C.SCARD_E_INVALID_PARAMETER, "SCARD_E_INVALID_PARAMETER", "parameters could not be properly interpreted"},
	C.SCARD_E_INVALID_TARGET:      {C.SCARD_E_INVALID_TARGET, "SCARD_E_INVALID_TARGET", "registry startup information is missing or invalid"},
	C.SCARD_E_NO_MEMORY:           {C.SCARD_E_NO_MEMORY, "SCARD_E_NO_MEMORY", "not enough memory available to complete command"},
	C.SCARD_F_WAITED_TOO_LONG:     {C.SCARD_F_WAITED_TOO_LONG, "SCARD_F_WAITED_TOO_LONG", "internal consistency timer expired"},
	C.SCARD_E_INSUFFICIENT_BUFFER: {C.SCARD_E_INSUFFICIENT_BUFFER, "SCARD_E_INSUFFICIENT_BUFFER", "data buffer too small for returned data"},
	C.SCARD_E_UNKNOWN_READER:      {C.SCARD_E_UNKNOWN_READER, "SCARD_E_UNKNOWN_READER", "specified reader name not recognized"},
	C.SCARD_E_TIMEOUT:             {C.SCARD_E_TIMEOUT, "SCARD_E_TIMEOUT", "user-specified timeout expired"},
	C.SCARD_E_SHARING_VIOLATION:   {C.SCARD_E_SHARING_VIOLATION, "SCARD_E_SHARING_VIOLATION", "smart card cannot be accessed due to other connections"},
	C.SCARD_E_NO_SMARTCARD:        {C.SCARD_E_NO_SMARTCARD, "SCARD_E_NO_SMARTCARD", "operation requires a smart card, but none is in the device"},
	C.SCARD_E_UNKNOWN_CARD:        {C.SCARD_E_UNKNOWN_CARD, "SCARD_E_UNKNOWN_CARD", "specified smart card name not recognized"},
	C.SCARD_E_CANT_DISPOSE:        {C.SCARD_E_CANT_DISPOSE, "SCARD_E_CANT_DISPOSE", "system could not dispose of the media as requested"},
	C.SCARD_E_PROTO_MISMATCH:      {C.SCARD_E_PROTO_MISMATCH, "SCARD_E_PROTO_MISMATCH", "requested protocols incompatible with card's protocols"},
	C.SCARD_E_NOT_READY:           {C.SCARD_E_NOT_READY, "SCARD_E_NOT_READY", "reader or smart card not ready to accept commands"},
	C.SCARD_E_INVALID_VALUE:       {C.SCARD_E_INVALID_VALUE, "SCARD_E_INVALID_VALUE", "one or more supplied parameter values could not be interpreted"},
	C.SCARD_E_SYSTEM_CANCELLED:    {C.SCARD_E_SYSTEM_CANCELLED, "SCARD_E_SYSTEM_CANCELLED", "action cancelled by the system"},
	C.SCARD_F_COMM_ERROR:          {C.SCARD_F_COMM_ERROR, "SCARD_F_COMM_ERROR", "internal communications error detected"},
	C.SCARD_F_UNKNOWN_ERROR:       {C.SCARD_F_UNKNOWN_ERROR, "SCARD_F_UNKNOWN_ERROR", "internal error detected, but source unknown"},
	C.SCARD_E_INVALID_ATR:         {C.SCARD_E_INVALID_ATR, "SCARD_E_INVALID_ATR", "ATR from registry is not a valid ATR string"},
	C.SCARD_E_NOT_TRANSACTED:      {C.SCARD_E_NOT_TRANSACTED, "SCARD_E_NOT_TRANSACTED", "attempt made to end a non-existent transaction"},
	C.SCARD_E_READER_UNAVAILABLE:  {C.SCARD_E_READER_UNAVAILABLE, "SCARD_E_READER_UNAVAILABLE", "specified reader not currently available"},
	C.SCARD_P_SHUTDOWN:            {C.SCARD_P_SHUTDOWN, "SCARD_P_SHUTDOWN", "operation aborted to allow server application to exit"},
	C.SCARD_E_PCI_TOO_SMALL:       {C.SCARD_E_PCI_TOO_SMALL, "SCARD_E_PCI_TOO_SMALL", "PCI receive buffer was too small"},
	C.SCARD_E_READER_UNSUPPORTED:  {C.SCARD_E_READER_UNSUPPORTED, "SCARD_E_READER_UNSUPPORTED", "reader driver does not meet minimal requirements"},
	C.SCARD_E_DUPLICATE_READER:    {C.SCARD_E_DUPLICATE_READER, "SCARD_E_DUPLICATE_READER", "reader driver did not produce a unique reader name"},
	C.SCARD_E_CARD_UNSUPPORTED:    {C.SCARD_E_CARD_UNSUPPORTED, "SCARD_E_CARD_UNSUPPORTED", "smart card does not meet minimal requirements"},
	C.SCARD_E_NO_SERVICE:          {C.SCARD_E_NO_SERVICE, "SCARD_E_NO_SERVICE", "smart card resource manager is not running"},
	C.SCARD_E_SERVICE_STOPPED:     {C.SCARD_E_SERVICE_STOPPED, "SCARD_E_SERVICE_STOPPED", "smart card resource manager has shut down"},
	C.SCARD_E_UNEXPECTED:          {C.SCARD_E_UNEXPECTED, "SCARD_E_UNEXPECTED", "unexpected card error occurred"},
	// C.SCARD_E_UNSUPPORTED_FEATURE:     {C.SCARD_E_UNSUPPORTED_FEATURE, "SCARD_E_UNSUPPORTED_FEATURE", "smart card does not support the requested feature"},
	C.SCARD_E_ICC_INSTALLATION:        {C.SCARD_E_ICC_INSTALLATION, "SCARD_E_ICC_INSTALLATION", "no primary provider can be found for the smart card"},
	C.SCARD_E_ICC_CREATEORDER:         {C.SCARD_E_ICC_CREATEORDER, "SCARD_E_ICC_CREATEORDER", "requested order of object creation not supported"},
	C.SCARD_E_DIR_NOT_FOUND:           {C.SCARD_E_DIR_NOT_FOUND, "SCARD_E_DIR_NOT_FOUND", "identified directory does not exist on the smart card"},
	C.SCARD_E_FILE_NOT_FOUND:          {C.SCARD_E_FILE_NOT_FOUND, "SCARD_E_FILE_NOT_FOUND", "identified file does not exist on the smart card"},
	C.SCARD_E_NO_DIR:                  {C.SCARD_E_NO_DIR, "SCARD_E_NO_DIR", "supplied path does not represent a smart card directory"},
	C.SCARD_E_NO_FILE:                 {C.SCARD_E_NO_FILE, "SCARD_E_NO_FILE", "supplied path does not represent a smart card file"},
	C.SCARD_E_NO_ACCESS:               {C.SCARD_E_NO_ACCESS, "SCARD_E_NO_ACCESS", "access denied to this file"},
	C.SCARD_E_WRITE_TOO_MANY:          {C.SCARD_E_WRITE_TOO_MANY, "SCARD_E_WRITE_TOO_MANY", "smart card does not have enough memory to store information"},
	C.SCARD_E_BAD_SEEK:                {C.SCARD_E_BAD_SEEK, "SCARD_E_BAD_SEEK", "error trying to set smart card file object pointer"},
	C.SCARD_E_INVALID_CHV:             {C.SCARD_E_INVALID_CHV, "SCARD_E_INVALID_CHV", "supplied PIN is incorrect"},
	C.SCARD_E_UNKNOWN_RES_MNG:         {C.SCARD_E_UNKNOWN_RES_MNG, "SCARD_E_UNKNOWN_RES_MNG", "unrecognized error code from a layered component"},
	C.SCARD_E_NO_SUCH_CERTIFICATE:     {C.SCARD_E_NO_SUCH_CERTIFICATE, "SCARD_E_NO_SUCH_CERTIFICATE", "requested certificate does not exist"},
	C.SCARD_E_CERTIFICATE_UNAVAILABLE: {C.SCARD_E_CERTIFICATE_UNAVAILABLE, "SCARD_E_CERTIFICATE_UNAVAILABLE", "requested certificate could not be obtained"},
	C.SCARD_E_NO_READERS_AVAILABLE:    {C.SCARD_E_NO_READERS_AVAILABLE, "SCARD_E_NO_READERS_AVAILABLE", "cannot find a smart card reader"},
	C.SCARD_E_COMM_DATA_LOST:          {C.SCARD_E_COMM_DATA_LOST, "SCARD_E_COMM_DATA_LOST", "communications error with smart card detected"},
	C.SCARD_E_NO_KEY_CONTAINER:        {C.SCARD_E_NO_KEY_CONTAINER, "SCARD_E_NO_KEY_CONTAINER", "requested key container does not exist on the smart card"},
	C.SCARD_E_SERVER_TOO_BUSY:         {C.SCARD_E_SERVER_TOO_BUSY, "SCARD_E_SERVER_TOO_BUSY", "smart card resource manager too busy to complete operation"},
	C.SCARD_W_UNSUPPORTED_CARD:        {C.SCARD_W_UNSUPPORTED_CARD, "SCARD_W_UNSUPPORTED_CARD", "reader cannot communicate with card due to ATR configuration conflicts"},
	C.SCARD_W_UNRESPONSIVE_CARD:       {C.SCARD_W_UNRESPONSIVE_CARD, "SCARD_W_UNRESPONSIVE_CARD", "smart card is not responding to a reset"},
	C.SCARD_W_UNPOWERED_CARD:          {C.SCARD_W_UNPOWERED_CARD, "SCARD_W_UNPOWERED_CARD", "power has been removed from the smart card"},
	C.SCARD_W_RESET_CARD:              {C.SCARD_W_RESET_CARD, "SCARD_W_RESET_CARD", "smart card has been reset"},
	C.SCARD_W_REMOVED_CARD:            {C.SCARD_W_REMOVED_CARD, "SCARD_W_REMOVED_CARD", "smart card has been removed"},
	C.SCARD_W_SECURITY_VIOLATION:      {C.SCARD_W_SECURITY_VIOLATION, "SCARD_W_SECURITY_VIOLATION", "access denied due to a security violation"},
	C.SCARD_W_WRONG_CHV:               {C.SCARD_W_WRONG_CHV, "SCARD_W_WRONG_CHV", "wrong PIN presented to the smart card"},
	C.SCARD_W_CHV_BLOCKED:             {C.SCARD_W_CHV_BLOCKED, "SCARD_W_CHV_BLOCKED", "maximum number of PIN entry attempts reached"},
	C.SCARD_W_EOF:                     {C.SCARD_W_EOF, "SCARD_W_EOF", "end of the smart card file reached"},
	C.SCARD_W_CANCELLED_BY_USER:       {C.SCARD_W_CANCELLED_BY_USER, "SCARD_W_CANCELLED_BY_USER", "user cancelled the Smart Card Selection Dialog"},
	C.SCARD_W_CARD_NOT_AUTHENTICATED:  {C.SCARD_W_CARD_NOT_AUTHENTICATED, "SCARD_W_CARD_NOT_AUTHENTICATED", "no PIN was presented to the smart card"},
}

type ScardScope dword

const (
	ScardScopeUser     ScardScope = C.SCARD_SCOPE_USER
	ScardScopeTerminal ScardScope = C.SCARD_SCOPE_TERMINAL
	ScardScopeGlobal   ScardScope = C.SCARD_SCOPE_GLOBAL
	ScardScopeSystem   ScardScope = C.SCARD_SCOPE_SYSTEM
)

// type ShareType int

// const (
// 	ScardShareShared    ShareType = C.SCARD_SHARE_SHARED
// 	ScardShareExclusive ShareType = C.SCARD_SHARE_EXCLUSIVE
// 	ScardShareDirect    ShareType = C.SCARD_SHARE_DIRECT
// )

// type DispositionType int

// const (
// 	ScardLeaveCard   DispositionType = C.SCARD_LEAVE_CARD
// 	ScardResetCard   DispositionType = C.SCARD_RESET_CARD
// 	ScardUnpowerCard DispositionType = C.SCARD_UNPOWER_CARD
// 	ScardEjectCard   DispositionType = C.SCARD_EJECT_CARD
// )

// type StateType int

// const (
// 	ScardStateUnaware     StateType = C.SCARD_STATE_UNAWARE
// 	ScardStateIgnore      StateType = C.SCARD_STATE_IGNORE
// 	ScardStateChanged     StateType = C.SCARD_STATE_CHANGED
// 	ScardStateUnknown     StateType = C.SCARD_STATE_UNKNOWN
// 	ScardStateUnavailable StateType = C.SCARD_STATE_UNAVAILABLE
// 	ScardStateEmpty       StateType = C.SCARD_STATE_EMPTY
// 	ScardStatePresent     StateType = C.SCARD_STATE_PRESENT
// 	ScardStateAtrMatch    StateType = C.SCARD_STATE_ATRMATCH
// 	ScardStateExclusive   StateType = C.SCARD_STATE_EXCLUSIVE
// 	ScardStateInUse       StateType = C.SCARD_STATE_INUSE
// 	ScardStateMute        StateType = C.SCARD_STATE_MUTE
// )

// type ProtocolType int

// const (
// 	ScardProtocolUnset     ProtocolType = C.SCARD_PROTOCOL_UNSET
// 	ScardProtocolT0        ProtocolType = C.SCARD_PROTOCOL_T0
// 	ScardProtocolT1        ProtocolType = C.SCARD_PROTOCOL_T1
// 	ScardProtocolRaw       ProtocolType = C.SCARD_PROTOCOL_RAW
// 	ScardProtocolT15       ProtocolType = C.SCARD_PROTOCOL_T15
// 	ScardProtocolAny       ProtocolType = C.SCARD_PROTOCOL_ANY
// 	ScardProtocolUndefined ProtocolType = C.SCARD_PROTOCOL_UNDEFINED
// )

// type ErrorCode int

// const Infinite int = C.INFINITE
