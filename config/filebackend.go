package config

import (
	logging "github.com/op/go-logging"
	"os"
)

type FileBackend struct {
	File *os.File
}

func NewFileBackend(file *os.File) *FileBackend {
	return &FileBackend{File: file}
}

func (filebackend *FileBackend) Log(level logging.Level, calldepth int, record *logging.Record) (err error) {
	line := record.Formatted(calldepth + 1)

	_, err = filebackend.File.WriteString(line + "\n")
	if err != nil {
		return err
	}

	err = filebackend.File.Sync()
	if err != nil {
		return err
	}

	return nil
}
