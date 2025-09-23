package logger

import (
	"io"
	"log"
)

type CustomLogger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

func NewCustomLogger(out io.Writer) *CustomLogger {
	return &CustomLogger{
		info:  log.New(out, "INFO: ", log.Ldate|log.Ltime),
		warn:  log.New(out, "WARN: ", log.Ldate|log.Ltime),
		error: log.New(out, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
