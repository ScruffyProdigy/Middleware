#	Sessioner
This middleware loads a session for any future middleware, and then saves the session before passing control back to previous middleware

## 	Dependencies
This uses the Gorilla Session module to implement sessions (code.google.com/p/gorilla/sessions)

## 	Installation
`go get github.com/HairyMezican/Middleware/sessioner`
you'll also need to use a `go get code.google.com/p/gorilla/sessions` to get the dependency if you haven't already

## 	Example

	package main

	import (
		"fmt"
		"github.com/HairyMezican/Middleware/sessioner"
		"github.com/HairyMezican/TheRack/rack"
		"net/http"
	)

	var HelloWorldWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		times, _ := vars.Apply(sessioner.Get("times")).(int)
		vars.Apply(sessioner.Set("times", times+1))
		message := fmt.Sprint("You have been here ", times, " time")
		if times != 1 {
			message += "s"
		}
		message += " before"

		return http.StatusOK, rack.NewHeader(), []byte(message)
	}

	func main() {
		conn := rack.HttpConnection(":3000")
		rack.Up.Add(sessioner.Middleware)
		rack.Up.Add(HelloWorldWare)
		rack.Run(conn, rack.Up)
	}
	
going to localhost:3000 will now display the number of times you've been to localhost:3000
Note that changing "Get" to "Clear" will result in the same behavior.  "Clear" will erase the variable after getting it, but it won't matter as it will immediately get overwritten in the following "Set" command