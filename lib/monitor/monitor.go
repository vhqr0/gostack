package monitor

import (
	"log"
	"net/http"

	"github.com/vhqr0/gostack/lib/stack"
)

var (
	monitorHandlerMap map[string]MonitorHandler = make(map[string]MonitorHandler)
)

type (
	Monitor        struct{ Stack *stack.Stack }
	MonitorHandler func(*Monitor, http.ResponseWriter, *http.Request)
)

func (m *Monitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if handler, ok := monitorHandlerMap[path]; ok {
		log.Printf("monitor recv " + path)
		handler(m, w, r)
	} else {
		log.Printf("monitor drop " + path)
		w.WriteHeader(http.StatusNotFound)
	}
}

func RegistorMonitorHandler(path string, handler MonitorHandler) {
	monitorHandlerMap[path] = handler
}
