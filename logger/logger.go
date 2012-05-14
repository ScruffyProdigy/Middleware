package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	log *log.Logger
}

const (
	loggerIndex = "Logger"
)

func (this Logger) Run(vars map[string]interface{}, next func()) {
	vars[loggerIndex] = this.log
	next()
}

func Set(out io.Writer, prefix string, flag int) Logger {
	return Logger{log.New(out, prefix, flag)}
}

type V map[string] interface{}

func (vars V) Get() *log.Logger {
	result, ok := vars["Logger"].(*log.Logger)
	if !ok {
		return nil
	}
	return result
}

var StandardLogger = Set(os.Stdout, "", log.LstdFlags)
