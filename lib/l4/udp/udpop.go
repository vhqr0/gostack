package udp

import (
	"net"

	"github.com/vhqr0/gostack/lib/sock"
	"github.com/vhqr0/gostack/lib/util"
)

const (
	UDPBindMaxRetry = 3
)

func (udpSock *UDPSock) ReadFrom(p []byte) (n int, addr *sock.Addr, err error) {
	if udpSock.Status == UDPClosed {
		err = &sock.OPOnClosedSock{OP: "ReadFrom"}
		return
	}

	if udpSock.Local == nil {
		err = &sock.OPOnFreeSock{OP: "ReadFrom"}
		return
	}

	if pkt, ok := <-udpSock.recvCh; !ok { // Notice: block
		err = &sock.OPOnClosedSock{OP: "ReadFrom"}
	} else {
		addr = pkt.Peer
		n = copy(p, pkt.Payload)
	}
	return
}

func (udpSock *UDPSock) Read(p []byte) (n int, err error) {
	n, _, err = udpSock.ReadFrom(p)
	return
}

func (udpSock *UDPSock) WriteTo(p []byte, addr *sock.Addr) (int, error) {
	if addr == nil {
		return 0, &sock.InvalidSockAddr{Addr: addr}
	}

	if err := udpSock.ValidateAddr(addr); err != nil {
		return 0, err
	}

	local, peer := udpSock.Local, udpSock.Peer // inhibit change

	if udpSock.Status == UDPClosed {
		return 0, &sock.OPOnClosedSock{OP: "WriteTo"}
	}

	if peer != nil && addr != peer { // Notice: same ptr, for Write
		return 0, &sock.OPOnBusySock{OP: "WriteTo"}
	}

	if local == nil {
		var err error
		if local, err = udpSock.Bind(nil); err != nil {
			return 0, err
		}
	}

	pkt := &UDPPkt{
		Local:   local,
		Peer:    addr,
		Payload: p,
	}

	return len(p), udpSock.Stack.udpSend(pkt) // Notice: block
}

func (udpSock *UDPSock) Write(p []byte) (int, error) {
	return udpSock.WriteTo(p, udpSock.Peer)
}

func (udpSock *UDPSock) Close() error {
	udpSock.Mutex.Lock()
	defer udpSock.Mutex.Unlock()

	if udpSock.Status == UDPClosed {
		return &sock.OPOnClosedSock{OP: "Close"}
	}

	udpSock.Status = UDPClosed
	close(udpSock.recvCh)

	if udpSock.Local != nil {
		udpSock.Stack.UDPTable.Del(udpSock.Local)
	}

	udpSock.Local = nil
	udpSock.Peer = nil

	return nil
}

func (udpSock *UDPSock) Bind(addr *sock.Addr) (*sock.Addr, error) {
	udpSock.Mutex.Lock()
	defer udpSock.Mutex.Unlock()

	if udpSock.Status == UDPClosed {
		return nil, &sock.OPOnClosedSock{OP: "Bind"}
	}

	if udpSock.Local != nil {
		return nil, &sock.OPOnBusySock{OP: "Bind"}
	}

	var addrClone *sock.Addr // don't modify origin addr
	if addr != nil {
		addrClone = &sock.Addr{IP: addr.IP, Port: addr.Port}
	}

	switch udpSock.Family {
	case sock.AFIP4:
		if addrClone == nil {
			addrClone = &sock.Addr{IP: net.IPv4zero.To4()}
		} else if len(addrClone.IP) != 4 {
			return nil, &sock.InvalidSockAddr{Addr: addrClone}
		}
		if addrClone.IP.Equal(net.IPv4zero) {
			addrClone.IP = udpSock.Stack.IPStack.EthStack.Host.Ifaces[0].IP4 // Notice: assume ifaces != nil
		}
	case sock.AFIP6:
		if addrClone == nil {
			addrClone = &sock.Addr{IP: net.IPv6zero}
		} else if len(addrClone.IP) != 16 {
			return nil, &sock.InvalidSockAddr{Addr: addrClone}
		}
		if addrClone.IP.Equal(net.IPv6zero) {
			addrClone.IP = udpSock.Stack.IPStack.EthStack.Host.Ifaces[0].IP6 // Notice: assume ifaces != nil
		}
	default:
		return nil, &sock.InvalidSockArgs{Family: udpSock.Family, Typ: udpSock.Typ}
	}

	if !udpSock.Stack.IPStack.IsInStack(addrClone.IP) {
		return nil, &sock.InvalidSockAddr{Addr: addrClone}
	}

	if addrClone.Port == 0 {
		var err error
		for i := 0; i < UDPBindMaxRetry; i++ {
			addrClone.Port = util.RandUint16()
			err = udpSock.Stack.UDPTable.Add(addrClone, udpSock)
			if err == nil {
				return addrClone, nil
			}
		}
		return nil, err
	}

	if err := udpSock.Stack.UDPTable.Add(addrClone, udpSock); err != nil {
		return nil, err
	}

	return addrClone, nil
}

func (udpSock *UDPSock) Connect(addr *sock.Addr) error {
	udpSock.Mutex.Lock()
	defer udpSock.Mutex.Unlock()

	if udpSock.Status == UDPClosed {
		return &sock.OPOnClosedSock{OP: "Connect"}
	}

	if udpSock.Peer != nil {
		return &sock.OPOnBusySock{OP: "Connect"}
	}

	if err := udpSock.ValidateAddr(addr); err != nil {
		return err
	}

	udpSock.Peer = addr

	return nil
}

func (udpSock *UDPSock) Shutdown(op int) error      { return &sock.InvalidSockOP{OP: "Shutdown"} }
func (udpSock *UDPSock) Listen() error              { return &sock.InvalidSockOP{OP: "Listen"} }
func (udpSock *UDPSock) Accept() (sock.Sock, error) { return nil, &sock.InvalidSockOP{OP: "Accept"} }

func (udpSock *UDPSock) GetSockName() *sock.Addr { return udpSock.Local }
func (udpSock *UDPSock) GetPeerName() *sock.Addr { return udpSock.Peer }
