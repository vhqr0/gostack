package host

import (
	"net"
)

func init() {
	RegisterAdaptor("udp", NewUDPAdaptor)
}

func NewUDPAdaptor(args map[string]string) (Adaptor, error) {
	localStr, err := AdaptorRequireArg("Local", args)
	if err != nil {
		return nil, err
	}
	peerStr, err := AdaptorRequireArg("Peer", args)
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
