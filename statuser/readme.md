#	Statuser
This Middleware looks at the HTTP status code sent down from future Middleware, and looks through the available templates to set an appropriate one.  This is typically used with Encapsulator (github.com/HairyMezican/Middleware/encapsulator) to actually apply the template

## 	Dependencies
This uses TheTemplater to figure out which templates are available (github.com/HairyMezican/TheTemplater)

## 	Installation
`go get github.com/HairyMezican/Middleware/statuser`
you'll also need to use a `go get github.com/HairyMezican/TheTemplater` to get the dependency if you haven't already

## 	Example

__main.go__

	package main

	import (
		"github.com/HairyMezican/Middleware/encapsulator"
		"github.com/HairyMezican/Middleware/statuser"
		"github.com/HairyMezican/TheRack/rack"
		"github.com/HairyMezican/TheTemplater/templater"
	)

	func main() {
		templater.LoadFromFiles("templates", nil)

		conn := rack.HttpConnection(":3000")
		rack.Up.Add(encapsulator.AddLayout)
		rack.Up.Add(statuser.SetErrorLayout)
		rack.Run(conn, rack.Up)
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
		"github.com/HairyMezican/TheTemplater/templater"
		"net/http"
	)

	var ErrorWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		return http.StatusInternalServerError, rack.NewHeader(), []byte("")
	}

	func main() {
		templater.LoadFromFiles("templates", nil)

		conn := rack.HttpConnection(":3000")
		rack.Up.Add(encapsulator.AddLayout)
		rack.Up.Add(statuser.SetErrorLayout)
		rack.Up.Add(ErrorWare)
		rack.Run(conn, rack.Up)
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