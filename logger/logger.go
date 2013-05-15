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

type V map[string]interface{}

func (vars V) Get() *log.Logger {
	result, ok := vars["Logger"].(*log.Logger)
	if !ok {
		return nil
	}
	return result
}

func (vars V) Print(v ...interface{}) {
	l := vars.Get()
	if l != nil {
		l.Print(v...)
	}
}

func (vars V) Printf(format string, v ...interface{}) {
	l := vars.Get()
	if l != nil {
		l.Printf(format, v...)
	}
}

func (vars V) Println(v ...interface{}) {
	l := vars.Get()
	if l != nil {
		l.Println(v...)
	}
}

var StandardLogger = Set(os.Stdout, "", log.LstdFlags)
