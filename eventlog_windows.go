package netlog

import (
	"errors"
	"fmt"
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
	d bool
}

func (w windowsLogger) Debug(format string, v ...interface{}) {
	if w.d {
		w.w.Info(1001, fmt.Sprintf(format, v...))
	}
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

func NewLogger(f Facility, tag string, debug bool, addr ...string) (logger Logger, err error) {
	var w *eventlog.Log
	if len(addr) > 0 {
		err = errors.New("not implemented")
	} else {
		w, err = eventlog.Open(tag)
	}
	return windowsLogger{w: w, d:debug}, err
}
