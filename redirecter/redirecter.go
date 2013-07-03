/*
redirecter makes it simple to redirect the user to another url - either as a middleware, or a vars operation
*/
package redirecter

import (
	"github.com/ScruffyProdigy/Middleware/logger"
	"github.com/ScruffyProdigy/TheRack/httper"
	"net/http"
)

/*
Redirecter is a Middleware struct that will return a response that redirects the user to Path
*/
type Redirecter struct {
	Path string
}

// Run() implements the rack.Middleware interface
func (this Redirecter) Run(vars map[string]interface{}, next func()) {
	V(vars).Redirect(this.Path)
}

// V is a type you can cast your vars to in order to access the following functions
type V map[string]interface{}

/*
Redirect will update the response to be a redirect response
*/
func (vars V) Redirect(path string) {
	(logger.V)(vars).Println("Redirecting to " + path)

	r := (httper.V)(vars).GetRequest()
	w := (httper.V)(vars).BlankResponse()
	http.Redirect(w, r, path, http.StatusFound)
	w.Save()
}
