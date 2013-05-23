#	Redirecter
This is a simple redirection middleware that will send down a http.StatusFound response with appropriate headers to redirect the user to the url you specified.  It also takes an array of VarFuncs to apply before the redirection, typically a bunch of flash or session setting functions
It also provides a more direct function that another middleware can call directly

## 	Installation
`go get github.com/ScruffyProdigy/Middleware/redirecter`

## 	Example

	package main

	import (
		"github.com/ScruffyProdigy/Middleware/redirecter"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
	)

	var GithubWare rack.Func = func(vars map[string]interface{}, next func()) {
		(redirecter.V)(vars).Redirect("http://github.com/ScruffyProdigy")
	}

	func main() {
		conn := httper.HttpConnection(":3000")
		conn.Go(GithubWare)
	}
	

Running this, and going to localhost:3000 should redirect you to my github page
The following code would do the same thing:

	package main

	import (
		"github.com/ScruffyProdigy/Middleware/redirecter"
		"github.com/ScruffyProdigy/TheRack/httper"
	)

	func main() {
		conn := httper.HttpConnection(":3000")
		conn.Go(redirecter.Redirecter{Path: "http://github.com/ScruffyProdigy"})
	}
	