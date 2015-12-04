package netlog

import (
	"log"
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

var (
	DefaultLogger Logger = consoleLogger{}
)

type consoleLogger struct{}

func (c consoleLogger) Info(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (c consoleLogger) Warning(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (c consoleLogger) Err(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (c consoleLogger) Crit(format string, v ...interface{}) {
	log.Printf(format, v...)
	os.Exit(2)
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
