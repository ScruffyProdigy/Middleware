#	Methoder
Allows basic html forms to issue put and delete requests by including "_method" as a hidden input with the requested method

##  Dependencies
The form values need to be parsed before this is ran, so a middleware such as parser would be useful (github.com/HairyMezican/Middleware/parser)

## 	Installation
`go get github.com/HairyMezican/Middleware/methoder`

## 	Example

    package main

	import (
		"github.com/HairyMezican/Middleware/methoder"
		"github.com/HairyMezican/Middleware/parser"
		"github.com/HairyMezican/TheRack/rack"
	)

	var HttpWare rack.Func = func(vars rack.Vars, next func()) {
		rack.SetMessageString(vars, "<html><head><title>Form!</title></head><body>")
		request := rack.GetRequest(vars)
		rack.AppendMessageString(vars, "<p>You used "+request.Method+"</p>")
		rack.AppendMessageString(vars, "<form action='/' method='post'><input type='hidden' name='_method' value='put' /><input type='submit' value='put'/></form>")
		rack.AppendMessageString(vars, "<form action='/' method='post'><input type='submit' value='post'/></form>")
		rack.AppendMessageString(vars, "<form action='/' method='post'><input type='hidden' name='_method' value='delete' /><input type='submit' value='delete'/></form>")
	}

	func main() {
		rackup := rack.New()
		rackup.Add(parser.Form)
		rackup.Add(methoder.Override)
		rackup.Add(HttpWare)

		conn := rack.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
	
Going to localhost:3000 should tell you what method was used to access the page, along with 3 buttons that allow you to access the same page with a different method