package log

import (
	log "github.com/cihub/seelog"
)

// InitLogger initialize logger
func InitLogger(opts string) (err error) {
	var logger log.LoggerInterface
	if opts == "" {
		logger = log.Default
	} else {
		logger, err = log.LoggerFromConfigAsString(opts)
		if err != nil {
			return err
		}
	}
	return log.UseLogger(logger)
}

// Tracef formats message according to format specifier
// and writes to log with level = Trace.
func Tracef(format string, params ...interface{}) {
	log.Tracef(format, params...)
}

// Debugf formats message according to format specifier
// and writes to log with level = Debug.
func Debugf(format string, params ...interface{}) {
	log.Debugf(format, params...)
}

// Infof formats message according to format specifier
// and writes to log with level = Info.
func Infof(format string, params ...interface{}) {
	log.Infof(format, params...)
}

// Warnf formats message according to format specifier
// and writes to log with level = Warn.
func Warnf(format string, params ...interface{}) error {
	return log.Warnf(format, params...)
}

// Errorf formats message according to format specifier
// and writes to log with level = Error.
func Errorf(format string, params ...interface{}) error {
	return log.Errorf(format, params...)
}

// Criticalf formats message according to format specifier
// and writes to log with level = Critical.
func Criticalf(format string, params ...interface{}) error {
	return log.Criticalf(format, params...)
}

// Trace formats message using the default formats for its operands
// and writes to log with level = Trace
func Trace(v ...interface{}) {
	log.Trace(v...)
}

// Debug formats message using the default formats for its operands
// and writes to log with level = Debug
func Debug(v ...interface{}) {
	log.Debug(v...)
}

// Info formats message using the default formats for its operands
// and writes to log with level = Info
func Info(v ...interface{}) {
	log.Info(v...)
}

// Warn formats message using the default formats for its operands
// and writes to log with level = Warn
func Warn(v ...interface{}) error {
	return log.Warn(v...)
}

// Error formats message using the default formats for its operands
// and writes to log with level = Error
func Error(v ...interface{}) error {
	return log.Error(v...)
}

// Critical formats message using the default formats for its operands
// and writes to log with level = Critical
func Critical(v ...interface{}) error {
	return log.Critical(v...)
}
