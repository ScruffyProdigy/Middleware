# Logger
This is used to set a logger that all other middleware will have access to

## Installation
`go get github.com/ScruffyProdigy/Middleware/logger`

## Usage

* Generally just add logger.StandardLogger to your rack before anything that will need to log data
* Alternatively call logger.Set() to get a non-standard logger, and then add that to your rack instead

## Example

	package main

	import (
		"github.com/ScruffyProdigy/Middleware/logger"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
		"log"
		"os"
	)

	var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
		(logger.V)(vars).Println("Hello World!")
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