/*
	The interceptor package creates a Middleware that does a lightweight lookup for a bunch of static URLs
*/
package interceptor

import (
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
)

type PreExistingInterceptorError struct {
	url string
}

func (this PreExistingInterceptorError) Error() string {
	return "Interception '" + this.url + "' already registered!"
}

type Interceptor map[string]rack.Middleware

func (this Interceptor) Intercept(url string, exec rack.Middleware) {
	if this[url] != nil {
		panic(PreExistingInterceptorError{url})
	}
	this[url] = exec
}

func (this Interceptor) Run(vars map[string]interface{}, next func()) {
	url := httper.V(vars).GetRequest().URL.Path
	exec := this[url]
	if exec != nil {
		exec.Run(vars, next)
	} else {
		next()
	}
}

func New() Interceptor {
	return make(Interceptor)
}
