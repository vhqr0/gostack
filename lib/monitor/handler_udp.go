package monitor

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vhqr0/gostack/lib/l4/udp"
)

type UDPTableSer struct {
	Entries []*UDPSockSer
}

type UDPSockSer struct {
	Local string
	Peer  string
}

func init() {
	RegistorMonitorHandler("/udp", handleUDP)
}

func marshalUDPTable(table *udp.UDPTable) ([]byte, error) {
	table.Mutex.RLock()
	defer table.Mutex.RUnlock()

	tableSer := &UDPTableSer{}

	for _, sock := range table.Entries {
		var local, peer string
		if sock.Local != nil {
			local = sock.Local.String()
		}
		if sock.Peer != nil {
			peer = sock.Peer.String()
		}
		sockSer := &UDPSockSer{
			Local: local,
			Peer:  peer,
		}
		tableSer.Entries = append(tableSer.Entries, sockSer)
	}
	return json.Marshal(tableSer)
}

func handleUDP(m *Monitor, w http.ResponseWriter, r *http.Request) {
	if buf, err := marshalUDPTable(m.Stack.UDPStack.SockTable); err != nil {
		log.Printf("monitor/udp error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.Write(buf)
	}
}
