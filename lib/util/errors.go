package util

import (
	"strconv"
)

type UnsupportedFeature struct{ Feature string }

func (err *UnsupportedFeature) Error() string {
	return "unsupported feature: " + err.Feature
}

type (
	InvalidIfaceName  struct{ Name string }
	IfaceNameOverflow struct{ Name string }
)

func (err *InvalidIfaceName) Error() string {
	return "invalid iface name: " + err.Name
}

func (err *IfaceNameOverflow) Error() string {
	return "iface name overflow: " + err.Name
}

type (
	InvalidIPVer struct{ Ver int }
	InvalidIPLen struct{ Len int }
	InvalidIPStr struct{ Str string }
)

func (err *InvalidIPVer) Error() string {
	return "invalid IP ver: " + strconv.FormatInt(int64(err.Ver), 10)
}

func (err *InvalidIPLen) Error() string {
	return "invalid IP len: " + strconv.FormatInt(int64(err.Len), 10)
}

func (err *InvalidIPStr) Error() string {
	return "invalid IP str: " + err.Str
}
