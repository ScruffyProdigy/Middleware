/*
	Encapsulator is a middleware that helps to add a layout to a response
*/
package encapsulator

import (
	"github.com/ScruffyProdigy/Middleware/logger"
	"github.com/ScruffyProdigy/Middleware/templater"
	"github.com/ScruffyProdigy/TheRack/httper"
	"html/template"
)

/*
	Encapsulator is a Middleware that uses TheTemplater to help create Layouts for your websites
	it will take the current body, and place it in the middle of a layou
*/
type Encapsulator struct {
	LayoutVar string //the variable this will look for to find the layout
	BodyVar   string //the variable this will set the old body into that the layout should look for to reapply
	Folder    string //the folder to look for the layouts in
}

// Implementation of rack.Middleware
func (this Encapsulator) Run(vars map[string]interface{}, next func()) {
	next()

	layout, ok := vars[this.LayoutVar].(string)

	if !ok {
		logger.V(vars).Println("Layout not set")
		return
	}

	location := this.Folder + "/" + layout

	if !templater.V(vars).Exists(location) {
		//no "layout", just let it through
		logger.V(vars).Println("Layout \"" + layout + "\" not found")
		return
	}

	vars[this.BodyVar] = template.HTML(httper.V(vars).ResetMessage())

	w := httper.V(vars).FilledResponse()
	templater.V(vars).Render(location, w)
	w.Save()
}

/*
	AddLayout is the default version of Encapsulator

	It will encapsulate the current body, within whichever template is in the "Layout" variable
	The layout will be found in the "layouts" folder, and will use {{.Body}} to specify the old body
*/
var AddLayout = &Encapsulator{LayoutVar: "Layout", BodyVar: "Body", Folder: "layouts"}
