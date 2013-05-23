#	Parser
Parser is a simple body parsing middleware that will parse a URL or Form, and store it in the request

## 	Installation
`go get github.com/ScruffyProdigy/Middleware/parser`

## 	Example

    package main

	import (
		"github.com/ScruffyProdigy/Middleware/parser"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
	)

	var HttpWare rack.Func = func(vars map[string]interface{}, next func()) {
		h := httper.V(vars)
		p := parser.V(vars)

		h.SetMessageString("<html><head><title>Form!</title></head><body>")
		name := p.FormValue("Name")
		if name != "" {
			h.AppendMessageString("<p>Welcome " + name + "</p>")
		}
		h.AppendMessageString("<form action='/' method='post'><label for='name'>Name:</label><input id='name' name='Name' type='text' value='' /><input type='submit' /></form>")
	}

	func main() {
		rackup := rack.New()
		rackup.Add(parser.Form)
		rackup.Add(HttpWare)

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
	
Going to localhost:3000 should present you with a form, filling out the form and submitting it should take you back to the same page, except it greets you with the name you entered