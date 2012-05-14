# Error Handler
This is used to trap uncaught errors from future middleware and return a 500 Internal Service Error instead of nothing
Previous Middleware are typically used to filter the error into something more presentable, and/or report the error to your error tracking service

## Installation
`go get github.com/HairyMezican/Middleware/errorhandler`

## Example

	package main

	import (
		"github.com/HairyMezican/Middleware/errorhandler"
		"github.com/HairyMezican/TheRack/httper"
		"github.com/HairyMezican/TheRack/rack"
	)

	var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
		array := make([]byte, 0)
		array[1] = 0 //this action results in a runtime error; we are indexing past the range of the slice
	}

	func main() {
		rackup := rack.New()
		rackup.Add(errorhandler.ErrorHandler)
		rackup.Add(HelloWorldWare)

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}
	

when run, it should simply display "runtime error: index out of range", as there is a runtime error inside of the HelloWorldWare.  The status code returned is a "500 - Internal Service Error"