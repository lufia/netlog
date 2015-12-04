package netlog

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
