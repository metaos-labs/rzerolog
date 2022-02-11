package rzerolog

import (
	"path/filepath"
	"strings"
	"time"
)

const (
	year         = "yyyy"
	yearT        = "2006"
	month        = "MM"
	monthT       = "01"
	day          = "dd"
	dayT         = "02"
	hour         = "HH"
	hourT        = "15"
	minute       = "mm"
	minuteT      = "04"
	second       = "ss"
	secondT      = "05"
	microsecond  = "sss"
	microsecondT = "000"
)

// ParseTimeFormat parse time format string.
func ParseTimeFormat(format string) string {
	dir, file := filepath.Split(format)
	file = strings.ReplaceAll(file, year, yearT)
	file = strings.ReplaceAll(file, month, monthT)
	file = strings.ReplaceAll(file, day, dayT)
	file = strings.ReplaceAll(file, hour, hourT)
	file = strings.ReplaceAll(file, minute, minuteT)
	file = strings.ReplaceAll(file, microsecond, microsecondT)
	file = strings.ReplaceAll(file, second, secondT)
	file = time.Now().Format(file)
	return filepath.Join(dir, file)
}
