#	Renderer
This is middleware simply renders a template and returns a http.StatusOK
It also provides a more direct function so that another Middleware can render the template instead

## 	Dependencies
This uses templates from TheTemplater to simplify rendering	(github.com/HairyMezican/TheTemplater)

## 	Installation
`go get github.com/HairyMezican/Middleware/renderer`
you'll also need to use a `go get github.com/HairyMezican/TheTemplater` to get the dependency if you haven't already

## 	Example

__main.go__

	package main

	import (
		"github.com/HairyMezican/Middleware/renderer"
		"github.com/HairyMezican/TheRack/rack"
		"github.com/HairyMezican/TheTemplater/templater"
		"net/http"
	)

	var HelloWorldWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		vars["Title"] = "Hello World"
		vars["Message"] = "Hello World"
		return renderer.Render("main",vars)
	}

	func main() {
		templater.LoadFromFiles("templates",nil)

		conn := rack.HttpConnection(":3000")
		rack.Up.Add(HelloWorldWare)
		rack.Run(conn, rack.Up)
	}
	
__templates/main.tmpl__

	<html>
		<head>
			<title>{{.Title}}</title>
		</head>
		<body>
			{{.Message}}
		</body>
	</html>
	
Running this will display an HTML file with the text "Hello World" in both the title and the body.  Changing main.go to the following would do the same thing:

	package main

	import (
		"github.com/HairyMezican/Middleware/renderer"
		"github.com/HairyMezican/TheRack/rack"
		"github.com/HairyMezican/TheTemplater/templater"
		"net/http"
	)

	var HelloWorldWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		vars["Title"] = "Hello World"
		vars["Message"] = "Hello World"
		return next()
	}

	func main() {
		templater.LoadFromFiles("templates", nil)

		conn := rack.HttpConnection(":3000")
		rack.Up.Add(HelloWorldWare)
		rack.Up.Add(renderer.Renderer{Template: "main"})
		rack.Run(conn, rack.Up)
	}
	