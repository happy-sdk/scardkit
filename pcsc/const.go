// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo -importpath github.com/happy-sdk/nfcsdk/pcsc -objdir /tmp/pcnfcsdk-1493475376 -godefs -- -I /usr/include/PCSC/ pcsc/const_c.go

package pcsc

const (
	MaxAtrSize = 0x21
	Infinite   = 0xffffffff
)

type (
	long         int64
	scardContext int64
	dword        uint64
	hContext     int64
	lpcvoid      *byte
	lpcstr       *int8
	lpstr        *int8
	lpdword      *uint64
	returnValue  long
)

type ScardScope dword

const (
	ScardScopeUser     ScardScope = 0x0
	ScardScopeTerminal ScardScope = 0x1
	ScardScopeGlobal   ScardScope = 0x3
	ScardScopeSystem   ScardScope = 0x2
)

type ScardState dword

const (
	ScardStateUnaware     ScardState = 0x0
	ScardStateIgnore      ScardState = 0x1
	ScardStateChanged     ScardState = 0x2
	ScardStateUnknown     ScardState = 0x4
	ScardStateUnavailable ScardState = 0x8
	ScardStateEmpty       ScardState = 0x10
	ScardStatePresent     ScardState = 0x20
	ScardStateAtrMatch    ScardState = 0x40
	ScardStateExclusive   ScardState = 0x80
	ScardStateInUse       ScardState = 0x100
	ScardStateMute        ScardState = 0x200
	ScardStateUnpowered   ScardState = 0x400
)

// Code generated by go generateor; DO NOT EDIT.// by calling go generate . in package root
var (
	ErrScardSSuccess                = Error{0x0, "SCARD_S_SUCCESS", "no error encountered", nil}
	ErrScardFInternalError          = Error{0x80100001, "SCARD_F_INTERNAL_ERROR", "internal consistency check failed", nil}
	ErrScardECancelled              = Error{0x80100002, "SCARD_E_CANCELLED", "action cancelled by SCardCancel request", nil}
	ErrScardEInvalidHandle          = Error{0x80100003, "SCARD_E_INVALID_HANDLE", "supplied handle was invalid", nil}
	ErrScardEInvalidParameter       = Error{0x80100004, "SCARD_E_INVALID_PARAMETER", "parameters could not be properly interpreted", nil}
	ErrScardEInvalidTarget          = Error{0x80100005, "SCARD_E_INVALID_TARGET", "registry startup information is missing or invalid", nil}
	ErrScardENoMemory               = Error{0x80100006, "SCARD_E_NO_MEMORY", "not enough memory available to complete command", nil}
	ErrScardFWaitedTooLong          = Error{0x80100007, "SCARD_F_WAITED_TOO_LONG", "internal consistency timer expired", nil}
	ErrScardEInsufficientBuffer     = Error{0x80100008, "SCARD_E_INSUFFICIENT_BUFFER", "data buffer too small for returned data", nil}
	ErrScardEUnknownReader          = Error{0x80100009, "SCARD_E_UNKNOWN_READER", "specified reader name not recognized", nil}
	ErrScardETimeout                = Error{0x8010000a, "SCARD_E_TIMEOUT", "user-specified timeout expired", nil}
	ErrScardESharingViolation       = Error{0x8010000b, "SCARD_E_SHARING_VIOLATION", "smart card cannot be accessed due to other connections", nil}
	ErrScardENoSmartcard            = Error{0x8010000c, "SCARD_E_NO_SMARTCARD", "operation requires a smart card, but none is in the device", nil}
	ErrScardEUnknownCard            = Error{0x8010000d, "SCARD_E_UNKNOWN_CARD", "specified smart card name not recognized", nil}
	ErrScardECantDispose            = Error{0x8010000e, "SCARD_E_CANT_DISPOSE", "system could not dispose of the media as requested", nil}
	ErrScardEProtoMismatch          = Error{0x8010000f, "SCARD_E_PROTO_MISMATCH", "requested protocols incompatible with card's protocols", nil}
	ErrScardENotReady               = Error{0x80100010, "SCARD_E_NOT_READY", "reader or smart card not ready to accept commands", nil}
	ErrScardEInvalidValue           = Error{0x80100011, "SCARD_E_INVALID_VALUE", "one or more supplied parameter values could not be interpreted", nil}
	ErrScardESystemCancelled        = Error{0x80100012, "SCARD_E_SYSTEM_CANCELLED", "action cancelled by the system", nil}
	ErrScardFCommError              = Error{0x80100013, "SCARD_F_COMM_ERROR", "internal communications error detected", nil}
	ErrScardFUnknownError           = Error{0x80100014, "SCARD_F_UNKNOWN_ERROR", "internal error detected, but source unknown", nil}
	ErrScardEInvalidAtr             = Error{0x80100015, "SCARD_E_INVALID_ATR", "ATR from registry is not a valid ATR string", nil}
	ErrScardENotTransacted          = Error{0x80100016, "SCARD_E_NOT_TRANSACTED", "attempt made to end a non-existent transaction", nil}
	ErrScardEReaderUnavailable      = Error{0x80100017, "SCARD_E_READER_UNAVAILABLE", "specified reader not currently available", nil}
	ErrScardPShutdown               = Error{0x80100018, "SCARD_P_SHUTDOWN", "operation aborted to allow server application to exit", nil}
	ErrScardEPciTooSmall            = Error{0x80100019, "SCARD_E_PCI_TOO_SMALL", "PCI receive buffer was too small", nil}
	ErrScardEReaderUnsupported      = Error{0x8010001a, "SCARD_E_READER_UNSUPPORTED", "reader driver does not meet minimal requirements", nil}
	ErrScardEDuplicateReader        = Error{0x8010001b, "SCARD_E_DUPLICATE_READER", "reader driver did not produce a unique reader name", nil}
	ErrScardECardUnsupported        = Error{0x8010001c, "SCARD_E_CARD_UNSUPPORTED", "smart card does not meet minimal requirements", nil}
	ErrScardENoService              = Error{0x8010001d, "SCARD_E_NO_SERVICE", "smart card resource manager is not running", nil}
	ErrScardEServiceStopped         = Error{0x8010001e, "SCARD_E_SERVICE_STOPPED", "smart card resource manager has shut down", nil}
	ErrScardEUnexpected             = Error{0x8010001f, "SCARD_E_UNEXPECTED", "unexpected card error occurred", nil}
	ErrScardEIccInstallation        = Error{0x80100020, "SCARD_E_ICC_INSTALLATION", "no primary provider can be found for the smart card", nil}
	ErrScardEIccCreateorder         = Error{0x80100021, "SCARD_E_ICC_CREATEORDER", "requested order of object creation not supported", nil}
	ErrScardEDirNotFound            = Error{0x80100023, "SCARD_E_DIR_NOT_FOUND", "identified directory does not exist on the smart card", nil}
	ErrScardEFileNotFound           = Error{0x80100024, "SCARD_E_FILE_NOT_FOUND", "identified file does not exist on the smart card", nil}
	ErrScardENoDir                  = Error{0x80100025, "SCARD_E_NO_DIR", "supplied path does not represent a smart card directory", nil}
	ErrScardENoFile                 = Error{0x80100026, "SCARD_E_NO_FILE", "supplied path does not represent a smart card file", nil}
	ErrScardENoAccess               = Error{0x80100027, "SCARD_E_NO_ACCESS", "access denied to this file", nil}
	ErrScardEWriteTooMany           = Error{0x80100028, "SCARD_E_WRITE_TOO_MANY", "smart card does not have enough memory to store information", nil}
	ErrScardEBadSeek                = Error{0x80100029, "SCARD_E_BAD_SEEK", "error trying to set smart card file object pointer", nil}
	ErrScardEInvalidChv             = Error{0x8010002a, "SCARD_E_INVALID_CHV", "supplied PIN is incorrect", nil}
	ErrScardEUnknownResMng          = Error{0x8010002b, "SCARD_E_UNKNOWN_RES_MNG", "unrecognized error code from a layered component", nil}
	ErrScardENoSuchCertificate      = Error{0x8010002c, "SCARD_E_NO_SUCH_CERTIFICATE", "requested certificate does not exist", nil}
	ErrScardECertificateUnavailable = Error{0x8010002d, "SCARD_E_CERTIFICATE_UNAVAILABLE", "requested certificate could not be obtained", nil}
	ErrScardENoReadersAvailable     = Error{0x8010002e, "SCARD_E_NO_READERS_AVAILABLE", "cannot find a smart card reader", nil}
	ErrScardECommDataLost           = Error{0x8010002f, "SCARD_E_COMM_DATA_LOST", "communications error with smart card detected", nil}
	ErrScardENoKeyContainer         = Error{0x80100030, "SCARD_E_NO_KEY_CONTAINER", "requested key container does not exist on the smart card", nil}
	ErrScardEServerTooBusy          = Error{0x80100031, "SCARD_E_SERVER_TOO_BUSY", "smart card resource manager too busy to complete operation", nil}
	ErrScardWUnsupportedCard        = Error{0x80100065, "SCARD_W_UNSUPPORTED_CARD", "reader cannot communicate with card due to ATR configuration conflicts", nil}
	ErrScardWUnresponsiveCard       = Error{0x80100066, "SCARD_W_UNRESPONSIVE_CARD", "smart card is not responding to a reset", nil}
	ErrScardWUnpoweredCard          = Error{0x80100067, "SCARD_W_UNPOWERED_CARD", "power has been removed from the smart card", nil}
	ErrScardWResetCard              = Error{0x80100068, "SCARD_W_RESET_CARD", "smart card has been reset", nil}
	ErrScardWRemovedCard            = Error{0x80100069, "SCARD_W_REMOVED_CARD", "smart card has been removed", nil}
	ErrScardWSecurityViolation      = Error{0x8010006a, "SCARD_W_SECURITY_VIOLATION", "access denied due to a security violation", nil}
	ErrScardWWrongChv               = Error{0x8010006b, "SCARD_W_WRONG_CHV", "wrong PIN presented to the smart card", nil}
	ErrScardWChvBlocked             = Error{0x8010006c, "SCARD_W_CHV_BLOCKED", "maximum number of PIN entry attempts reached", nil}
	ErrScardWEof                    = Error{0x8010006d, "SCARD_W_EOF", "end of the smart card file reached", nil}
	ErrScardWCancelledByUser        = Error{0x8010006e, "SCARD_W_CANCELLED_BY_USER", "user cancelled the Smart Card Selection Dialog", nil}
	ErrScardWCardNotAuthenticated   = Error{0x8010006f, "SCARD_W_CARD_NOT_AUTHENTICATED", "no PIN was presented to the smart card", nil}
)

var cerrors = map[returnValue]Error{
	0x0:        ErrScardSSuccess,
	0x80100001: ErrScardFInternalError,
	0x80100002: ErrScardECancelled,
	0x80100003: ErrScardEInvalidHandle,
	0x80100004: ErrScardEInvalidParameter,
	0x80100005: ErrScardEInvalidTarget,
	0x80100006: ErrScardENoMemory,
	0x80100007: ErrScardFWaitedTooLong,
	0x80100008: ErrScardEInsufficientBuffer,
	0x80100009: ErrScardEUnknownReader,
	0x8010000a: ErrScardETimeout,
	0x8010000b: ErrScardESharingViolation,
	0x8010000c: ErrScardENoSmartcard,
	0x8010000d: ErrScardEUnknownCard,
	0x8010000e: ErrScardECantDispose,
	0x8010000f: ErrScardEProtoMismatch,
	0x80100010: ErrScardENotReady,
	0x80100011: ErrScardEInvalidValue,
	0x80100012: ErrScardESystemCancelled,
	0x80100013: ErrScardFCommError,
	0x80100014: ErrScardFUnknownError,
	0x80100015: ErrScardEInvalidAtr,
	0x80100016: ErrScardENotTransacted,
	0x80100017: ErrScardEReaderUnavailable,
	0x80100018: ErrScardPShutdown,
	0x80100019: ErrScardEPciTooSmall,
	0x8010001a: ErrScardEReaderUnsupported,
	0x8010001b: ErrScardEDuplicateReader,
	0x8010001c: ErrScardECardUnsupported,
	0x8010001d: ErrScardENoService,
	0x8010001e: ErrScardEServiceStopped,
	0x8010001f: ErrScardEUnexpected,
	0x80100020: ErrScardEIccInstallation,
	0x80100021: ErrScardEIccCreateorder,
	0x80100023: ErrScardEDirNotFound,
	0x80100024: ErrScardEFileNotFound,
	0x80100025: ErrScardENoDir,
	0x80100026: ErrScardENoFile,
	0x80100027: ErrScardENoAccess,
	0x80100028: ErrScardEWriteTooMany,
	0x80100029: ErrScardEBadSeek,
	0x8010002a: ErrScardEInvalidChv,
	0x8010002b: ErrScardEUnknownResMng,
	0x8010002c: ErrScardENoSuchCertificate,
	0x8010002d: ErrScardECertificateUnavailable,
	0x8010002e: ErrScardENoReadersAvailable,
	0x8010002f: ErrScardECommDataLost,
	0x80100030: ErrScardENoKeyContainer,
	0x80100031: ErrScardEServerTooBusy,
	0x80100065: ErrScardWUnsupportedCard,
	0x80100066: ErrScardWUnresponsiveCard,
	0x80100067: ErrScardWUnpoweredCard,
	0x80100068: ErrScardWResetCard,
	0x80100069: ErrScardWRemovedCard,
	0x8010006a: ErrScardWSecurityViolation,
	0x8010006b: ErrScardWWrongChv,
	0x8010006c: ErrScardWChvBlocked,
	0x8010006d: ErrScardWEof,
	0x8010006e: ErrScardWCancelledByUser,
	0x8010006f: ErrScardWCardNotAuthenticated,
}
