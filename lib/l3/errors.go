package l3

import (
	"fmt"
	"net"
	"strconv"
)

type (
	DstUnreach  struct{ Dst net.IP }
	HostUnreach struct{ Host net.IP }
	PktTooBig   struct{ Len int }
)

func (err *DstUnreach) Error() string {
	return fmt.Sprintf("dst unreach: %v", err.Dst)
}

func (err *HostUnreach) Error() string {
	return fmt.Sprintf("host unreach: %v", err.Host)
}

func (err *PktTooBig) Error() string {
	return "pkt too big: " + strconv.FormatInt(int64(err.Len), 10)
}
