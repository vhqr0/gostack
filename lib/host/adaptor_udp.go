package host

import (
	"net"
)

func init() {
	RegisterIfaceAdaptor("udp", NewUDPIfaceAdaptor)
}

func NewUDPIfaceAdaptor(args map[string]string) (IfaceAdaptor, error) {
	localStr, err := IfaceAdaptorRequireArg("Local", args)
	if err != nil {
		return nil, err
	}
	peerStr, err := IfaceAdaptorRequireArg("Peer", args)
	if err != nil {
		return nil, err
	}

	local, err := net.ResolveUDPAddr("udp", localStr)
	if err != nil {
		return nil, err
	}
	peer, err := net.ResolveUDPAddr("udp", peerStr)
	if err != nil {
		return nil, err
	}

	return net.DialUDP("udp", local, peer)
}
