package renderer

import (
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheTemplater/templater"
)

type Renderer struct {
	Template string
}

func (this Renderer) Run(vars map[string]interface{}, next func()) {
	V(vars).Render(this.Template)
}

type V map[string] interface{}

func (vars V) Render(s string) {
	w := httper.V(vars).BlankResponse()
	t, err := templater.Get(s)
	if err != nil {
		panic(err)
	}
	t.Execute(w, vars)
	w.Save()
}
