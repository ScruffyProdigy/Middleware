/*
methoder creates a way for you to allow HTML forms to submit PUT and DELETE requests, which are useful to allow for a complete restful vocabulary
*/
package methoder

import (
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"strings"
)

var legal = map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}

func isLegal(s string) bool {
	return legal[s]
}

/*
Override is a Middleware that will override the method used in the request if a _method argument is passed in
*/
var Override rack.Func = func(vars map[string]interface{}, next func()) {
	r := httper.V(vars).GetRequest()
	method := strings.ToUpper(r.Form.Get("_method"))
	if isLegal(method) {
		r.Method = method
	}
	next()
}
