package parser

import (
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheRack/rack"
)

var Form rack.Func = func(vars map[string]interface{}, next func()) {
	V(vars).Parse()
	next()
}

type Multipart struct {
	MaxSize int64
}

func (this Multipart) Run(vars map[string]interface{}, next func()) {
	V(vars).ParseMultipart(this.MaxSize)
	next()
}

type V map[string] interface{}

func (vars V) Parse() {
	r := (httper.V)(vars).GetRequest()
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
}

func (vars V) ParseMultipart(maxSize int64) {
	r := (httper.V)(vars).GetRequest()
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		panic(err)
	}
}

func (vars V) FormValue(key string) string {
	r := (httper.V)(vars).GetRequest()
	f := r.Form
	if f == nil {
		vars.Parse()
		f := r.Form
		if f == nil {
			return ""
		}
	}

	return f.Get(key)
}
