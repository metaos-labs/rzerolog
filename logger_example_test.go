package rzerolog

func Example() {
	// new default logger
	logger := NewDefaultRZeroLogger()

	// new custom logger
	logger = NewRZeroLogger(
		// All Options below are optional.
		// set level
		WithLevel(InfoLevel),
		// disable console output
		DisableConsolePrint(),
		// console printing with no color
		WithNoConsolePrintColor(),
		// write logs to files
		EnableLogFiles(),
		// log file name
		WithLogFileName("filename.log"),
		// log files path
		WithLogFilePath("./testdata"),
		// time rolling rules of log files
		//rzerolog.WithLogFileName("yyyyMMdd.log"),
		EnableTimeRolling(),
		// size rolling rules of log files
		WithSizeRolling(100<<10, 10),
		// no caller info printed
		WithNoCaller(),
		// log record format in log file (not in console)
		WithLogFormat(LogFormatJSON),
		// set label name, usually used to mark modules
		WithLabel("label name"),
	)

	// print log
	logger.Info().Msg("info message")
	// Output: {"level":"info","label":"label name","time":"2022-02-11 17:39:11.597 CST +0800","message":"info message"}

	// get labeled sub logger
	subLogger := logger.GetLabeledSubLogger("other")
	subLogger.Info().Msg("info message")
	// Output: {"level":"info","label":"other","time":"2022-02-11 17:39:11.597 CST +0800","message":"info message"}

	// see the log file : ./testdata/filename.log
}
