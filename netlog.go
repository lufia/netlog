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
	LOG_SYSTEM      = facilitySystem
	LOG_APPLICATION = facilityApplication
	LOG_SERVICE     = facilityService
	LOG_SECURITY    = facilitySecurity
)

func parseFacility(s string) (Facility, error) {
	switch s {
	case "sys", "system":
		return LOG_SYSTEM, nil
	case "app", "application":
		return LOG_APPLICATION, nil
	case "service":
		return LOG_SERVICE, nil
	case "security":
		return LOG_SECURITY, nil
	default:
		return LOG_SYSTEM, errors.New("unknown scheme")
	}
}

var (
	DefaultLogger Logger = consoleLogger{w: os.Stderr}
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
	facility := LOG_APPLICATION
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
		DefaultLogger = consoleLogger{w: fp, d: isDebug}
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

func Debug(format string, v ...interface{}) {
	DefaultLogger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	DefaultLogger.Info(format, v...)
}

func Warning(format string, v ...interface{}) {
	DefaultLogger.Warning(format, v...)
}

func Err(format string, v ...interface{}) {
	DefaultLogger.Err(format, v...)
}

func Crit(format string, v ...interface{}) {
	DefaultLogger.Crit(format, v...)
}
