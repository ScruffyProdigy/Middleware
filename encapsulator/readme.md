# 	Encapsulator
This is mostly used to add a layout to the output of future middleware

## 	Dependencies
1.	This uses TheTemplater for access to the templates (github.com/HairyMezican/TheTemplater)
2. 	This tries to use the Logger Middleware to report any problems (github.com/HairyMezican/Middleware/logger)

##	Installation
`go get github.com/HairyMezican/Middleware/encapsulator`
you'll also need to use a `go get` for each of the dependencies

## Example

__main.go__

	package main

	import (
		"github.com/HairyMezican/Middleware/encapsulator"
		"github.com/HairyMezican/TheRack/rack"
		"github.com/HairyMezican/TheTemplater/templater"
		"net/http"
	)

	var HelloWorldWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		vars["Layout"] = "base"
		vars["Title"] = "Hello World"
		return http.StatusOK, rack.NewHeader(), []byte("Hello World")
	}

	func main() {
		templater.LoadFromFiles("templates", nil)

		conn := rack.HttpConnection(":3000")
		rack.Up.Add(encapsulator.AddLayout)
		rack.Up.Add(HelloWorldWare)
		rack.Run(conn, rack.Up)
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