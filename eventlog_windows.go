package netlog

import (
	"fmt"
	"log"
	"os"

	"github.com/golang/sys/windows/svc/eventlog"
)

const (
	facilitySystem Facility = iota
	facilityApplication
	facilityService
	facilitySecurity
)

type windowsLogger struct {
	w *eventlog.Log
}

func (w windowsLogger) Info(format string, v ...interface{}) {
	w.w.Info(1001, fmt.Sprintf(format, v...))
}

func (w windowsLogger) Warning(format string, v ...interface{}) {
	w.w.Warning(2001, fmt.Sprintf(format, v...))
}

func (w windowsLogger) Err(format string, v ...interface{}) {
	w.w.Error(3001, fmt.Sprintf(format, v...))
}

func (w windowsLogger) Crit(format string, v ...interface{}) {
	w.w.Error(4001, fmt.Sprintf(format, v...))
	os.Exit(2)
}

func NewLogger(facility Facility, tag string) Logger {
	w, err := eventlog.Open(tag)
	if err != nil {
		log.Fatal(err)
	}
	return windowsLogger{w: w}
}
