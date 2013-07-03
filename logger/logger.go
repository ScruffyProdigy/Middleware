/*
Logger helps give a way for middleware to output debugging information or other logged info
*/
package logger

import (
	"io"
	"log"
	"os"
)

// The Logger will set up logging information for later Middleware
type Logger struct {
	log *log.Logger
}

const (
	loggerIndex = "Logger"
)

// Run() implements the rack.Middleware interface
func (this Logger) Run(vars map[string]interface{}, next func()) {
	vars[loggerIndex] = this.log
	next()
}

// New() will create a new Logger with the information you specify
func New(out io.Writer, prefix string, flag int) Logger {
	return Logger{log.New(out, prefix, flag)}
}

// V is a type you can cast your vars to in order to access the following functions
type V map[string]interface{}

// Get() will get you direct access to a log.Logger.
// in general, you will usually be able to get by with some of the more direct functions
func (vars V) Get() *log.Logger {
	result, ok := vars["Logger"].(*log.Logger)
	if !ok {
		return nil
	}
	return result
}

// Print() is the equivalent of the log.Print()
func (vars V) Print(v ...interface{}) {
	l := vars.Get()
	if l != nil {
		l.Print(v...)
	}
}

// Printf() is the equivalent of log.Printf()
func (vars V) Printf(format string, v ...interface{}) {
	l := vars.Get()
	if l != nil {
		l.Printf(format, v...)
	}
}

// Println() is the equivalent of log.Println()
func (vars V) Println(v ...interface{}) {
	l := vars.Get()
	if l != nil {
		l.Println(v...)
	}
}

// StandardLogger is a Logger with the default variables set that you can add into your rack without worrying about the settings
var StandardLogger = New(os.Stdout, "", log.LstdFlags)
