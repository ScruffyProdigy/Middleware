# Interceptor
This is a lightweight routing system that uses a simple lookup on the request URL to route the request to a different middleware.  If none is found, the request passes right through this Middleware.  Please note that "/users" is considered different than "/users/" for the purpose of this middleware.  This is mostly used for APIs

## Dependencies
None

## Installation
`go get github.com/HairyMezican/Middleware/interceptor`

## Example

	package main

	import (
		"github.com/HairyMezican/Middleware/interceptor"
		"github.com/HairyMezican/TheRack/rack"
		"net/http"
	)

	var HelloWorldWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		return http.StatusOK, rack.NewHeader(), []byte("Hello World!")
	}

	var RootWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		return http.StatusOK, rack.NewHeader(), []byte("<html>Check out my <a href=\"/helloworld\">Hello World</a></html>")
	}

	func main() {
		cept := interceptor.New()
		cept.Intercept("/", RootWare)
		cept.Intercept("/helloworld", HelloWorldWare)

		conn := rack.HttpConnection(":3000")
		rack.Up.Add(cept)
		rack.Up.Add(HelloWorldWare)
		rack.Run(conn, rack.Up)
	}
	
running this, http://localhost:3000/ should provide you to a link, and following that link, should give you a hello world