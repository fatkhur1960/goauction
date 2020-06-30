package monitor

import (
	"log"
	"reflect"
	"time"
)

// Monitor ... Abstraksi untuk sistem monitor
type Monitor interface {
	Start()
	Stop()
}

// StartMonitors Run all monitors
func StartMonitors() {
	monitors := []Monitor{NewProductMonitor()}

	time.Sleep(5 * time.Second)
	for _, monitor := range monitors {
		log.Printf("Monitor] Starting `%s`...\n", reflect.TypeOf(monitor).String())
		go monitor.Start()
	}
}
