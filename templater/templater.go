/*
templater helps set up templates for later middleware to render
*/
package templater

import (
	"errors"
	"github.com/ScruffyProdigy/Middleware/logger"
	"github.com/ScruffyProdigy/TheRack/rack"
	"github.com/ScruffyProdigy/TheTemplater/templater"
	"io"
)

const (
	template_index = "middleware_template_index"
)

// A Group is a set of templates
type Group struct {
	*templater.Group
}

// Get Templates will return a Middleware that will load a specified folder of templates into your rack environment
func GetTemplates(loc string) rack.Middleware {
	t, errs := templater.New(loc)

	return rack.Func(func(vars map[string]interface{}, next func()) {
		if len(errs) > 0 {
			for _, err := range errs {
				(logger.V)(vars).Println("Template Loading - Warning - " + err.Error())
			}
		}
		vars[template_index] = Group{t}
		next()
	})
}

// V is a type you can cast your vars to in order to access the following functions
type V map[string]interface{}

// Render() will render one template to an io.Writer
func (vars V) Render(template_name string, out io.Writer) error {
	result, ok := vars[template_index].(Group)
	if !ok {
		return errors.New("Templates not loaded - did you forget to add the \"GetTemplates()\" Middleware?")
	}
	return result.Render(template_name, out, vars)
}

// Exists() checks to see whether or not a template exists
func (vars V) Exists(template_name string) bool {
	templates, ok := vars[template_index].(Group)
	if !ok {
		return false
	}
	template := templates.Get(template_name)
	return template != nil
}
