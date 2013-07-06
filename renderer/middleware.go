/*
	renderer renders a template for you, either as a middleware or a vars method
*/
package renderer

import (
	"github.com/ScruffyProdigy/Middleware/logger"
	"github.com/ScruffyProdigy/Middleware/templater"
	"github.com/ScruffyProdigy/TheRack/httper"
)

/*
Renderer is a middleware that will just render a template onto the response, and return it
*/
type Renderer struct {
	Template string
}

//Run implements the rack.Middleware interface
func (this Renderer) Run(vars map[string]interface{}, next func()) {
	V(vars).Render(this.Template)
}

// V is a type you can cast your vars to in order to access the following functions
type V map[string]interface{}

/*
Render() is a vars operation that will render a template onto the response
*/
func (vars V) Render(s string) {
	w := httper.V(vars).BlankResponse()
	err := (templater.V)(vars).Render(s, w)
	w.Save()

	if err != nil {
		(logger.V)(vars).Println(err.Error())
	}
}
