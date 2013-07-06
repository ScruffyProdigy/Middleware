#	Parser
Parser is a simple body parsing middleware that will parse a URL or Form, and store it in the request

## 	Installation
`go get github.com/ScruffyProdigy/Middleware/parser`

##  Documentation
http://godoc.org/github.com/ScruffyProdigy/Middleware/parser

##  Usage

* If you are reading a generic form (ie from a POST or a PUT), just add parser.Form to your rack
	* or alternatively, it will be parsed anyways if you call any of the Form____() vars methods
* If you are loading a file, and know how large to expect it to be, use a parser.Multipart to load it
	* You will need to set it's MaxSize
		* Generally, it should look something like rackup.Add(parser.Multipart{MaxSize:5000})
* If you are loading a file, don't know immediately how large to expect it to be, but can figure it out on the fly, use a parser.VarMultipart to load it
	* In a previous middleware, you will need to set a variable in the vars indicating how large to expect the file to be
	* In the VarMultipart, you will need to set it's MaxBytesVar to the name of the vars variable that contains the expected size of the file
		* Generally, it should look something like rackup.Add(parser.VarMultipart{MaxSizeVar:"File Size"})

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