package rzerolog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

var (
	_ io.Writer  = (*ConsoleWriter)(nil)
	_ FileWriter = (*ConsoleWriter)(nil)

	consoleBufPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 100))
		},
	}
)

const (
	consoleDefaultTimeFormat = DefaultTimeFormat
)

// Formatter transforms the input into a formatted string.
type Formatter func(interface{}) string

// ConsoleWriter parses the JSON input and writes it in an
// (optionally) colorized, human-friendly format to Out.
type ConsoleWriter struct {
	Enable bool
	// Out is the output destination.
	Out io.Writer

	// NoColor disables the colorized output.
	NoColor bool

	// TimeFormat specifies the format for timestamp in output.
	TimeFormat string

	// PartsOrder defines the order of parts in output.
	PartsOrder []string

	// PartsExclude defines parts to not display in output.
	PartsExclude []string

	FormatTimestamp     Formatter
	FormatLevel         Formatter
	FormatLabel         Formatter
	FormatCaller        Formatter
	FormatMessage       Formatter
	FormatFieldName     Formatter
	FormatFieldValue    Formatter
	FormatErrFieldName  Formatter
	FormatErrFieldValue Formatter
}

type ConsoleWriterOption func(w *ConsoleWriter)

func EnableConsoleWriter() ConsoleWriterOption {
	return func(w *ConsoleWriter) {
		w.Enable = true
	}
}

func WithNoColor() ConsoleWriterOption {
	return func(w *ConsoleWriter) {
		w.NoColor = true
	}
}

// NewConsoleWriter creates and initializes a new ConsoleWriter.
func NewConsoleWriter(options ...ConsoleWriterOption) ConsoleWriter {
	w := ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: consoleDefaultTimeFormat,
		PartsOrder: consoleDefaultPartsOrder(),
	}

	for _, opt := range options {
		opt(&w)
	}

	return w
}

func (w *ConsoleWriter) SetOutput(file *os.File) error {
	w.Out = file
	return nil
}

// Write transforms the JSON input with formatters and appends to w.Out.
func (w *ConsoleWriter) Write(p []byte) (n int, err error) {
	if !w.Enable {
		return len(p), nil
	}
	if w.PartsOrder == nil {
		w.PartsOrder = consoleDefaultPartsOrder()
	}

	var buf = consoleBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		consoleBufPool.Put(buf)
	}()

	var evt map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&evt)
	if err != nil {
		return n, fmt.Errorf("cannot decode event: %s", err)
	}

	for _, p := range w.PartsOrder {
		w.writePart(buf, evt, p)
	}

	w.writeFields(evt, buf)

	err = buf.WriteByte('\n')
	if err != nil {
		return n, err
	}
	_, err = buf.WriteTo(w.Out)
	return len(p), err
}

// writeFields appends formatted key-value pairs to buf.
func (w *ConsoleWriter) writeFields(evt map[string]interface{}, buf *bytes.Buffer) {
	var fields = make([]string, 0, len(evt))
	for field := range evt {
		switch field {
		case zerolog.LevelFieldName,
			zerolog.TimestampFieldName,
			zerolog.MessageFieldName,
			zerolog.CallerFieldName,
			LabelFieldName:
			continue
		}
		fields = append(fields, field)
	}
	sort.Strings(fields)

	if len(fields) > 0 {
		buf.WriteByte(' ')
	}

	// Move the "error" field to the front
	ei := sort.Search(len(fields), func(i int) bool { return fields[i] >= zerolog.ErrorFieldName })
	if ei < len(fields) && fields[ei] == zerolog.ErrorFieldName {
		fields[ei] = ""
		fields = append([]string{zerolog.ErrorFieldName}, fields...)
		var xfields = make([]string, 0, len(fields))
		for _, field := range fields {
			if field == "" { // Skip empty fields
				continue
			}
			xfields = append(xfields, field)
		}
		fields = xfields
	}

	for i, field := range fields {
		var fn Formatter
		var fv Formatter

		if field == zerolog.ErrorFieldName {
			if w.FormatErrFieldName == nil {
				fn = consoleDefaultFormatErrFieldName(w.NoColor)
			} else {
				fn = w.FormatErrFieldName
			}

			if w.FormatErrFieldValue == nil {
				fv = consoleDefaultFormatErrFieldValue(w.NoColor)
			} else {
				fv = w.FormatErrFieldValue
			}
		} else {
			if w.FormatFieldName == nil {
				fn = consoleDefaultFormatFieldName(w.NoColor)
			} else {
				fn = w.FormatFieldName
			}

			if w.FormatFieldValue == nil {
				fv = consoleDefaultFormatFieldValue
			} else {
				fv = w.FormatFieldValue
			}
		}

		buf.WriteString(fn(field))

		switch fValue := evt[field].(type) {
		case string:
			if needsQuote(fValue) {
				buf.WriteString(fv(strconv.Quote(fValue)))
			} else {
				buf.WriteString(fv(fValue))
			}
		case json.Number:
			buf.WriteString(fv(fValue))
		default:
			b, err := json.Marshal(fValue)
			if err != nil {
				fmt.Fprintf(buf, colorize("[error: %v]", colorRed, w.NoColor), err)
			} else {
				fmt.Fprint(buf, fv(b))
			}
		}

		if i < len(fields)-1 { // Skip space for last field
			buf.WriteByte(' ')
		}
	}
}

// writePart appends a formatted part to buf.
func (w *ConsoleWriter) writePart(buf *bytes.Buffer, evt map[string]interface{}, p string) {
	var f Formatter

	if w.PartsExclude != nil && len(w.PartsExclude) > 0 {
		for _, exclude := range w.PartsExclude {
			if exclude == p {
				return
			}
		}
	}

	switch p {
	case zerolog.LevelFieldName:
		if w.FormatLevel == nil {
			f = consoleDefaultFormatLevel(w.NoColor)
		} else {
			f = w.FormatLevel
		}
	case zerolog.TimestampFieldName:
		if w.FormatTimestamp == nil {
			f = consoleDefaultFormatTimestamp(w.TimeFormat, w.NoColor)
		} else {
			f = w.FormatTimestamp
		}
	case zerolog.MessageFieldName:
		if w.FormatMessage == nil {
			f = consoleDefaultFormatMessage
		} else {
			f = w.FormatMessage
		}
	case LabelFieldName:
		if w.FormatLabel == nil {
			f = consoleDefaultFormatLabel(w.NoColor)
		} else {
			f = w.FormatMessage
		}
	case zerolog.CallerFieldName:
		if w.FormatCaller == nil {
			f = consoleDefaultFormatCaller(w.NoColor)
		} else {
			f = w.FormatCaller
		}
	default:
		if w.FormatFieldValue == nil {
			f = consoleDefaultFormatFieldValue
		} else {
			f = w.FormatFieldValue
		}
	}

	var s = f(evt[p])

	if len(s) > 0 {
		buf.WriteString(s)
		if p != w.PartsOrder[len(w.PartsOrder)-1] { // Skip space for last part
			buf.WriteByte(' ')
		}
	}
}

// needsQuote returns true when the string s should be quoted in output.
func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}

// colorize returns the string s wrapped in ANSI code c, unless disabled is true.
func colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

// ----- DEFAULT FORMATTERS ---------------------------------------------------

func consoleDefaultPartsOrder() []string {
	return []string{
		zerolog.TimestampFieldName,
		zerolog.LevelFieldName,
		LabelFieldName,
		zerolog.CallerFieldName,
		zerolog.MessageFieldName,
	}
}

func consoleDefaultFormatTimestamp(timeFormat string, noColor bool) Formatter {
	if timeFormat == "" {
		timeFormat = consoleDefaultTimeFormat
	}
	return func(i interface{}) string {
		t := "<nil>"
		switch tt := i.(type) {
		case string:
			ts, err := time.Parse(zerolog.TimeFieldFormat, tt)
			if err != nil {
				t = tt
			} else {
				t = ts.Format(timeFormat)
			}
		case json.Number:
			i, err := tt.Int64()
			if err != nil {
				t = tt.String()
			} else {
				var sec, nsec int64 = i, 0
				switch zerolog.TimeFieldFormat {
				case zerolog.TimeFormatUnixMs:
					nsec = int64(time.Duration(i) * time.Millisecond)
					sec = 0
				case zerolog.TimeFormatUnixMicro:
					nsec = int64(time.Duration(i) * time.Microsecond)
					sec = 0
				}
				ts := time.Unix(sec, nsec).UTC()
				t = ts.Format(timeFormat)
			}
		}
		return colorize(t, colorDarkGray, noColor)
	}
}

func consoleDefaultFormatLevel(noColor bool) Formatter {
	return func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "trace":
				l = colorize("TRC", colorMagenta, noColor)
			case "debug":
				l = colorize("DBG", colorYellow, noColor)
			case "info":
				l = colorize("INF", colorGreen, noColor)
			case "warn":
				l = colorize("WRN", colorRed, noColor)
			case "error":
				l = colorize(colorize("ERR", colorRed, noColor), colorBold, noColor)
			case "fatal":
				l = colorize(colorize("FTL", colorRed, noColor), colorBold, noColor)
			case "panic":
				l = colorize(colorize("PNC", colorRed, noColor), colorBold, noColor)
			default:
				l = colorize("???", colorBold, noColor)
			}
		} else {
			if i == nil {
				l = colorize("???", colorBold, noColor)
			} else {
				l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
			}
		}
		return l
	}
}

func consoleDefaultFormatLabel(noColor bool) Formatter {
	return func(i interface{}) string {
		if i == nil {
			return ""
		}
		return colorize(fmt.Sprintf("%s", i), colorBlue, noColor)
	}
}

func consoleDefaultFormatCaller(noColor bool) Formatter {
	return func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			if cwd, err := os.Getwd(); err == nil {
				if rel, err := filepath.Rel(cwd, c); err == nil {
					ss := strings.Split(cwd, string(filepath.Separator))
					ssl := len(ss)
					if ssl > 0 {
						c = ss[ssl-1] + "/" + rel
					} else {
						c = rel
					}
				}
			}
			c = colorize(c, colorBold, noColor) + colorize(" >", colorCyan, noColor)
		}
		return c
	}
}

func consoleDefaultFormatMessage(i interface{}) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%s", i)
}

func consoleDefaultFormatFieldName(noColor bool) Formatter {
	return func(i interface{}) string {
		return colorize(fmt.Sprintf("%s=", i), colorCyan, noColor)
	}
}

func consoleDefaultFormatFieldValue(i interface{}) string {
	return fmt.Sprintf("%s", i)
}

func consoleDefaultFormatErrFieldName(noColor bool) Formatter {
	return func(i interface{}) string {
		return colorize(fmt.Sprintf("%s=", i), colorRed, noColor)
	}
}

func consoleDefaultFormatErrFieldValue(noColor bool) Formatter {
	return func(i interface{}) string {
		return colorize(fmt.Sprintf("%s", i), colorRed, noColor)
	}
}
