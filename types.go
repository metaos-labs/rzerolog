package rzerolog

import (
	"io"
	"os"
)

// FileWriter is a writer for writing log into files.
type FileWriter interface {
	io.Writer
	// SetOutput change the file which be written in.
	SetOutput(file *os.File) error
}
