# Logger
This is used to set a logger that all other middleware will have access to

## Dependencies
None

## Installation
`go get github.com/HairyMezican/Middleware/logger`

## Example

	package main

	import (
		"github.com/HairyMezican/Middleware/logger"
		"github.com/HairyMezican/TheRack/rack"
		"log"
		"net/http"
		"os"
	)

	var HelloWorldWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		lg, isLogger := vars.Apply(logger.Get).(*log.Logger)
		if isLogger {
			lg.Println("Hello World!")
		}
		return http.StatusOK, rack.NewHeader(), []byte("Hello World!")
	}

	func main() {
		conn := rack.HttpConnection(":3000")
		rack.Up.Add(logger.Set(os.Stdout, "Log Test - ", log.LstdFlags))
		rack.Up.Add(HelloWorldWare)
		rack.Run(conn, rack.Up)
	}
	
running this should display a "Hello World!" message at localhost:3000, but when it does so, it should also display a 