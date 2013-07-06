# Error Handler
This is used to trap uncaught errors from future middleware and return a 500 Internal Service Error instead of nothing
Previous Middleware are typically used to filter the error into something more presentable, and/or report the error to your error tracking service

## Installation
`go get github.com/ScruffyProdigy/Middleware/errorhandler`

## Docuemnation
http://godoc.org/github.com/ScruffyProdigy/Middleware/errorhandler

## Usage

* Just add errorhandler.ErrorHandler to your rack

## Example

	package main

	import (
		"github.com/ScruffyProdigy/Middleware/errorhandler"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
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