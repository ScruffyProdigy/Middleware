#	Redirecter
This is a simple redirection middleware that will send down a http.StatusFound response with appropriate headers to redirect the user to the url you specified.  It also takes an array of VarFuncs to apply before the redirection, typically a bunch of flash or session setting functions
It also provides a more direct function that another middleware can call directly

## 	Dependencies
None

## 	Installation
`go get github.com/HairyMezican/Middleware/redirecter`

## 	Example

	package main

	import (
		"github.com/HairyMezican/Middleware/redirecter"
		"github.com/HairyMezican/TheRack/rack"
		"net/http"
	)

	var GithubWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		return redirecter.Redirect(r, vars, "http://github.com/HairyMezican/")
	}

	func main() {
		conn := rack.HttpConnection(":3000")
		rack.Up.Add(GithubWare)
		rack.Run(conn, rack.Up)
	}

Running this, and going to localhost:3000 should redirect you to my github page
The following code would do the same thing:

	package main

	import (
		"../Middleware/redirecter"
		"github.com/HairyMezican/TheRack/rack"
	)

	func main() {
		conn := rack.HttpConnection(":3000")
		rack.Run(conn, redirecter.Redirecter{Path:"http://github.com/HairyMezican"})
	}
	