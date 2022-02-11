package log

import (
	"errors"
	"github.com/sophon-labs/rzerolog"
	"os"
	"path/filepath"
)

func setup() {
	logger := rzerolog.NewRZeroLogger(rzerolog.EnableLogFiles(), rzerolog.WithLevel(rzerolog.TraceLevel))
	SetGlobalLogger(logger)
}

// Example setting global logger with a custom one.
func ExampleSetGlobalLogger() {
	logger := rzerolog.NewRZeroLogger(
		rzerolog.WithLevel(rzerolog.InfoLevel),
		rzerolog.WithLabel("Label name"),
	)
	SetGlobalLogger(logger)
}

func ExamplePrint() {
	setup()

	Print("hello world")
	// Output:
	// Console: 2022-02-11 16:37:40.789 CST +0800 DBG log/log_example_test.go:26 > hello world
	// File(TEXT): 2022-02-11 16:37:40.789 CST +0800 DBG log/log_example_test.go:26 > hello world
	// File(JSON): {"level":"debug","time":"2022-02-11 16:37:40.789 CST +0800","caller":"rzerolog/log/log_example_test.go:26","message":"hello world"}
}

func ExamplePrintf() {
	setup()

	Printf("hello %s", "world")
	// Output:
	// Console: 2022-02-11 16:22:44.021 CST +0800 DBG log/log_example_test.go:35 > hello world
	// File(TEXT): 2022-02-11 16:22:44.021 CST +0800 DBG log/log_example_test.go:35 > hello world
	// File(JSON): {"level":"debug","time":"2022-02-11 16:22:44.021 CST +0800","caller":"rzerolog/log/log_example_test.go:35","message":"hello world"}
}

func ExampleLog() {
	setup()

	Log().Msg("hello world")
	// Output:
	// Console: 2022-02-11 16:25:10.427 CST +0800 ??? log/log_example_test.go:43 > hello world
	// File(TEXT): 2022-02-11 16:25:10.427 CST +0800 ??? log/log_example_test.go:43 > hello world
	// File(JSON): {"time":"2022-02-11 16:25:10.427 CST +0800","caller":"rzerolog/log/log_example_test.go:43","message":"hello world"}
}

func ExampleErr() {
	setup()

	err := errors.New("some error")
	Err(err).Msg("hello world")
	Err(nil).Msg("hello world")
	// Output:
	// Console:
	// 2022-02-11 16:28:13.651 CST +0800 ERR log/log_example_test.go:57 > hello world error="some error"
	// 2022-02-11 16:28:13.689 CST +0800 INF log/log_example_test.go:58 > hello world
	// File(TEXT):
	// 2022-02-11 16:28:13.651 CST +0800 ERR log/log_example_test.go:57 > hello world error="some error"
	// 2022-02-11 16:28:13.689 CST +0800 INF log/log_example_test.go:58 > hello world
	// File(JSON):
	// {"level":"error","error":"some error","time":"2022-02-11 16:28:13.651 CST +0800","caller":"rzerolog/log/log_example_test.go:57","message":"hello world"}
	// {"level":"info","time":"2022-02-11 16:28:13.689 CST +0800","caller":"rzerolog/log/log_example_test.go:58","message":"hello world"}
}

func ExampleTrace() {
	setup()

	Trace().Msg("hello world")
	// Output:
	// Console:
	// 2022-02-11 16:33:25.613 CST +0800 TRC log/log_example_test.go:74 > hello world
	// File(TEXT):
	// 2022-02-11 16:33:25.613 CST +0800 TRC log/log_example_test.go:74 > hello world
	// File(JSON):
	// {"level":"trace","time":"2022-02-11 16:33:25.613 CST +0800","caller":"rzerolog/log/log_example_test.go:74","message":"hello world"}
}

func ExampleDebug() {
	setup()

	Debug().Msg("hello world")
	// Output:
	// Console:
	// 2022-02-11 16:35:42.250 CST +0800 DBG log/log_example_test.go:87 > hello world
	// File(TEXT):
	// 2022-02-11 16:35:42.250 CST +0800 DBG log/log_example_test.go:87 > hello world
	// File(JSON):
	// {"level":"debug","time":"2022-02-11 16:35:42.250 CST +0800","caller":"rzerolog/log/log_example_test.go:87","message":"hello world"}
}

func ExampleInfo() {
	setup()

	Info().Msg("hello world")
	// Output:
	// Console:
	// 2022-02-11 16:38:47.684 CST +0800 INF log/log_example_test.go:100 > hello world
	// File(TEXT):
	// 2022-02-11 16:38:47.684 CST +0800 INF log/log_example_test.go:100 > hello world
	// File(JSON):
	// {"level":"info","time":"2022-02-11 16:38:47.684 CST +0800","caller":"rzerolog/log/log_example_test.go:100","message":"hello world"}
}

func ExampleWarn() {
	setup()

	Warn().Msg("hello world")
	// Output:
	// Console:
	// 2022-02-11 16:39:44.946 CST +0800 WRN log/log_example_test.go:113 > hello world
	// File(TEXT):
	// 2022-02-11 16:39:44.946 CST +0800 WRN log/log_example_test.go:113 > hello world
	// File(JSON):
	// {"level":"warn","time":"2022-02-11 16:39:44.946 CST +0800","caller":"rzerolog/log/log_example_test.go:113","message":"hello world"}
}

func ExampleError() {
	setup()

	Error().Msg("hello world")
	// Output:
	// Console:
	// 2022-02-11 16:41:10.384 CST +0800 ERR log/log_example_test.go:126 > hello world
	// File(TEXT):
	// 2022-02-11 16:41:10.384 CST +0800 ERR log/log_example_test.go:126 > hello world
	// File(JSON):
	// {"level":"error","time":"2022-02-11 16:41:10.384 CST +0800","caller":"rzerolog/log/log_example_test.go:126","message":"hello world"}
}

func ExampleFatal() {
	setup()

	Fatal().Err(errors.New("some error")).Str("some key", "key value").Msg("hello world")
	// Output:
	// File(JSON):
	// {"level":"fatal","error":"some error","some key":"key value","time":"2022-02-11 16:43:16.137 CST +0800","caller":"rzerolog/log/log_example_test.go:139","message":"hello world"}
}

func ExamplePanic() {
	setup()

	Panic().Err(errors.New("some error")).Str("some key", "key value").Msg("hello world")
	// Output:
	// File(JSON):
	// {"level":"panic","error":"some error","some key":"key value","time":"2022-02-11 16:45:11.172 CST +0800","caller":"rzerolog/log/log_example_test.go:148","message":"hello world"}
}

func Example() {
	// get global logger
	logger := GlobalLogger()

	// new default logger
	logger = rzerolog.NewDefaultRZeroLogger()

	// new custom logger
	logger = rzerolog.NewRZeroLogger(
		// All Options below are optional.
		// set level
		rzerolog.WithLevel(rzerolog.InfoLevel),
		// disable console output
		rzerolog.DisableConsolePrint(),
		// console printing with no color
		rzerolog.WithNoConsolePrintColor(),
		// write logs to files
		rzerolog.EnableLogFiles(),
		// log file name
		rzerolog.WithLogFileName("filename.log"),
		// log files path
		rzerolog.WithLogFilePath("../testdata"),
		// time rolling rules of log files
		//rzerolog.WithLogFileName("yyyyMMdd.log"),
		rzerolog.EnableTimeRolling(),
		// size rolling rules of log files
		rzerolog.WithSizeRolling(100<<10, 10),
		// no caller info printed
		rzerolog.WithNoCaller(),
		// log record format in log file (not in console)
		rzerolog.WithLogFormat(rzerolog.LogFormatJSON),
		// set label name, usually used to mark modules
		rzerolog.WithLabel("label name"),
	)

	// print log
	logger.Info().Msg("info message")
	// Output: {"level":"info","label":"label name","time":"2022-02-11 17:39:11.597 CST +0800","message":"info message"}

	// get labeled sub logger
	subLogger := logger.GetLabeledSubLogger("other")
	subLogger.Info().Msg("info message")
	// Output: {"level":"info","label":"other","time":"2022-02-11 17:39:11.597 CST +0800","message":"info message"}

	defer func() {
		_ = os.RemoveAll(filepath.Join("../testdata/filename.log"))
	}()
}
