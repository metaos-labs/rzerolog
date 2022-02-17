package rzerolog

import (
	"github.com/rs/zerolog"
)

type Level zerolog.Level

const (
	// DebugLevel defines debug log level.
	DebugLevel = Level(zerolog.DebugLevel)
	// InfoLevel defines info log level.
	InfoLevel = Level(zerolog.InfoLevel)
	// WarnLevel defines warn log level.
	WarnLevel = Level(zerolog.WarnLevel)
	// ErrorLevel defines error log level.
	ErrorLevel = Level(zerolog.ErrorLevel)
	// FatalLevel defines fatal log level.
	FatalLevel = Level(zerolog.FatalLevel)
	// PanicLevel defines panic log level.
	PanicLevel = Level(zerolog.PanicLevel)
	// NoLevel defines an absent log level.
	NoLevel = Level(zerolog.NoLevel)
	// Disabled disables the logger.
	Disabled = Level(zerolog.Disabled)

	// TraceLevel defines trace log level.
	TraceLevel = Level(zerolog.TraceLevel)
)

func init() {
	zerolog.TimeFieldFormat = DefaultTimeFormat
}

// RZeroLogger is a logger wrapped with zerolog.
type RZeroLogger struct {
	zerolog.Logger
	label string
}

func newRZeroLogger(cfg loggerPrepare) *RZeroLogger {
	if err := cfg.fw.initBase(); err != nil {
		panic(err)
	}
	multi := zerolog.MultiLevelWriter(cfg.cw, cfg.fw)
	ctx := zerolog.New(multi).Level(zerolog.Level(cfg.level)).With().Timestamp()
	if cfg.caller {
		ctx = ctx.Caller()
	}
	zeroLog := ctx.Logger()
	return &RZeroLogger{
		Logger: zeroLog,
		label:  cfg.label,
	}
}

func NewDefaultRZeroLogger() *RZeroLogger {
	cfg := defaultConfig()
	if err := cfg.fw.initBase(); err != nil {
		panic(err)
	}
	return newRZeroLogger(cfg)
}

func NewRZeroLogger(opts ...Option) *RZeroLogger {
	cfg := defaultConfig()
	cfg.apply(opts...)
	return newRZeroLogger(cfg)
}

// GetLabeledSubLogger create a new sub logger with a new label given.
//
// The internal logger is the same as parent.
func (l *RZeroLogger) GetLabeledSubLogger(label string) *RZeroLogger {
	return &RZeroLogger{
		Logger: l.Logger,
		label:  label,
	}
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func (l *RZeroLogger) Trace() *Event {
	return &Event{
		Event: l.Logger.Trace(),
		label: l.label,
	}
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func (l *RZeroLogger) Debug() *Event {
	return &Event{
		Event: l.Logger.Debug(),
		label: l.label,
	}
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func (l *RZeroLogger) Info() *Event {
	return &Event{
		Event: l.Logger.Info(),
		label: l.label,
	}
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func (l *RZeroLogger) Warn() *Event {
	return &Event{
		Event: l.Logger.Warn(),
		label: l.label,
	}
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func (l *RZeroLogger) Error() *Event {
	return &Event{
		Event: l.Logger.Error(),
		label: l.label,
	}
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func (l *RZeroLogger) Err(err error) *Event {
	return &Event{
		Event: l.Logger.Err(err),
		label: l.label,
	}
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method, which terminates the program immediately.
//
// You must call Msg on the returned event in order to send the event.
func (l *RZeroLogger) Fatal() *Event {
	return &Event{
		Event: l.Logger.Fatal(),
		label: l.label,
	}
}

// Panic starts a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine.
//
// You must call Msg on the returned event in order to send the event.
func (l *RZeroLogger) Panic() *Event {
	return &Event{
		Event: l.Logger.Panic(),
		label: l.label,
	}
}

// WithLevel starts a new message with level. Unlike Fatal and Panic
// methods, WithLevel does not terminate the program or stop the ordinary
// flow of a gourotine when used with their respective levels.
//
// You must call Msg on the returned event in order to send the event.
func (l *RZeroLogger) WithLevel(level Level) *Event {
	return &Event{
		Event: l.Logger.WithLevel(zerolog.Level(level)),
		label: l.label,
	}
}

// Log starts a new message with no level. Setting GlobalLevel to Disabled
// will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func (l *RZeroLogger) Log() *Event {
	return &Event{
		Event: l.Logger.Log(),
		label: l.label,
	}
}
