#	Sessioner
This middleware loads a session for any future middleware, and then saves the session before passing control back to previous middleware

## 	Dependencies
This uses the Gorilla Session module to implement sessions (code.google.com/p/gorilla/sessions)

## 	Installation
`go get github.com/ScruffyProdigy/Middleware/sessioner`

## 	Example

	package main

	import (
		"github.com/ScruffyProdigy/Middleware/sessioner"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
		"fmt"
	)

	var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
		s := (sessioner.V)(vars)
		times, ok := s.Get("times").(int)
		if !ok {
			times = 0
		}

		s.Set("times", times+1)

		h := httper.V(vars)
		h.SetMessageString(fmt.Sprint("You have been here ", times, " time"))
		if times != 1 {
			h.AppendMessageString("s before")
		} else {
			h.AppendMessageString(" before")
		}
	}

	func main() {
		rackup := rack.New()
		rackup.Add(sessioner.Middleware)
		rackup.Add(HelloWorldWare)

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
	
going to localhost:3000 will now display the number of times you've been to localhost:3000
Note that changing "Get" to "Clear" will result in the same behavior.  "Clear" will erase the variable after getting it, but it won't matter as it will immediately get overwritten in the following "Set" command