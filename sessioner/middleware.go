package sessioner

import (
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheRack/rack"
)

const (
	sessionIndex = "session"
	flashesIndex = "flashes"
)

/*
	Middleware is the Middleware function that inserts a Session variable as "Session" into Rack variables
	This allows all later Middleware to have persistent effects
*/
var Middleware rack.Func = func(vars map[string]interface{}, next func()) {
	r := (httper.V)(vars).GetRequest()

	session := get(r)
	vars[sessionIndex] = Session(session)

	vars[flashesIndex] = session.Clear(flashesIndex)
	_, isStrings := vars[flashesIndex].([]string)
	if !isStrings {
		vars[flashesIndex] = []string{}
	}

	next()

	w := (httper.V)(vars).FilledResponse()
	session.save(w)
	w.Save()
}

type V map[string]interface{}

func (vars V) Set(k, v interface{}) {
	vars[sessionIndex].(Session).Set(k, v)
}

func (vars V) Get(k interface{}) interface{} {
	return vars[sessionIndex].(Session).Get(k)
}

func (vars V) Clear(k interface{}) interface{} {
	return vars[sessionIndex].(Session).Clear(k)
}

func (vars V) AddFlash(s string) {
	a, isStrings := vars.Get(flashesIndex).([]string)
	if !isStrings {
		a = []string{s}
	} else {
		a = append(a, s)
	}
	vars.Set(flashesIndex, a)
}
