package sock

import (
	"fmt"
)

type (
	InvalidSockArgs struct{ Family, Typ uint32 }
	InvalidSockAddr struct{ Addr *Addr }
	BusySockAddr    struct{ Addr *Addr }
	InvalidSockOP   struct{ OP string }
)

func (err *InvalidSockArgs) Error() string {
	return fmt.Sprintf("invalid sock args: family=%d typ=%d", err.Family, err.Typ)
}

func (err *InvalidSockAddr) Error() string { return "invalid sock addr: " + err.Addr.String() }
func (err *BusySockAddr) Error() string    { return "busy sock addr: " + err.Addr.String() }
func (err *InvalidSockOP) Error() string   { return "invalid sock op: " + err.OP }

type (
	OPOnClosedSock struct{ OP string }
	OPOnFreeSock   struct{ OP string }
	OPOnBusySock   struct{ OP string }
)

func (err *OPOnClosedSock) Error() string { return "op on closed sock: " + err.OP }
func (err *OPOnFreeSock) Error() string   { return "op on free sock: " + err.OP }
func (err *OPOnBusySock) Error() string   { return "op on busy sock: " + err.OP }
