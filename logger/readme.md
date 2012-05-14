# Logger
This is used to set a logger that all other middleware will have access to

## Installation
`go get github.com/HairyMezican/Middleware/logger`

## Example

	package main

	import (
		"github.com/HairyMezican/Middleware/logger"
		"github.com/HairyMezican/TheRack/httper"
		"github.com/HairyMezican/TheRack/rack"
		"log"
		"os"
	)

	var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
		lg := (logger.V)(vars).Get()
		if lg != nil {
			lg.Println("Hello World!")
		}
		(httper.V)(vars).SetMessageString("Hello World!")
	}

	func main() {
		rackup := rack.New()
		rackup.Add(logger.Set(os.Stdout, "Log Test - ", log.LstdFlags))
		rackup.Add(HelloWorldWare)

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
	
running this should display a "Hello World!" message at localhost:3000, but when it does so, it should also display the same message to the standard output with "Log Test - " and the time appended