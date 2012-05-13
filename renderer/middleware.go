package renderer

import (
	"github.com/HairyMezican/TheRack/rack"
	"github.com/HairyMezican/TheTemplater/templater"
)

type Renderer struct {
	Template string
}

func (this Renderer) Run(vars rack.Vars, next func()) {
	Render(vars, this.Template)
}

func Render(vars rack.Vars, s string) {
	w := rack.BlankResponse(vars)
	t, err := templater.Get(s)
	if err != nil {
		panic(err)
	}
	t.Execute(w, vars)
	w.Save()
}
