# 	Encapsulator
This is mostly used to add a layout to the output of future middleware

## 	Dependencies
1.	This uses TheTemplater for access to the templates (github.com/ScruffyProdigy/TheTemplater)
2. 	This tries to use the Logger Middleware to report any problems (github.com/ScruffyProdigy/Middleware/logger)

##	Installation
`go get github.com/ScruffyProdigy/Middleware/encapsulator`

## Documentation
http://godoc.org/github.com/ScruffyProdigy/Middleware/encapsulator

## Example

__main.go__

	package main

	import (
		"github.com/ScruffyProdigy/Middleware/encapsulator"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
		"github.com/ScruffyProdigy/TheTemplater/templater"
	)

	var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
		vars["Layout"] = "base"
		vars["Title"] = "Hello World"
		(httper.V)(vars).AppendMessageString("Hello World!")
	}

	func main() {
		templater.LoadFromFiles("templates", nil)

		rackup := rack.New()
		rackup.Add(encapsulator.AddLayout)
		rackup.Add(HelloWorldWare)

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}
	

__templates/layouts/base.templ__

	<html>
		<head>
			<title>{{.Title}}</title>
		</head>
		<body>
			{{.Body}}
		</body>
	</html>
	
when run, opening http://localhost:3000 should give you a complete html document instead of the fragment that the HelloWorldWare by itself would render, and the title should be set the "Hello World"