#	Parser
Parser is a simple body parsing middleware that will parse a URL or Form, and store it in the request

## 	Installation
`go get github.com/HairyMezican/Middleware/parser`

## 	Example

    package main

	import (
		"github.com/HairyMezican/Middleware/parser"
		"github.com/HairyMezican/TheRack/rack"
	)

	var HttpWare rack.Func = func(vars rack.Vars, next func()) {
		rack.SetMessageString(vars, "<html><head><title>Form!</title></head><body>")
		name := parser.FormValue(vars, "Name")
		if name != "" {
			rack.AppendMessageString(vars, "<p>Welcome "+name+"</p>")
		}
		rack.AppendMessageString(vars, "<form action='/' method='post'><label for='name'>Name:</label><input id='name' name='Name' type='text' value='' /><input type='submit' /></form>")
	}

	func main() {
		rackup := rack.New()
		rackup.Add(parser.Form)
		rackup.Add(HttpWare)

		conn := rack.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
Going to localhost:3000 should present you with a form, filling out the form and submitting it should take you back to the same page, except it greets you with the name you entered