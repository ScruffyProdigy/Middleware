package parser

import (
	"github.com/HairyMezican/TheRack/rack"
)

var Form rack.Func = func(vars rack.Vars, next func()) {
	r := rack.GetRequest(vars)
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	next()
}

type Multipart struct {
	MaxSize int64
}

func (this Multipart) Run(vars rack.Vars, next func()) {
	r := rack.GetRequest(vars)
	err := r.ParseMultipartForm(this.MaxSize)
	if err != nil {
		panic(err)
	}
	next()
}

func FormValue(vars rack.Vars, key string) string {
	r := rack.GetRequest(vars)
	f := r.Form
	if f == nil {
		r.ParseForm()
		f := r.Form
		if f == nil {
			return ""
		}
	}

	return f.Get(key)
}
