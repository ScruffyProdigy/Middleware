#	Statuser
This Middleware looks at the HTTP status code sent down from future Middleware, and looks through the available templates to set an appropriate one.  This is typically used with Encapsulator (github.com/ScruffyProdigy/Middleware/encapsulator) to actually apply the template

## 	Installation
`go get github.com/ScruffyProdigy/Middleware/statuser`

## Usage

* Generally, you can just add statuser.SetErrorLayout to your rack
	* should add it after both templater and encapsulator in order for it to work properly
	* should add it before any middleware that might set the status on the way back
	* will look inside the "layouts" template folder, set the "Layout" variable to the proper template if found, and set the "Error" variable to the appropriate status code
	* if you want to use a different setup than that, you will need to use a custom setting
* To create a custom setting, you can add a statuser.Statuser to your rack, and fill in the settings manually
	* "Folder" is the template folder it will look into
	* "LayoutVar" is the name of the variable it will set to the best template it finds
	* "ErrorVar" is the variable it will set to the response status code
	* you should add this after the templater middleware in order for it to act properly
	* you should add it before any middleware that might set the status on the way back


## 	Example

__main.go__

	package main

	import (
		"github.com/ScruffyProdigy/Middleware/encapsulator"
		"github.com/ScruffyProdigy/Middleware/statuser"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
		"github.com/ScruffyProdigy/TheTemplater/templater"
	)

	func main() {
		templater.LoadFromFiles("templates", nil)

		rackup := rack.New()
		rackup.Add(encapsulator.AddLayout)
		rackup.Add(statuser.SetErrorLayout)

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
	
__templates/layouts/404.tmpl__

	<html>
		<head>
			<title>404 - Not Found</title>
		</head>
		<body>
			Try Again!
		</body>
	</html>
	
when this is run, and you go to localhost:3000, none of the middleware will resolve the request, and so a 404 error is passed down from the top, but because of the middleware, you should see an appropriate 404 page rendered

Similarly, the following code:

__main.go__

	package main

	import (
		"github.com/ScruffyProdigy/Middleware/encapsulator"
		"github.com/ScruffyProdigy/Middleware/statuser"
		"github.com/ScruffyProdigy/Middleware/templater"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
	)

	var ErrorWare rack.Func = func(vars map[string]interface{}, next func()) {
		v := httper.V(vars)
		v.StatusError()
		v.SetMessageString("An unknown error has occurred")
	}

	func main() {
		rackup := rack.New()
		rackup.Add(templater.GetTemplates("test_templates"))
		rackup.Add(encapsulator.AddLayout)
		rackup.Add(statuser.SetErrorLayout)
		rackup.Add(ErrorWare)

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}
	

__templates/layouts/5xx.tmpl__

	<html>
		<head>
			<title>Error!</title>
		</head>
		<body>
			{{.Error}} - We Messed Up!
		</body>
	</html>

This code will display an appropriate error page for the user