package encapsulator

import (
	"github.com/HairyMezican/Middleware/logger"
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheTemplater/templater"
	"html/template"
)

/*
	Encapsulator is a Middleware to be used with templater
	it will encapsulate the current body, within a specified template
*/
type Encapsulator struct {
	LayoutVar string //the variable this will look for to find the layout
	BodyVar   string //the variable this will set the old body into that the layout should look for to reapply
	Folder    string //the folder to look for the layouts in
}

func (this Encapsulator) Run(vars map[string]interface{}, next func()) {
	next()

	layout, hasLayout := vars[this.LayoutVar].(string)
	if !hasLayout {
		//no "layout", just let it through
		return
	}

	vars[this.BodyVar] = template.HTML(httper.V(vars).ResetMessage())
	w := httper.V(vars).FilledResponse()

	L, err := templater.Get(this.Folder + "/" + layout)
	if err != nil {
		//layout not found
		//either log the error and let it through, or panic
		l := logger.V(vars).Get()
		if l != nil {
			l.Println(err.Error())
			return
		} else {
			panic(err)
		}
	}

	L.Execute(w, vars)

	w.Save()
}

/*
	AddLayout is the default version of Encapsulator

	It will encapsulate the current body, within whichever template is in the "Layout" variable
	The layout will be found in the "layouts" folder, and will use {{.Body}} to specify the old body
*/
var AddLayout = Encapsulator{LayoutVar: "Layout", BodyVar: "Body", Folder: "layouts"}
