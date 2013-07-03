/*
	The interceptor package creates a Middleware that does a lightweight lookup for a bunch of static URLs
*/
package interceptor

import (
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
)

//PreExistingInterceptorError is used when you try to intercept a URL that is already being intercepted
type PreExistingInterceptorError struct {
	url string
}

// implements the error interface
func (this PreExistingInterceptorError) Error() string {
	return "Interception '" + this.url + "' already registered!"
}

// The Interceptor is a Middleware that handles requests by matching them with a list of URLs and passing them off to other corresponding Middleware
type Interceptor map[string]rack.Middleware

// Intercept tells the interceptor to trap requests matching a specific URL and to pass them to a corresponding Middleware
func (this Interceptor) Intercept(url string, exec rack.Middleware) error {
	if this[url] != nil {
		return PreExistingInterceptorError{url}
	}
	this[url] = exec
	return nil
}

// Run implements the rack.Middleware interface
func (this Interceptor) Run(vars map[string]interface{}, next func()) {
	url := httper.V(vars).GetRequest().URL.Path
	exec := this[url]
	if exec != nil {
		exec.Run(vars, next)
	} else {
		next()
	}
}

// New gives you a blank interceptor
func New() Interceptor {
	return make(Interceptor)
}
