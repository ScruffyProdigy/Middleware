#	Methoder
Allows basic html forms to issue put and delete requests by including "_method" as a hidden input with the requested method

##  Dependencies
The form values need to be parsed before this is ran, so a middleware such as parser would be useful (github.com/ScruffyProdigy/Middleware/parser)

## 	Installation
`go get github.com/ScruffyProdigy/Middleware/methoder`

## Usage

* Add methoder.Override to your rack
	* Somewhere after the form has been parsed
	* Somewhere before any middleware that inspects the method of the request

## 	Example

    package main

	import (
		"github.com/ScruffyProdigy/Middleware/methoder"
		"github.com/ScruffyProdigy/Middleware/parser"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
	)

	var HttpWare rack.Func = func(vars map[string]interface{}, next func()) {
		v := httper.V(vars)

		v.SetMessageString("<html><head><title>Form!</title></head><body>")
		request := v.GetRequest()
		v.AppendMessageString("<p>You used " + request.Method + "</p>")
		v.AppendMessageString("<form action='/' method='post'><input type='hidden' name='_method' value='put' /><input type='submit' value='put'/></form>")
		v.AppendMessageString("<form action='/' method='post'><input type='submit' value='post'/></form>")
		v.AppendMessageString("<form action='/' method='post'><input type='hidden' name='_method' value='delete' /><input type='submit' value='delete'/></form>")
	}

	func main() {
		rackup := rack.New()
		rackup.Add(parser.Form)
		rackup.Add(methoder.Override)
		rackup.Add(HttpWare)

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
	
Going to localhost:3000 should tell you what method was used to access the page, along with 3 buttons that allow you to access the same page with a different method