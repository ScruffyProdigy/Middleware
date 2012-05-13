/*
	The interceptor package creates a Middleware that does a lightweight lookup for a bunch of static URLs
*/
package interceptor

import (
	"github.com/HairyMezican/TheRack/rack"
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

func (this Interceptor) Run(vars rack.Vars, next func()) {
	url := rack.GetRequest(vars).URL.Path
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
