package statuser

import (
	"github.com/HairyMezican/TheRack/rack"
	"github.com/HairyMezican/TheTemplater/templater"
	"net/http"
	"strconv"
)

/*
	Statuser is a Middleware tool that is to be used with Encapsulator and Templater

	It checks the status of the response, and sets the layout to status number if it is found
	if it is not found, it sets it to a more general layout (if found) such as 40x, or 50x
*/
type Statuser struct {
	ErrorVar  string //the variable to store the error code in
	LayoutVar string //the variable to store the layout in
	Folder    string //the folder where the layouts are kept
}

func (this Statuser) Run(r *http.Request, vars rack.Vars, next rack.Next) (status int, header http.Header, message []byte) {
	status, header, message = next()

	layout := strconv.Itoa(status)
	if templater.Available(this.Folder + "/" + layout) {
		vars[this.LayoutVar] = layout
		return
	}

	layout = strconv.Itoa(status/100) + "xx"
	if templater.Available(this.Folder + "/" + layout) {
		vars[this.ErrorVar] = strconv.Itoa(status)
		vars[this.LayoutVar] = layout
		return
	}

	return
}

/*
	SetErrorLayout is the default version of Statuser, that will set the "Layout" var to an appropriate template from the layouts folder
	and will set the "Error" var if we are using a generic template
*/
var SetErrorLayout = Statuser{LayoutVar: "Layout", Folder: "layouts", ErrorVar: "Error"}
