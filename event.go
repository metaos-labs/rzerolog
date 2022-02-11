package rzerolog

import (
	"fmt"
	"github.com/rs/zerolog"
	"net"
	"time"
)

type Event struct {
	*zerolog.Event
	label string
}

// Msg sends the *Event with msg added as the message field if not empty.
//
// NOTICE: once this method is called, the *Event should be disposed.
// Calling Msg twice can have unexpected result.
func (e *Event) Msg(msg string) {
	if e.label != "" {
		e.Event.CallerSkipFrame(1).Str(LabelFieldName, e.label).Msg(msg)
	} else {
		e.Event.CallerSkipFrame(1).Msg(msg)
	}
}

// Send is equivalent to calling Msg("").
//
// NOTICE: once this method is called, the *Event should be disposed.
func (e *Event) Send() {
	if e.label != "" {
		e.Event.CallerSkipFrame(1).Str(LabelFieldName, e.label).Msg("")
	} else {
		e.Event.CallerSkipFrame(1).Msg("")
	}
}

// Msgf sends the event with formatted msg added as the message field if not empty.
//
// NOTICE: once this method is called, the *Event should be disposed.
// Calling Msgf twice can have unexpected result.
func (e *Event) Msgf(format string, v ...interface{}) {
	if e.label != "" {
		e.Event.CallerSkipFrame(1).Str(LabelFieldName, e.label).Msg(fmt.Sprintf(format, v...))
	} else {
		e.Event.CallerSkipFrame(1).Msg(fmt.Sprintf(format, v...))
	}
}

// Discard disables the event so Msg(f) won't print it.
func (e *Event) Discard() *Event {
	e.Event.Discard()
	return e
}

// Fields is a helper function to use a map to set fields using type assertion.
func (e *Event) Fields(fields map[string]interface{}) *Event {
	e.Event.Fields(fields)
	return e
}

// Dict adds the field key with a dict to the event context.
// Use zerolog.Dict() to create the dictionary.
func (e *Event) Dict(key string, dict *Event) *Event {
	e.Event.Dict(key, dict.Event)
	return e
}

// Array adds the field key with an array to the event context.
// Use zerolog.Arr() to create the array or pass a type that
// implement the LogArrayMarshaler interface.
func (e *Event) Array(key string, arr zerolog.LogArrayMarshaler) *Event {
	e.Event.Array(key, arr)
	return e
}

// Object marshals an object that implement the LogObjectMarshaler interface.
func (e *Event) Object(key string, obj zerolog.LogObjectMarshaler) *Event {
	e.Event.Object(key, obj)
	return e
}

// Func allows an anonymous func to run only if the event is enabled.
func (e *Event) Func(f func(ee *Event)) *Event {
	f2 := func(_ *zerolog.Event) {
		f(e)
	}
	e.Event.Func(f2)
	return e
}

// EmbedObject marshals an object that implement the LogObjectMarshaler interface.
func (e *Event) EmbedObject(obj zerolog.LogObjectMarshaler) *Event {
	e.Event.EmbedObject(obj)
	return e
}

// Str adds the field key with val as a string to the *Event context.
func (e *Event) Str(key, val string) *Event {
	e.Event.Str(key, val)
	return e
}

// Strs adds the field key with vals as a []string to the *Event context.
func (e *Event) Strs(key string, vals []string) *Event {
	e.Event.Strs(key, vals)
	return e
}

// Stringer adds the field key with val.String() (or null if val is nil) to the *Event context.
func (e *Event) Stringer(key string, val fmt.Stringer) *Event {
	e.Event.Stringer(key, val)
	return e
}

// Bytes adds the field key with val as a string to the *Event context.
//
// Runes outside of normal ASCII ranges will be hex-encoded in the resulting
// JSON.
func (e *Event) Bytes(key string, val []byte) *Event {
	e.Event.Bytes(key, val)
	return e
}

// Hex adds the field key with val as a hex string to the *Event context.
func (e *Event) Hex(key string, val []byte) *Event {
	e.Event.Hex(key, val)
	return e
}

// RawJSON adds already encoded JSON to the log line under key.
//
// No sanity check is performed on b; it must not contain carriage returns and
// be valid JSON.
func (e *Event) RawJSON(key string, b []byte) *Event {
	e.Event.RawJSON(key, b)
	return e
}

// AnErr adds the field key with serialized err to the *Event context.
// If err is nil, no field is added.
func (e *Event) AnErr(key string, err error) *Event {
	e.Event.AnErr(key, err)
	return e
}

// Errs adds the field key with errs as an array of serialized errors to the
// *Event context.
func (e *Event) Errs(key string, errs []error) *Event {
	e.Event.Errs(key, errs)
	return e
}

// Err adds the field "error" with serialized err to the *Event context.
// If err is nil, no field is added.
//
// To customize the key name, change zerolog.ErrorFieldName.
//
// If Stack() has been called before and zerolog.ErrorStackMarshaler is defined,
// the err is passed to ErrorStackMarshaler and the result is appended to the
// zerolog.ErrorStackFieldName.
func (e *Event) Err(err error) *Event {
	e.Event.Err(err)
	return e
}

// Stack enables stack trace printing for the error passed to Err().
//
// ErrorStackMarshaler must be set for this method to do something.
func (e *Event) Stack() *Event {
	e.Event.Stack()
	return e
}

// Bool adds the field key with val as a bool to the *Event context.
func (e *Event) Bool(key string, b bool) *Event {
	e.Event.Bool(key, b)
	return e
}

// Bools adds the field key with val as a []bool to the *Event context.
func (e *Event) Bools(key string, b []bool) *Event {
	e.Event.Bools(key, b)
	return e
}

// Int adds the field key with i as a int to the *Event context.
func (e *Event) Int(key string, i int) *Event {
	e.Event.Int(key, i)
	return e
}

// Ints adds the field key with i as a []int to the *Event context.
func (e *Event) Ints(key string, i []int) *Event {
	e.Event.Ints(key, i)
	return e
}

// Int8 adds the field key with i as a int8 to the *Event context.
func (e *Event) Int8(key string, i int8) *Event {
	e.Event.Int8(key, i)
	return e
}

// Ints8 adds the field key with i as a []int8 to the *Event context.
func (e *Event) Ints8(key string, i []int8) *Event {
	e.Event.Ints8(key, i)
	return e
}

// Int16 adds the field key with i as a int16 to the *Event context.
func (e *Event) Int16(key string, i int16) *Event {
	e.Event.Int16(key, i)
	return e
}

// Ints16 adds the field key with i as a []int16 to the *Event context.
func (e *Event) Ints16(key string, i []int16) *Event {
	e.Event.Ints16(key, i)
	return e
}

// Int32 adds the field key with i as a int32 to the *Event context.
func (e *Event) Int32(key string, i int32) *Event {
	e.Event.Int32(key, i)
	return e
}

// Ints32 adds the field key with i as a []int32 to the *Event context.
func (e *Event) Ints32(key string, i []int32) *Event {
	e.Event.Ints32(key, i)
	return e
}

// Int64 adds the field key with i as a int64 to the *Event context.
func (e *Event) Int64(key string, i int64) *Event {
	e.Event.Int64(key, i)
	return e
}

// Ints64 adds the field key with i as a []int64 to the *Event context.
func (e *Event) Ints64(key string, i []int64) *Event {
	e.Event.Ints64(key, i)
	return e
}

// Uint adds the field key with i as a uint to the *Event context.
func (e *Event) Uint(key string, i uint) *Event {
	e.Event.Uint(key, i)
	return e
}

// Uints adds the field key with i as a []int to the *Event context.
func (e *Event) Uints(key string, i []uint) *Event {
	e.Event.Uints(key, i)
	return e
}

// Uint8 adds the field key with i as a uint8 to the *Event context.
func (e *Event) Uint8(key string, i uint8) *Event {
	e.Event.Uint8(key, i)
	return e
}

// Uints8 adds the field key with i as a []int8 to the *Event context.
func (e *Event) Uints8(key string, i []uint8) *Event {
	e.Event.Uints8(key, i)
	return e
}

// Uint16 adds the field key with i as a uint16 to the *Event context.
func (e *Event) Uint16(key string, i uint16) *Event {
	e.Event.Uint16(key, i)
	return e
}

// Uints16 adds the field key with i as a []int16 to the *Event context.
func (e *Event) Uints16(key string, i []uint16) *Event {
	e.Event.Uints16(key, i)
	return e
}

// Uint32 adds the field key with i as a uint32 to the *Event context.
func (e *Event) Uint32(key string, i uint32) *Event {
	e.Event.Uint32(key, i)
	return e
}

// Uints32 adds the field key with i as a []int32 to the *Event context.
func (e *Event) Uints32(key string, i []uint32) *Event {
	e.Event.Uints32(key, i)
	return e
}

// Uint64 adds the field key with i as a uint64 to the *Event context.
func (e *Event) Uint64(key string, i uint64) *Event {
	e.Event.Uint64(key, i)
	return e
}

// Uints64 adds the field key with i as a []int64 to the *Event context.
func (e *Event) Uints64(key string, i []uint64) *Event {
	e.Event.Uints64(key, i)
	return e
}

// Float32 adds the field key with f as a float32 to the *Event context.
func (e *Event) Float32(key string, f float32) *Event {
	e.Event.Float32(key, f)
	return e
}

// Floats32 adds the field key with f as a []float32 to the *Event context.
func (e *Event) Floats32(key string, f []float32) *Event {
	e.Event.Floats32(key, f)
	return e
}

// Float64 adds the field key with f as a float64 to the *Event context.
func (e *Event) Float64(key string, f float64) *Event {
	e.Event.Float64(key, f)
	return e
}

// Floats64 adds the field key with f as a []float64 to the *Event context.
func (e *Event) Floats64(key string, f []float64) *Event {
	e.Event.Floats64(key, f)
	return e
}

// Timestamp adds the current local time as UNIX timestamp to the *Event context with the "time" key.
// To customize the key name, change zerolog.TimestampFieldName.
//
// NOTE: It won't dedupe the "time" key if the *Event (or *Context) has one
// already.
func (e *Event) Timestamp() *Event {
	e.Event.Timestamp()
	return e
}

// Time adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (e *Event) Time(key string, t time.Time) *Event {
	e.Event.Time(key, t)
	return e
}

// Times adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (e *Event) Times(key string, t []time.Time) *Event {
	e.Event.Times(key, t)
	return e
}

// Dur adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Dur(key string, d time.Duration) *Event {
	e.Event.Dur(key, d)
	return e
}

// Durs adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Durs(key string, d []time.Duration) *Event {
	e.Event.Durs(key, d)
	return e
}

// TimeDiff adds the field key with positive duration between time t and start.
// If time t is not greater than start, duration will be 0.
// Duration format follows the same principle as Dur().
func (e *Event) TimeDiff(key string, t time.Time, start time.Time) *Event {
	e.Event.TimeDiff(key, t, start)
	return e
}

// Interface adds the field key with i marshaled using reflection.
func (e *Event) Interface(key string, i interface{}) *Event {
	e.Event.Interface(key, i)
	return e
}

// CallerSkipFrame instructs any future Caller calls to skip the specified number of frames.
// This includes those added via hooks from the context.
func (e *Event) CallerSkipFrame(skip int) *Event {
	e.Event.CallerSkipFrame(skip)
	return e
}

// Caller adds the file:line of the caller with the zerolog.CallerFieldName key.
// The argument skip is the number of stack frames to ascend
// Skip If not passed, use the global variable CallerSkipFrameCount
func (e *Event) Caller(skip ...int) *Event {
	e.Event = e.Event.Caller(skip...)
	return e
}

// IPAddr adds IPv4 or IPv6 Address to the event
func (e *Event) IPAddr(key string, ip net.IP) *Event {
	e.Event.IPAddr(key, ip)
	return e
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the event
func (e *Event) IPPrefix(key string, pfx net.IPNet) *Event {
	e.Event.IPPrefix(key, pfx)
	return e
}

// MACAddr adds MAC address to the event
func (e *Event) MACAddr(key string, ha net.HardwareAddr) *Event {
	e.Event.MACAddr(key, ha)
	return e
}
