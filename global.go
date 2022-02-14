package rzerolog

const (
	LabelFieldName = "label"
)

const (
	DefaultTimeFormat   = "2006-01-02 15:04:05.000"
	DefaultFilePath     = "./"
	DefaultFileName     = "RZero.log"
	DefaultFileSize     = 100 * 1 << 10 // 100M
	DefaultMaxFileCount = 10
	DefaultLevel        = DebugLevel

	LogFormatJSON        = "json"
	LogFormatConsoleText = "text"
	DefaultLogFormat     = LogFormatJSON
)
