package monitor

import (
	"log"
	"net/http"

	"github.com/vhqr0/gostack/lib/conf"
)

func init() {
	RegistorMonitorHandler("/conf", handleConf)
}

func handleConf(m *Monitor, w http.ResponseWriter, r *http.Request) {
	stackConf := conf.ConfFromStack(m.Stack)
	if buf, err := stackConf.Marshal(); err != nil {
		log.Printf("monitor[conf] error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(buf)
	}
}
