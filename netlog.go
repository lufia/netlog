package netlog

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"
)

type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Err(format string, v ...interface{})
	Crit(format string, v ...interface{})

	// SetDebug can enable/disable debug mode.
	SetDebug(bool)
}

const (
	logHeaderDebug   = "debug: "
	logHeaderInfo    = "info:  "
	logHeaderWarning = "warn:  "
	logHeaderErr     = "error: "
	logHeaderCrit    = "crit:  "
)

type Facility int

const (
	LogSystem      = facilitySystem
	LogApplication = facilityApplication
	LogService     = facilityService
	LogSecurity    = facilitySecurity
)

func parseFacility(s string) (Facility, error) {
	switch s {
	case "sys", "system":
		return LogSystem, nil
	case "app", "application":
		return LogApplication, nil
	case "service":
		return LogService, nil
	case "security":
		return LogSecurity, nil
	default:
		return LogSystem, errors.New("unknown scheme")
	}
}

var (
	DefaultLogger Logger = &consoleLogger{w: os.Stderr}
)

const stampFormat = "2006/01/02 15:04:05.000000"

type consoleLogger struct {
	w io.Writer
	d bool
}

func (c consoleLogger) stamp(t time.Time) {
	fmt.Fprintf(c.w, "%s ", t.Format(stampFormat))
}

func (c consoleLogger) print(format string, v ...interface{}) {
	c.stamp(time.Now())
	fmt.Fprintf(c.w, format, v...)
	if format[len(format)-1] != '\n' {
		fmt.Fprintf(c.w, "\n")
	}
}

func (c consoleLogger) Debug(format string, v ...interface{}) {
	if c.d {
		c.print(logHeaderDebug+format, v...)
	}
}

func (c consoleLogger) Info(format string, v ...interface{}) {
	c.print(logHeaderInfo+format, v...)
}

func (c consoleLogger) Warning(format string, v ...interface{}) {
	c.print(logHeaderWarning+format, v...)
}

func (c consoleLogger) Err(format string, v ...interface{}) {
	c.print(logHeaderErr+format, v...)
}

func (c consoleLogger) Crit(format string, v ...interface{}) {
	c.print(logHeaderCrit+format, v...)
	os.Exit(2)
}

// SetDebug can enable/disable debug mode.
func (c *consoleLogger) SetDebug(status bool) {
	c.d = status
}

// SetOutputURL is to set output for netlog.
func SetOutputURL(s string, debug ...bool) (err error) {
	var u *url.URL
	if u, err = url.Parse(s); err != nil {
		return
	}

	var isDebug bool
	if len(debug) > 0 {
		isDebug = debug[0]
	}

	q := u.Query()
	facility := LogApplication
	t := q.Get("facility")
	if t != "" {
		if facility, err = parseFacility(t); err != nil {
			return
		}
	}
	tag := q.Get("tag")

	switch u.Scheme {
	case "file":
		// file:///var/log/xxx.log
		var fp *os.File
		if fp, err = os.OpenFile(u.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err != nil {
			return
		}
		DefaultLogger = &consoleLogger{w: fp, d: isDebug}
		return nil
	case "net":
		// net:///?facility=x&tag=x
		DefaultLogger, err = NewLogger(facility, tag, isDebug)
		return
	case "tcp":
		// tcp://localhost:port/?facility=x&tag=x
		DefaultLogger, err = NewLogger(facility, tag, isDebug, u.Host)
		return
	case "tcp4", "tcp6":
		return errors.New("not implemented")
	default:
		// ??
		return errors.New("unsupported scheme")
	}
}

// Debug outputs debug level log output.
// It is usually used to output detailed debug information.
// If debug status is not enabled, no output is generated.
// This log level need not be treated as an anomaly.
// The debug status can be set from SetDebug.
func Debug(format string, v ...interface{}) {
	DefaultLogger.Debug(format, v...)
}

// Info outputs information level log output.
// It is usually used to output interesting events.
// This log level need not be treated as an anomaly.
func Info(format string, v ...interface{}) {
	DefaultLogger.Info(format, v...)
}

// Warning outputs warning level log output.
// It is usually used to output exceptional occurrences that are not errors.
// This log level need not be treated as an anomaly.
func Warning(format string, v ...interface{}) {
	DefaultLogger.Warning(format, v...)
}

// Err outputs error level log output.
// It is usually used to output execution-time errors that do not require
// immediate action but should typically be logged and monitored.
// This log level need be treated as an anomaly.
func Err(format string, v ...interface{}) {
	DefaultLogger.Err(format, v...)
}

// Crit outputs critical level log output.
// It is usually used to output critical conditions.
// When this method is executed, the process abends after outputting the log.
// This log level need be treated as an anomaly.
func Crit(format string, v ...interface{}) {
	DefaultLogger.Crit(format, v...)
}

// SetDebug can enable/disable debug mode.
func SetDebug(state bool) {
	DefaultLogger.SetDebug(state)
}
