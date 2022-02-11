package rzerolog

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
)

var (
	_ io.Writer  = (*LogFileWriter)(nil)
	_ FileWriter = (*osFileWriter)(nil)
)

type osFileWriter struct {
	*os.File
}

func (o *osFileWriter) SetOutput(file *os.File) error {
	o.File = file
	return nil
}

type LogFileWriter struct {
	enable bool
	writer FileWriter

	fillPath string
	// time rolling params
	timeRolling bool
	// size rolling params
	sizeRolling  bool
	fileSize     int64
	maxFileCount int
	// file params
	logFileName     string
	currentFileName string
	file            *os.File

	lockC chan struct{}
}

func (f *LogFileWriter) initBase() error {
	if !f.enable {
		return nil
	}
	if f.lockC == nil {
		f.lockC = make(chan struct{}, 1)
	}
	if f.file == nil {
		// check log file path
		err := os.MkdirAll(f.fillPath, 0777)
		if err != nil {
			return err
		}
		newFileName := f.logFileName
		if f.timeRolling {
			newFileName = ParseTimeFormat(newFileName)
		}

		// create new file
		newFile, err := os.OpenFile(filepath.Join(f.fillPath, newFileName),
			os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return err
		}

		f.file = newFile
		if err = f.writer.SetOutput(f.file); err != nil {
			return err
		}
		f.currentFileName = newFileName
	}
	return nil
}

func (f *LogFileWriter) Write(p []byte) (n int, err error) {
	if !f.enable {
		return len(p), nil
	}
	f.lockC <- struct{}{}
	if err = f.doTimeRolling(); err != nil {
		return 0, err
	}
	n, err = f.writer.Write(p)
	if err != nil {
		return 0, err
	}
	if err = f.doSizeRolling(len(p)); err != nil {
		return 0, err
	}
	<-f.lockC
	return
}

func (f *LogFileWriter) doTimeRolling() error {
	if !f.timeRolling {
		return nil
	}

	newFileName := ParseTimeFormat(f.logFileName)
	if newFileName != f.currentFileName {
		// create new file
		newFile, err := os.OpenFile(filepath.Join(f.fillPath, newFileName),
			os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_SYNC, 0666)
		if err != nil {
			return err
		}
		oldFile := f.file
		f.file = newFile
		if err = f.writer.SetOutput(f.file); err != nil {
			return err
		}
		f.currentFileName = newFileName
		defer oldFile.Close()
	}
	return nil
}

func (f *LogFileWriter) currFileSize() (int64, error) {
	fInfo, err := f.file.Stat()
	if err != nil {
		return 0, err
	}
	return fInfo.Size(), nil
}

func (f *LogFileWriter) doSizeRolling(append int) error {
	if !f.sizeRolling || f.fileSize == 0 {
		return nil
	}

	// check file size
	currSize, err := f.currFileSize()
	if err != nil {
		return err
	}
	if currSize+int64(append) <= f.fileSize {
		return nil
	}
	// rename full file
	err = f.renameCurrentFile()
	if err != nil {
		return err
	}

	// create new file
	newFile, err := os.OpenFile(filepath.Join(f.fillPath, f.currentFileName),
		os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		return err
	}
	f.file = newFile
	if err = f.writer.SetOutput(f.file); err != nil {
		return err
	}
	fileLoc := filepath.Join(f.fillPath, f.currentFileName)
	go f.renameOldFiles(fileLoc)
	return nil
}

func (f *LogFileWriter) renameCurrentFile() error {
	if err := f.file.Sync(); err != nil {
		return err
	}
	if err := f.file.Close(); err != nil {
		return err
	}
	fileLoc := filepath.Join(f.fillPath, f.currentFileName)
	return os.Rename(fileLoc, fileLoc+".0")
}

func (f *LogFileWriter) renameOldFiles(fileLoc string) {
	for i := f.maxFileCount; i > 0; i-- {
		curr := fileLoc + "." + strconv.Itoa(i-1)
		now := fileLoc + "." + strconv.Itoa(i)
		if i == f.maxFileCount {
			_ = os.Remove(now)
			continue
		}
		_ = os.Rename(curr, now)
	}
}
