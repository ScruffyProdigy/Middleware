package redirecter

import (
	"log"
	"github.com/HairyMezican/TheRack/rack"
	"net/http"
)

type Redirecter struct {
	Apply []rack.VarFunc
	Path  string
}


func (this Redirecter) Run(r *http.Request, vars rack.Vars, next rack.Next) (status int, header http.Header, message []byte) {
	return Redirect(r, vars, this.Path, this.Apply...)
}

func Redirect(r *http.Request, vars rack.Vars, path string, apply ...rack.VarFunc) (status int, header http.Header, message []byte) {
	info,isLogger := vars["Logger"].(log.Logger)
	if isLogger {
		info.Println("Redirecting to "+path)
	}

	for _, a := range apply {
		vars.Apply(a)
	}

	w := rack.BlankResponse()
	http.Redirect(w, r, path, http.StatusFound)
	return w.Results()
}
