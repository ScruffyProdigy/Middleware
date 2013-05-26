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
		"github.com/ScruffyProdigy/Middleware/templater"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
	)

	var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
		vars["Title"] = "Hello World"
		vars["Message"] = "Hello World"
		(renderer.V)(vars).Render("main")
	}

	func main() {
		rackup := rack.New()
		rackup.Add(templater.GetTemplates("./templates"))
		rackup.Add(HelloWorldWare)
		
		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
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
		"github.com/ScruffyProdigy/Middleware/templater"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
	)

	var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
		vars["Title"] = "Hello World"
		vars["Message"] = "Hello World"
		next()
	}

	func main() {
		rackup := rack.New()
		rackup.Add(templater.GetTemplates("./templates"))
		rackup.Add(HelloWorldWare)
		rackup.Add(renderer.Renderer{Template: "main"})

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}
	