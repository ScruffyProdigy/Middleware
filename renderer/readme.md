#	Renderer
This is middleware simply renders a template and returns a http.StatusOK
It also provides a more direct function so that another Middleware can render the template instead

## 	Dependencies
This uses templates from TheTemplater to simplify rendering	(github.com/ScruffyProdigy/TheTemplater)

## 	Installation
`go get github.com/ScruffyProdigy/Middleware/renderer`

## 	Example

__main.go__

	package main

	import (
		"github.com/ScruffyProdigy/Middleware/renderer"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
		"github.com/ScruffyProdigy/TheTemplater/templater"
	)

	var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
		vars["Title"] = "Hello World"
		vars["Message"] = "Hello World"
		(renderer.V)(vars).Render("main")
	}

	func main() {
		templater.LoadFromFiles("templates", nil)

		conn := httper.HttpConnection(":3000")
		conn.Go(HelloWorldWare)
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
		"github.com/ScruffyProdigy/Middleware/renderer"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
		"github.com/ScruffyProdigy/TheTemplater/templater"
	)

	var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
		vars["Title"] = "Hello World"
		vars["Message"] = "Hello World"
		next()
	}

	func main() {
		templater.LoadFromFiles("templates", nil)

		rackup := rack.New()
		rackup.Add(HelloWorldWare)
		rackup.Add(renderer.Renderer{Template: "main"})

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}
	