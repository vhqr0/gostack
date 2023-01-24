package l3

import (
	"net"
	"strconv"
)

type (
	DstUnreach  struct{ Dst net.IP }
	HostUnreach struct{ Host net.IP }
	PktTooBig   struct{ Len int }
)

func (err *DstUnreach) Error() string {
	return "dst unreach: " + err.Dst.String()
}

func (err *HostUnreach) Error() string {
	return "host unreach: " + err.Host.String()
}

func (err *PktTooBig) Error() string {
	return "pkt too big: " + strconv.FormatInt(int64(err.Len), 10)
}

type (
	InvalidIPVer struct{ Ver int }
	InvalidIPLen struct{ Len int }
)

func (err *InvalidIPVer) Error() string {
	return "invalid ip ver: " + strconv.FormatInt(int64(err.Ver), 10)
}

func (err *InvalidIPLen) Error() string {
	return "invalid ip len: " + strconv.FormatInt(int64(err.Len), 10)
}
