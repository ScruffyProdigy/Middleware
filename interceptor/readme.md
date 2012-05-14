# Interceptor
This is a lightweight routing system that uses a simple lookup on the request URL to route the request to a different middleware.  If none is found, the request passes right through this Middleware.  Please note that "/users" is considered different than "/users/" for the purpose of this middleware.  This is mostly used for APIs

## Installation
`go get github.com/HairyMezican/Middleware/interceptor`

## Example

	package main

	import (
		"github.com/HairyMezican/Middleware/interceptor"
		"github.com/HairyMezican/TheRack/httper"
		"github.com/HairyMezican/TheRack/rack"
	)

	var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
		(httper.V)(vars).SetMessageString("Hello World")
	}

	var RootWare rack.Func = func(vars map[string]interface{}, next func()) {
		(httper.V)(vars).SetMessageString("<html>Check out my <a href=\"/helloworld\">Hello World</a></html>")
	}

	func main() {
		cept := interceptor.New()
		cept.Intercept("/", RootWare)
		cept.Intercept("/helloworld", HelloWorldWare)

		rackup := rack.New()
		rackup.Add(cept)
		rackup.Add(HelloWorldWare)

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
	
running this, http://localhost:3000/ should provide you to a link, and following that link, should give you a hello world