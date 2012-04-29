package renderer

import (
	"github.com/HairyMezican/TheRack/rack"
	"github.com/HairyMezican/TheTemplater/templater"
	"net/http"
)

type Renderer struct {
	Template string
}

func (this Renderer) Run(r *http.Request, vars rack.Vars, next rack.Next) (status int, header http.Header, message []byte) {
	return Render(this.Template, vars)
}

func Render(s string, vars rack.Vars) (status int, header http.Header, message []byte) {
	w := rack.BlankResponse()
	t, err := templater.Get(s)
	if err != nil {
		panic(err)
	}
	t.Execute(w, vars)
	return w.Results()
}
