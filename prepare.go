package rzerolog

import (
	"fmt"
	"os"
)

type loggerPrepare struct {
	cw *ConsoleWriter
	fw *LogFileWriter

	level     Level
	logFormat string
	label     string
	caller    bool
}

func defaultConfig() loggerPrepare {
	consoleWriter := &ConsoleWriter{
		Enable:  true,
		Out:     os.Stdout,
		NoColor: false,
	}

	fw := &LogFileWriter{
		enable:          false,
		writer:          &osFileWriter{},
		fillPath:        DefaultFilePath,
		timeRolling:     false,
		sizeRolling:     false,
		fileSize:        DefaultFileSize,
		maxFileCount:    DefaultMaxFileCount,
		logFileName:     DefaultFileName,
		currentFileName: DefaultFileName,
		file:            nil,
	}
	cfg := loggerPrepare{
		cw:        consoleWriter,
		fw:        fw,
		level:     DefaultLevel,
		logFormat: DefaultLogFormat,
		label:     "",
		caller:    true,
	}
	return cfg
}

type Option func(cfg *loggerPrepare)

func (lc *loggerPrepare) apply(opts ...Option) {
	for _, opt := range opts {
		opt(lc)
	}
}

// WithLevel set logger level.
func WithLevel(l Level) Option {
	return func(cfg *loggerPrepare) {
		cfg.level = l
	}
}

// DisableConsolePrint will disable console print.
func DisableConsolePrint() Option {
	return func(cfg *loggerPrepare) {
		cfg.cw.Enable = false
	}
}

// WithNoConsolePrintColor will print console with no color font.
func WithNoConsolePrintColor() Option {
	return func(cfg *loggerPrepare) {
		cfg.cw.NoColor = true
	}
}

// EnableLogFiles will make logger to write logs to log files.
func EnableLogFiles() Option {
	return func(cfg *loggerPrepare) {
		cfg.fw.enable = true
	}
}

// WithLogFilePath set the path which log files will be written to.
func WithLogFilePath(path string) Option {
	return func(cfg *loggerPrepare) {
		cfg.fw.fillPath = path
	}
}

// WithLogFileName set the filename of log files.
//
// NOTE: If EnableTimeRolling() invoked, the final log file name will be parsed by the time format parser.
// eg:
// "yyyyMMddHH.log" => "2022021116.log"
func WithLogFileName(name string) Option {
	return func(cfg *loggerPrepare) {
		cfg.fw.logFileName = name
	}
}

// EnableTimeRolling enable rolling the log files on rules implicit in LogFileName set.
// eg:
// "yyyyMMddHH.log" => "2022021116.log"
func EnableTimeRolling() Option {
	return func(cfg *loggerPrepare) {
		cfg.fw.timeRolling = true
		cfg.fw.enable = true
	}
}

// WithSizeRolling enable rolling the log files on rules bounded by file size.
// This option helps prevent log files from taking up too much disk space.
// When the number of log files cut reaches the threshold,
// the redundant old log files will be automatically removed.
// NOTE: The unit of the filesize parameter is Kb.
func WithSizeRolling(fileSize int64, maxFileCount int) Option {
	return func(cfg *loggerPrepare) {
		cfg.fw.sizeRolling = true
		cfg.fw.enable = true
		cfg.fw.maxFileCount = maxFileCount
		cfg.fw.fileSize = fileSize * 1 << 10
	}
}

// WithNoCaller will prevent the logger caller information from printing.
func WithNoCaller() Option {
	return func(cfg *loggerPrepare) {
		cfg.caller = false
	}
}

// WithLogFormat set the output format when logger printing.
// Current supporting:"text","json"
func WithLogFormat(format string) Option {
	var w FileWriter = &osFileWriter{}
	switch format {
	case LogFormatJSON:
	case LogFormatConsoleText:
		w = &ConsoleWriter{Enable: true, NoColor: true, Out: w}
	default:
		panic(fmt.Sprintln("unsupported log format. supporting:", LogFormatJSON, LogFormatConsoleText))
	}
	return func(cfg *loggerPrepare) {
		cfg.fw.writer = w
	}
}

// WithLabel set the logger label.
// The label will be print to log records automatically.
// It is usually used to mark modules.
func WithLabel(label string) Option {
	return func(cfg *loggerPrepare) {
		cfg.label = label
	}
}
