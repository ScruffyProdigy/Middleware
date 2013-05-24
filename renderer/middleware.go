/*
renderer will adjust the response to render a template, either as a middleware or a vars operation
*/
package renderer

import (
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"github.com/ScruffyProdigy/TheTemplater/templater"
)

const (
	template_index = "middleware_template_index"
)

func GetTemplates(loc string) rack.Middleware {
	t := templater.New(loc)
	return rack.Func(func(vars map[string]interface{},next func()) {
		vars[template_index] = t
	})
}

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
	t,ok := vars[template_index].(*templater.Group)
	if !ok {
		(logger.V)(vars).Println("Can't load templates - did you forget to add the 'GetTemplates()' Middleware?")
	}
	
	w := httper.V(vars).BlankResponse()
	
	t, err := t.Get(s)
	if err != nil {
		(logger.V)(vars).Println("Can't load template: "+s)
	}
	t.Execute(w, vars)
	
	w.Save()
}
