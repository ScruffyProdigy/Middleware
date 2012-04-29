# Error Handler
This is used to trap uncaught errors from future middleware and return a 500 Internal Service Error instead of nothing
Previous Middleware are typically used to filter the error into something more presentable, and/or report the error to your error tracking service

## Dependencies
None

## Installation
`go get github.com/HairyMezican/Middleware/errorhandler`

## Example

	package main

	import (
		"github.com/HairyMezican/Middleware/errorhandler"
		"github.com/HairyMezican/TheRack/rack"
		"net/http"
	)

	var HelloWorldWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		array := make([]byte, 0)
		array[1] = 0 //this action results in a runtime error; we are indexing past the range of the slice
		return http.StatusOK, rack.NewHeader(), array
	}

	func main() {
		conn := rack.HttpConnection(":3000")
		rack.Up.Add(errorhandler.ErrorHandler)
		rack.Up.Add(HelloWorldWare)
		rack.Run(conn, rack.Up)
	}
	

when run, it should simply display "runtime error: index out of range", as there is a runtime error inside of the HelloWorldWare.  The status code returned is a "500 - Internal Service Error"