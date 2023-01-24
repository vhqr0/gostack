package monitor

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/vhqr0/gostack/lib/l3"
	"github.com/vhqr0/gostack/lib/stack"
)

type NeighTableSer struct {
	TS        int64
	ValidSecs int64
	Entries   []*NeighEntrySer
}

type NeighEntrySer struct {
	Iface string
	IP    string
	MAC   string
	TS    int64
}

func init() {
	RegistorMonitorHandler("/neigh", handleNeigh)
}

func marshalNeighTable(vstack *stack.Stack) ([]byte, error) {
	table := vstack.IPStack.NeighTable

	table.Mutex.RLock()
	defer table.Mutex.RUnlock()

	tableSer := &NeighTableSer{
		TS:        time.Now().Unix(),
		ValidSecs: l3.NeighValidSecs,
	}

	for k, v := range table.Entries {
		kbytes := []byte(k)
		ifidx := binary.BigEndian.Uint16(kbytes[:2])
		ip := net.IP(kbytes[2:])
		mac := v.MAC
		ts := v.TS
		if len(ip) != 4 && len(ip) != 16 {
			return nil, &l3.InvalidIPLen{Len: len(ip)}
		}
		entrySer := &NeighEntrySer{
			Iface: vstack.Host.Ifaces[ifidx].Name,
			IP:    ip.String(),
			MAC:   mac.String(),
			TS:    ts,
		}
		tableSer.Entries = append(tableSer.Entries, entrySer)
	}

	return json.Marshal(tableSer)
}

func handleNeigh(m *Monitor, w http.ResponseWriter, r *http.Request) {
	if buf, err := marshalNeighTable(m.Stack); err != nil {
		log.Printf("monitor/neigh error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.Write(buf)
	}
}
