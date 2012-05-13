package redirecter

import (
	"github.com/HairyMezican/Middleware/logger"
	"github.com/HairyMezican/TheRack/rack"
	"net/http"
)

type Redirecter struct {
	Apply []rack.VarFunc
	Path  string
}

func (this Redirecter) Run(vars rack.Vars, next func()) {
	Redirect(vars, this.Path, this.Apply...)
}

func Redirect(vars rack.Vars, path string, apply ...rack.VarFunc) {
	info := logger.Get(vars)
	if info != nil {
		info.Println("Redirecting to " + path)
	}

	for _, a := range apply {
		vars.Apply(a)
	}

	r := rack.GetRequest(vars)
	w := rack.BlankResponse(vars)
	http.Redirect(w, r, path, http.StatusFound)
	w.Save()
}
