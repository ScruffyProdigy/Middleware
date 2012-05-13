# 	Encapsulator
This is mostly used to add a layout to the output of future middleware

## 	Dependencies
1.	This uses TheTemplater for access to the templates (github.com/HairyMezican/TheTemplater)
2. 	This tries to use the Logger Middleware to report any problems (github.com/HairyMezican/Middleware/logger)

##	Installation
`go get github.com/HairyMezican/Middleware/encapsulator`

## Example

__main.go__

	package main

	import (
		"github.com/HairyMezican/Middleware/encapsulator"
		"github.com/HairyMezican/TheRack/rack"
		"github.com/TheTemplater/templater"
	)

	var HelloWorldWare rack.Func = func(vars rack.Vars, next func()) {
		vars["Layout"] = "base"
		vars["Title"] = "Hello World"
		rack.AppendMessageString(vars, "Hello World!")
	}

	func main() {
		templater.LoadFromFiles("templates", nil)

		rackup := rack.New()
		rackup.Add(encapsulator.AddLayout)
		rackup.Add(HelloWorldWare)

		conn := rack.HttpConnection(":3000")
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