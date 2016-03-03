package netlog

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
)

type Logger interface {
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Err(format string, v ...interface{})
	Crit(format string, v ...interface{})
}

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

type consoleLogger struct {
	w io.Writer
}

func (c consoleLogger) Info(format string, v ...interface{}) {
	fmt.Fprintf(c.w, format, v...)
	if format[len(format)-1] != '\n' {
		fmt.Fprintf(c.w, "\n")
	}
}

func (c consoleLogger) Warning(format string, v ...interface{}) {
	fmt.Fprintf(c.w, format, v...)
	if format[len(format)-1] != '\n' {
		fmt.Fprintf(c.w, "\n")
	}
}

func (c consoleLogger) Err(format string, v ...interface{}) {
	fmt.Fprintf(c.w, format, v...)
	if format[len(format)-1] != '\n' {
		fmt.Fprintf(c.w, "\n")
	}
}

func (c consoleLogger) Crit(format string, v ...interface{}) {
	fmt.Fprintf(c.w, format, v...)
	if format[len(format)-1] != '\n' {
		fmt.Fprintf(c.w, "\n")
	}
	os.Exit(2)
}

// SetOutputURL is to set output for netlog.
func SetOutputURL(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	q := u.Query()
	facility := LOG_APPLICATION
	t := q.Get("facility")
	if t != "" {
		facility, err = parseFacility(t)
		if err != nil {
			return err
		}
	}
	tag := q.Get("tag")

	switch u.Scheme {
	case "file":
		// file:///var/log/xxx.log
		fout, err := os.OpenFile(u.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		DefaultLogger = consoleLogger{w: fout}
		return nil
	case "net":
		// net:///?facility=x&tag=x
		DefaultLogger = NewLogger(facility, tag)
		return nil
	case "tcp", "tcp4", "tcp6":
		// tcp://localhost/?facility=x&tag=x
		return errors.New("not implemented")
	default:
		// ??
		return errors.New("unsupported scheme")
	}
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
