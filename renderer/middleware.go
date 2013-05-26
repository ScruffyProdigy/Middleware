/*
renderer will adjust the response to render a template, either as a middleware or a vars operation
*/
package renderer

import (
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/Middleware/templater"
	"github.com/ScruffyProdigy/Middleware/logger"
)

/*
Renderer is a middleware that will just render a template onto the response, and return it
*/
type Renderer struct {
	Template string
}

func (this Renderer) Run(vars map[string]interface{}, next func()) {
	V(vars).Render(this.Template)
}

type V map[string]interface{}

/*
Render() is a vars operation that will render a template onto the response
*/
func (vars V) Render(s string) {
	w := httper.V(vars).BlankResponse()
	err := (templater.V)(vars).Render(s,w)
	w.Save()
	
	if err != nil {
		(logger.V)(vars).Println(err.Error())
	}
}
