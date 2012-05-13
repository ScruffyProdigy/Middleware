#Middleware

These are middleware that are each designed to be part of a Rack specified in https://github.com/HairyMezican/TheRack
Look at the documentation in each folder for more information about how they are used

*	**encapsulator** - Typically used to wrap the current body in a layout  
*	**errorhandler** - Used to catch any panics from any later middleware, and, if so, passes a 500 Internal Service Error to the previous ones  
*	**interceptor** - A lightweight router that uses a map[string]Middleware lookup on the requested URL to pass the request on to other middleware  
*	**logger** - sets a logger in the vars for other middleware to report potential problems  
*	**methoder** - allows html forms to make put or delete requests
*	**oauther** - provides an interface for oauth providers; sets up appropriate routes (requires interceptor)  
*	**parser** - parses body or url forms
*	**redirecter** - returns an appropriate redirect response  
*	**renderer** - renders an appropriate template (requires TheTemplater)  
*	**routes** - a branching based router  
*	**session** - provides a middleware wrapper for Gorilla based Sessions  
*	**statuser** - sets appropriate template variables based on http statuses sent from later middleware - useful for setting error layouts
*	**websocketer** - provides an rack-based interface for websockets

##Installation
You can install them all by simply running `go get github.com/HairyMezican/Middleware/...`