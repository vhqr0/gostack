package sock

import (
	"io"
)

const (
	AFIP4      uint32 = 2
	AFIP6      uint32 = 10
	SockStream uint32 = 1
	SockDgram  uint32 = 2

	ShutR  = 0
	ShutW  = 1
	ShutRW = 2
)

type Sock interface {
	io.ReadWriteCloser

	ReadFrom([]byte) (int, *Addr, error)
	WriteTo([]byte, *Addr) (int, error)
	Shutdown(int) error

	Bind(*Addr) (*Addr, error)
	Connect(*Addr) error
	Listen() error
	Accept() (Sock, error)

	GetSockName() *Addr
	GetPeerName() *Addr
}
