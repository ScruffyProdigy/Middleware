/*
Error Handler Package recovers from any downstream errors, returning an error status instead
*/
package errorhandler

import (
	"fmt"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
)

func getErrorString(rec interface{}) string {
	err, isError := rec.(error)

	if isError {
		return err.Error()
	}
	return fmt.Sprint(rec)
}

//ErrorHandler is the middleware that you insert into your rack.
// if any downstream Middleware panics, ErrorHandler will catch it and recover
var ErrorHandler rack.Func = func(vars map[string]interface{}, next func()) {
	defer func() {
		rec := recover()
		if rec != nil {
			httper.V(vars).StatusError()
			httper.V(vars).SetMessageString(getErrorString(rec))
		}
	}()
	next()
}
