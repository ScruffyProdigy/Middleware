#	Statuser
This Middleware looks at the HTTP status code sent down from future Middleware, and looks through the available templates to set an appropriate one.  This is typically used with Encapsulator (github.com/ScruffyProdigy/Middleware/encapsulator) to actually apply the template

## 	Dependencies
This uses TheTemplater to figure out which templates are available (github.com/ScruffyProdigy/TheTemplater)

## 	Installation
`go get github.com/ScruffyProdigy/Middleware/statuser`

## 	Example

__main.go__

	package main

	import (
		"github.com/ScruffyProdigy/Middleware/encapsulator"
		"github.com/ScruffyProdigy/Middleware/statuser"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
		"github.com/ScruffyProdigy/TheTemplater/templater"
	)

	func main() {
		templater.LoadFromFiles("templates", nil)

		rackup := rack.New()
		rackup.Add(encapsulator.AddLayout)
		rackup.Add(statuser.SetErrorLayout)

		conn := httper.HttpConnection(":3000")
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
		"github.com/ScruffyProdigy/Middleware/encapsulator"
		"github.com/ScruffyProdigy/Middleware/statuser"
		"github.com/ScruffyProdigy/Middleware/templater"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
	)

	var ErrorWare rack.Func = func(vars map[string]interface{}, next func()) {
		v := httper.V(vars)
		v.StatusError()
		v.SetMessageString("An unknown error has occurred")
	}

	func main() {
		rackup := rack.New()
		rackup.Add(templater.GetTemplates("test_templates"))
		rackup.Add(encapsulator.AddLayout)
		rackup.Add(statuser.SetErrorLayout)
		rackup.Add(ErrorWare)

		conn := httper.HttpConnection(":3000")
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