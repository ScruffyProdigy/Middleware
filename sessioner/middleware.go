package sessioner

import (
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
var Middleware rack.Func = func(vars rack.Vars, next func()) {
	r := rack.GetRequest(vars)

	session := get(r)
	vars[sessionIndex] = Session(session)

	vars[flashesIndex] = session.Clear(flashesIndex)
	_, isStrings := vars[flashesIndex].([]string)
	if !isStrings {
		vars[flashesIndex] = []string{}
	}

	next()

	w := rack.CreateResponse(vars)
	session.save(w)
	w.Save()
}

func Set(vars rack.Vars, k, v interface{}) {
	vars[sessionIndex].(Session).Set(k, v)
}

func Get(vars rack.Vars, k interface{}) interface{} {
	return vars[sessionIndex].(Session).Get(k)
}

func Clear(vars rack.Vars, k interface{}) interface{} {
	return vars[sessionIndex].(Session).Clear(k)
}

func AddFlash(vars rack.Vars, s string) {
	a, isStrings := Get(vars, "flash").([]string)
	if !isStrings {
		a = []string{s}
	} else {
		a = append(a, s)
	}
	Set(vars, "flash", a)
}
