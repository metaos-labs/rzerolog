package log

import (
	"fmt"
	"github.com/sophon-labs/rzerolog"
)

var (
	_log *rzerolog.RZeroLogger
)

func init() {
	_log = rzerolog.NewDefaultRZeroLogger()
}

// GlobalLogger get the global logger instance.
func GlobalLogger() *rzerolog.RZeroLogger {
	return _log
}

// SetGlobalLogger set the global logger with a custom logger.
func SetGlobalLogger(logger *rzerolog.RZeroLogger) {
	_log = logger
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func Err(err error) *rzerolog.Event {
	return _log.Err(err)
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func Trace() *rzerolog.Event {
	return _log.Trace()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *rzerolog.Event {
	return _log.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *rzerolog.Event {
	return _log.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *rzerolog.Event {
	return _log.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *rzerolog.Event {
	return _log.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *rzerolog.Event {
	return _log.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *rzerolog.Event {
	return _log.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func WithLevel(level rzerolog.Level) *rzerolog.Event {
	return _log.WithLevel(level)
}

// Log starts a new message with no level. Setting rzerolog.GlobalLevel to
// rzerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func Log() *rzerolog.Event {
	return _log.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	_log.Debug().CallerSkipFrame(1).Msg(fmt.Sprint(v...))
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	_log.Debug().CallerSkipFrame(1).Msgf(format, v...)
}
