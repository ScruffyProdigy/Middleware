/*
Error Handler Package recovers from any downstream errors, and will return a 500 status, and set the body the Error Message
*/
package errorhandler

import (
	"github.com/HairyMezican/TheRack/rack"
)

func getErrorString(rec interface{}) string {
	err, isError := rec.(error)
	str, isString := rec.(string)

	if isError {
		return err.Error()
	} else if isString {
		return str
	}
	return "Unknown Error"
}

var ErrorHandler rack.Func = func(vars rack.Vars, next func()) {
	defer func() {
		rec := recover()
		if rec != nil {
			rack.StatusError(vars)
			rack.SetMessageString(vars, getErrorString(rec))
		}
	}()
	next()
}
