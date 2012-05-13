#	Sessioner
This middleware loads a session for any future middleware, and then saves the session before passing control back to previous middleware

## 	Dependencies
This uses the Gorilla Session module to implement sessions (code.google.com/p/gorilla/sessions)

## 	Installation
`go get github.com/HairyMezican/Middleware/sessioner`

## 	Example

	package main

	import (
		"github.com/HairyMezican/Middleware/sessioner"
		"github.com/HairyMezican/TheRack/rack"
		"fmt"
	)

	var HelloWorldWare rack.Func = func(vars rack.Vars, next func()) {
		times, ok := sessioner.Get(vars, "times").(int)
		if !ok {
			times = 0
		}

		sessioner.Set(vars, "times", times+1)

		rack.SetMessageString(vars, fmt.Sprint("You have been here ", times, " time"))
		if times != 1 {
			rack.AppendMessageString(vars, "s before")
		} else {
			rack.AppendMessageString(vars, " before")
		}
	}

	func main() {
		rackup := rack.New()
		rackup.Add(sessioner.Middleware)
		rackup.Add(HelloWorldWare)

		conn := rack.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
	
going to localhost:3000 will now display the number of times you've been to localhost:3000
Note that changing "Get" to "Clear" will result in the same behavior.  "Clear" will erase the variable after getting it, but it won't matter as it will immediately get overwritten in the following "Set" command