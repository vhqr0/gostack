package sock

import (
	"fmt"
)

type (
	OPUnsupported  struct{ OP string }
	OPOnClosedSock struct{ OP string }
	OPOnFreeSock   struct{ OP string }
	OPOnBusySock   struct{ OP string }
)

func (err *OPUnsupported) Error() string  { return "op unsupported: " + err.OP }
func (err *OPOnClosedSock) Error() string { return "op on closed sock: " + err.OP }
func (err *OPOnFreeSock) Error() string   { return "op on free sock: " + err.OP }
func (err *OPOnBusySock) Error() string   { return "op on busy sock: " + err.OP }

type (
	InvalidSockFamilyOrTyp struct{ Family, Typ uint32 }
	InvalidSockAddr        struct{ Addr *Addr }
	BusySockAddr           struct{ Addr *Addr }
)

func (err *InvalidSockFamilyOrTyp) Error() string {
	return fmt.Sprintf("invalid sock family or typ: %d %d", err.Family, err.Typ)
}

func (err *InvalidSockAddr) Error() string { return "invalid sock addr: " + err.Addr.String() }
func (err *BusySockAddr) Error() string    { return "busy sock addr: " + err.Addr.String() }
