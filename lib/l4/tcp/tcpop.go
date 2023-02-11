package tcp

import (
	"github.com/vhqr0/gostack/lib/sock"
)

func (tcpSock *TCPSock) Read(p []byte) (n int, err error) {
	return 0, &sock.InvalidSockOP{}
}

func (tcpSock *TCPSock) ReadFrom(p []byte) (n int, addr *sock.Addr, err error) {
	n, err = tcpSock.Read(p)
	return n, tcpSock.Peer, err
}

func (tcpSock *TCPSock) Write(p []byte) (n int, err error) {
	return 0, &sock.InvalidSockOP{}
}

func (tcpSock *TCPSock) WriteTo(p []byte, addr *sock.Addr) (int, error) {
	return 0, &sock.InvalidSockOP{}
}

func (tcpSock *TCPSock) Close() error      { return &sock.InvalidSockOP{} }
func (tcpSock *TCPSock) CloseRead() error  { return &sock.InvalidSockOP{} }
func (tcpSock *TCPSock) CloseWrite() error { return &sock.InvalidSockOP{} }

func (tcpSock *TCPSock) Bind(adddr *sock.Addr) (*sock.Addr, error) {
	return nil, &sock.InvalidSockOP{}
}

func (tcpSock *TCPSock) Connect(addr *sock.Addr) error {
	return &sock.InvalidSockOP{}
}

func (tcpSock *TCPSock) Listen() error {
	return  &sock.InvalidSockOP{}
}

func (tcpSock *TCPSock) Accept() (sock.Sock, error) {
	return nil, &sock.InvalidSockOP{}
}

func (tcpSock *TCPSock) GetSockName() *sock.Addr {return tcpSock.Local}
func (tcpSock *TCPSock) GetPeerName() *sock.Addr { return tcpSock.Peer }
