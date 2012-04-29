/*
	The interceptor package creates a Middleware that does a lightweight lookup for a bunch of static URLs
*/
package interceptor

import (
	"github.com/HairyMezican/TheRack/rack"
	"net/http"
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

func (this Interceptor) Run(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
	url := r.URL.Path
	exec := this[url]
	if exec == nil {
		return next()
	}
	return exec.Run(r, vars, next)
}

func New() Interceptor {
	return make(Interceptor)
}
