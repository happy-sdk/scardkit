//go:build ignore
// +build ignore

package pcsc

/*
#cgo pkg-config: libpcsclite
#include <PCSC/winscard.h>
#include <PCSC/wintypes.h>
*/
import "C"

type ScopeType int

const (
	ScardScopeUser     ScopeType = C.SCARD_SCOPE_USER
	ScardScopeTerminal ScopeType = C.SCARD_SCOPE_TERMINAL
	ScardScopeSystem   ScopeType = C.SCARD_SCOPE_SYSTEM
)

type ShareType int

const (
	ScardShareShared    ShareType = C.SCARD_SHARE_SHARED
	ScardShareExclusive ShareType = C.SCARD_SHARE_EXCLUSIVE
	ScardShareDirect    ShareType = C.SCARD_SHARE_DIRECT
)

type DispositionType int

const (
	ScardLeaveCard   DispositionType = C.SCARD_LEAVE_CARD
	ScardResetCard   DispositionType = C.SCARD_RESET_CARD
	ScardUnpowerCard DispositionType = C.SCARD_UNPOWER_CARD
	ScardEjectCard   DispositionType = C.SCARD_EJECT_CARD
)

type StateType int

const (
	ScardStateUnaware     StateType = C.SCARD_STATE_UNAWARE
	ScardStateIgnore      StateType = C.SCARD_STATE_IGNORE
	ScardStateChanged     StateType = C.SCARD_STATE_CHANGED
	ScardStateUnknown     StateType = C.SCARD_STATE_UNKNOWN
	ScardStateUnavailable StateType = C.SCARD_STATE_UNAVAILABLE
	ScardStateEmpty       StateType = C.SCARD_STATE_EMPTY
	ScardStatePresent     StateType = C.SCARD_STATE_PRESENT
	ScardStateAtrMatch    StateType = C.SCARD_STATE_ATRMATCH
	ScardStateExclusive   StateType = C.SCARD_STATE_EXCLUSIVE
	ScardStateInUse       StateType = C.SCARD_STATE_INUSE
	ScardStateMute        StateType = C.SCARD_STATE_MUTE
)

type ProtocolType int

const (
	ScardProtocolUnset     ProtocolType = C.SCARD_PROTOCOL_UNSET
	ScardProtocolT0        ProtocolType = C.SCARD_PROTOCOL_T0
	ScardProtocolT1        ProtocolType = C.SCARD_PROTOCOL_T1
	ScardProtocolRaw       ProtocolType = C.SCARD_PROTOCOL_RAW
	ScardProtocolT15       ProtocolType = C.SCARD_PROTOCOL_T15
	ScardProtocolAny       ProtocolType = C.SCARD_PROTOCOL_ANY
	ScardProtocolUndefined ProtocolType = C.SCARD_PROTOCOL_UNDEFINED
)

type ResponseCode int

const (
	ScardSSuccess                ResponseCode = C.SCARD_S_SUCCESS
	ScardFInternalError          ResponseCode = C.SCARD_F_INTERNAL_ERROR
	ScardECancelled              ResponseCode = C.SCARD_E_CANCELLED
	ScardEInvalidHandle          ResponseCode = C.SCARD_E_INVALID_HANDLE
	ScardEInvalidParameter       ResponseCode = C.SCARD_E_INVALID_PARAMETER
	ScardEInvalidTarget          ResponseCode = C.SCARD_E_INVALID_TARGET
	ScardENoMemory               ResponseCode = C.SCARD_E_NO_MEMORY
	ScardFWaitedTooLong          ResponseCode = C.SCARD_F_WAITED_TOO_LONG
	ScardEInsufficientBuffer     ResponseCode = C.SCARD_E_INSUFFICIENT_BUFFER
	ScardEUnknownReader          ResponseCode = C.SCARD_E_UNKNOWN_READER
	ScardETimeout                ResponseCode = C.SCARD_E_TIMEOUT
	ScardESharingViolation       ResponseCode = C.SCARD_E_SHARING_VIOLATION
	ScardENoSmartcard            ResponseCode = C.SCARD_E_NO_SMARTCARD
	ScardEUnknownCard            ResponseCode = C.SCARD_E_UNKNOWN_CARD
	ScardECantDispose            ResponseCode = C.SCARD_E_CANT_DISPOSE
	ScardEProtoMismatch          ResponseCode = C.SCARD_E_PROTO_MISMATCH
	ScardENotReady               ResponseCode = C.SCARD_E_NOT_READY
	ScardEInvalidValue           ResponseCode = C.SCARD_E_INVALID_VALUE
	ScardESystemCancelled        ResponseCode = C.SCARD_E_SYSTEM_CANCELLED
	ScardFCommError              ResponseCode = C.SCARD_F_COMM_ERROR
	ScardFUnknownError           ResponseCode = C.SCARD_F_UNKNOWN_ERROR
	ScardEInvalidAtr             ResponseCode = C.SCARD_E_INVALID_ATR
	ScardENotTransacted          ResponseCode = C.SCARD_E_NOT_TRANSACTED
	ScardEReaderUnavailable      ResponseCode = C.SCARD_E_READER_UNAVAILABLE
	ScardEPciTooSmall            ResponseCode = C.SCARD_E_PCI_TOO_SMALL
	ScardEReaderUnsupported      ResponseCode = C.SCARD_E_READER_UNSUPPORTED
	ScardEDuplicateReader        ResponseCode = C.SCARD_E_DUPLICATE_READER
	ScardECardUnsupported        ResponseCode = C.SCARD_E_CARD_UNSUPPORTED
	ScardENoService              ResponseCode = C.SCARD_E_NO_SERVICE
	ScardEServiceStopped         ResponseCode = C.SCARD_E_SERVICE_STOPPED
	ScardENoReadersAvailable     ResponseCode = C.SCARD_E_NO_READERS_AVAILABLE
	ScardEUnsupportedFeature     ResponseCode = C.SCARD_E_UNSUPPORTED_FEATURE
	ScardWUnsupportedCard        ResponseCode = C.SCARD_W_UNSUPPORTED_CARD
	ScardWUnresponsiveCard       ResponseCode = C.SCARD_W_UNRESPONSIVE_CARD
	ScardWUnpoweredCard          ResponseCode = C.SCARD_W_UNPOWERED_CARD
	ScardWResetCard              ResponseCode = C.SCARD_W_RESET_CARD
	ScardWRemovedCard            ResponseCode = C.SCARD_W_REMOVED_CARD
	ScardWSecurityViolation      ResponseCode = C.SCARD_W_SECURITY_VIOLATION
	ScardWWrongChv               ResponseCode = C.SCARD_W_WRONG_CHV
	ScardWChvBlocked             ResponseCode = C.SCARD_W_CHV_BLOCKED
	ScardWEof                    ResponseCode = C.SCARD_W_EOF
	ScardWCancelledByUser        ResponseCode = C.SCARD_W_CANCELLED_BY_USER
	ScardWCardNotAuthenticated   ResponseCode = C.SCARD_W_CARD_NOT_AUTHENTICATED
	ScardEUnexpected             ResponseCode = C.SCARD_E_UNEXPECTED
	ScardEIccInstallation        ResponseCode = C.SCARD_E_ICC_INSTALLATION
	ScardEIccCreateorder         ResponseCode = C.SCARD_E_ICC_CREATEORDER
	ScardEDirNotFound            ResponseCode = C.SCARD_E_DIR_NOT_FOUND
	ScardEFileNotFound           ResponseCode = C.SCARD_E_FILE_NOT_FOUND
	ScardENoDir                  ResponseCode = C.SCARD_E_NO_DIR
	ScardENoFile                 ResponseCode = C.SCARD_E_NO_FILE
	ScardENoAccess               ResponseCode = C.SCARD_E_NO_ACCESS
	ScardEWriteTooMany           ResponseCode = C.SCARD_E_WRITE_TOO_MANY
	ScardEBadSeek                ResponseCode = C.SCARD_E_BAD_SEEK
	ScardEInvalidChv             ResponseCode = C.SCARD_E_INVALID_CHV
	ScardEUnknownResMng          ResponseCode = C.SCARD_E_UNKNOWN_RES_MNG
	ScardENoSuchCertificate      ResponseCode = C.SCARD_E_NO_SUCH_CERTIFICATE
	ScardECertificateUnavailable ResponseCode = C.SCARD_E_CERTIFICATE_UNAVAILABLE
	ScardECommDataLost           ResponseCode = C.SCARD_E_COMM_DATA_LOST
	ScardENoKeyContainer         ResponseCode = C.SCARD_E_NO_KEY_CONTAINER
	ScardEServerTooBusy          ResponseCode = C.SCARD_E_SERVER_TOO_BUSY
	ScardPShutdown               ResponseCode = C.SCARD_P_SHUTDOWN
)

const Infinite int = C.INFINITE
