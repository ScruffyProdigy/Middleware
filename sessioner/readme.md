#	Sessioner
This middleware loads a session for any future middleware, and then saves the session before passing control back to previous middleware

## 	Installation
`go get github.com/ScruffyProdigy/Middleware/sessioner`

##  Usage

* Add sessioner.Middleware to your rack
	* Must be added before any middleware that takes advantage of any of the following functions:
* To save to the session, call (sessioner.V)(vars).Set()
	* The first parameter is the name of the variable that you are setting
	* The second parameter is the value you are setting it to
* To access the session, call (sessioner.V)(vars).Get()
	* The only parameter is the name of the variable that you are getting
	* returns the value of that parameter, or nil if it doesn't exist
* To clear a session variable call (sessioner.V)(vars).Clear()
	* The only parameter is the name of the variable that you are clearing
	* returns the old value of that parameter, or nil if it never existed
* To add "flashes", call (sessioner.V)(vars).AddFlash()
	* The only parameter is a string, which will get appended to the flashes currently available
	* flashes are only available upon a redirect or reload, and are then destroyed
* To retrieve "flashes", call (sessioner.V)(vars).Flashes()
	* Will return a slice of flashes
	* flashes are only available upon a redirect or reload, and are then destroyed

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