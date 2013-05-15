package encapsulator

import (
	"../logger"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheTemplater/templater"
	"html/template"
)

/*
	Encapsulator is a Middleware that uses TheTemplater to help create Layouts for your websites
	it will take the current body, and place it in the middle of a layou
*/
type Encapsulator struct {
	Templates *templater.Group
	LayoutVar string //the variable this will look for to find the layout
	BodyVar   string //the variable this will set the old body into that the layout should look for to reapply
	Folder    string //the folder to look for the layouts in
}

func (this Encapsulator) Run(vars map[string]interface{}, next func()) {
	next()

	layout, hasLayout := vars[this.LayoutVar].(string)
	if !hasLayout {
		//no "layout", just let it through
		logger.V(vars).Println("Layout not set")
		return
	}
	
	L := this.Templates.Get(this.Folder + "/" + layout)
	if L == nil {
		//still no layout, just let it through
		logger.V(vars).Println("Layout \""+layout+"\" not found")
		return
	}

	vars[this.BodyVar] = template.HTML(httper.V(vars).ResetMessage())

	w := httper.V(vars).FilledResponse()
	L.Execute(w, vars)
	w.Save()
}

/*
	AddLayout is the default version of Encapsulator

	It will encapsulate the current body, within whichever template is in the "Layout" variable
	The layout will be found in the "layouts" folder, and will use {{.Body}} to specify the old body
*/
func AddLayout(t *templater.Group) *Encapsulator {
	return &Encapsulator{Templates: t, LayoutVar: "Layout", BodyVar: "Body", Folder: "layouts"}
}
