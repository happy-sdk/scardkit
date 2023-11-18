//go:build ignore
// +build ignore

package pcsc

/*
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
