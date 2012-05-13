package logger

import (
	"github.com/HairyMezican/TheRack/rack"
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

func (this Logger) Run(vars rack.Vars, next func()) {
	vars[loggerIndex] = this.log
	next()
}

func Set(out io.Writer, prefix string, flag int) Logger {
	return Logger{log.New(out, prefix, flag)}
}

func Get(vars rack.Vars) *log.Logger {
	result, ok := vars["Logger"].(*log.Logger)
	if !ok {
		return nil
	}
	return result
}

var StandardLogger = Set(os.Stdout, "", log.LstdFlags)
