package redirecter

import (
	"github.com/ScruffyProdigy/Middleware/logger"
	"github.com/ScruffyProdigy/TheRack/httper"
	"net/http"
)

type Redirecter struct {
	Path string
}

func (this Redirecter) Run(vars map[string]interface{}, next func()) {
	V(vars).Redirect(this.Path)
}

type V map[string]interface{}

func (vars V) Redirect(path string) {
	info := (logger.V)(vars).Get()
	if info != nil {
		info.Println("Redirecting to " + path)
	}

	r := (httper.V)(vars).GetRequest()
	w := (httper.V)(vars).BlankResponse()
	http.Redirect(w, r, path, http.StatusFound)
	w.Save()
}
