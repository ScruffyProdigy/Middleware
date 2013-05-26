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

type Group struct {
	*templater.Group
}

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

type V map[string]interface{}

func (vars V) Render(template_name string, out io.Writer) error {
	result, ok := vars[template_index].(Group)
	if !ok {
		return errors.New("Templates not loaded - did you forget to add the \"GetTemplates()\" Middleware?")
	}
	return result.Render(template_name, out, vars)
}

func (vars V) Exists(template_name string) bool {
	templates, ok := vars[template_index].(Group)
	if !ok {
		return false
	}
	template := templates.Get(template_name)
	return template != nil
}
