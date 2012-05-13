#	Statuser
This Middleware looks at the HTTP status code sent down from future Middleware, and looks through the available templates to set an appropriate one.  This is typically used with Encapsulator (github.com/HairyMezican/Middleware/encapsulator) to actually apply the template

## 	Dependencies
This uses TheTemplater to figure out which templates are available (github.com/HairyMezican/TheTemplater)

## 	Installation
`go get github.com/HairyMezican/Middleware/statuser`

## 	Example

__main.go__

	package main

	import (
		"github.com/HairyMezican/Middleware/encapsulator"
		"github.com/HairyMezican/Middleware/statuser"
		"github.com/HairyMezican/TheRack/rack"
		"github.com/HairyMezican/TheRack/templater"
	)

	func main() {
		templater.LoadFromFiles("templates", nil)

		rackup := rack.New()
		rackup.Add(encapsulator.AddLayout)
		rackup.Add(statuser.SetErrorLayout)

		conn := rack.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
__templates/layouts/404.tmpl__

	<html>
		<head>
			<title>404 - Not Found</title>
		</head>
		<body>
			Try Again!
		</body>
	</html>
	
when this is run, and you go to localhost:3000, none of the middleware will resolve the request, and so a 404 error is passed down from the top, but because of the middleware, you should see an appropriate 404 page rendered

Similarly, the following code:

__main.go__

	package main

	import (
		"github.com/HairyMezican/Middleware/encapsulator"
		"github.com/HairyMezican/Middleware/statuser"
		"github.com/HairyMezican/TheRack/rack"
		"github.com/HairyMezican/TheRack/templater"
	)

	var ErrorWare rack.Func = func(vars rack.Vars, next func()) {
		rack.StatusError(vars)
		rack.SetMessageString(vars, "An unknown error has occurred")
	}

	func main() {
		templater.LoadFromFiles("templates", nil)

		rackup := rack.New()
		rackup.Add(encapsulator.AddLayout)
		rackup.Add(statuser.SetErrorLayout)
		rackup.Add(ErrorWare)

		conn := rack.HttpConnection(":3000")
		conn.Go(rackup)
	}
	

__templates/layouts/5xx.tmpl__

	<html>
		<head>
			<title>Error!</title>
		</head>
		<body>
			{{.Error}} - We Messed Up!
		</body>
	</html>

This code will display an appropriate error page for the user