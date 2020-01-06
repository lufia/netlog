// +build darwin linux

package netlog

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
)

const (
	facilitySystem      = Facility(syslog.LOG_KERN)
	facilityApplication = Facility(syslog.LOG_USER)
	facilityService     = Facility(syslog.LOG_DAEMON)
	facilitySecurity    = Facility(syslog.LOG_AUTHPRIV)
)

type unixLogger struct {
	w *syslog.Writer
	d bool
}

func (w unixLogger) Debug(format string, v ...interface{}) {
	if w.d {
		w.w.Debug(fmt.Sprintf(format, v...))
	}
}

func (w unixLogger) Info(format string, v ...interface{}) {
	w.w.Info(fmt.Sprintf(format, v...))
}

func (w unixLogger) Warning(format string, v ...interface{}) {
	w.w.Warning(fmt.Sprintf(format, v...))
}

func (w unixLogger) Err(format string, v ...interface{}) {
	w.w.Err(fmt.Sprintf(format, v...))
}

func (w unixLogger) Crit(format string, v ...interface{}) {
	w.w.Crit(fmt.Sprintf(format, v...))
	os.Exit(2)
}

func NewLogger(f Facility, tag string, debug bool) Logger {
	w, err := syslog.New(syslog.Priority(f), tag)
	if err != nil {
		log.Fatal(err)
	}
	return unixLogger{w: w, d: debug}
}
