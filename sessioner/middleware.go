/*
Session is a middleware that will load and save session variables
*/
package sessioner

import (
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
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

// V is a type you can cast your vars to in order to access the following functions
type V map[string]interface{}

// Session() will get access to the actual session.
// Usually you should just use one of the later functions to directly do whatever you want to do
func (vars V) Session() Session {
	session, _ := vars[sessionIndex].(Session)
	return session
}

// Set() will set a session variable
func (vars V) Set(k, v interface{}) {
	if session := vars.Session(); session != nil {
		session.Set(k, v)
	}

}

// Get() will get a session variable
func (vars V) Get(k interface{}) interface{} {
	if session := vars.Session(); session != nil {
		return session.Get(k)
	}
	return nil
}

// Clear will delete a session variable (and return its old value)
func (vars V) Clear(k interface{}) interface{} {
	if session := vars.Session(); session != nil {
		return session.Clear(k)
	}
	return nil
}

// AddFlash will add a Flash.
// This is useful right before a redirect
func (vars V) AddFlash(s string) {
	a, isStrings := vars.Get(flashesIndex).([]string)
	if !isStrings {
		a = []string{s}
	} else {
		a = append(a, s)
	}
	vars.Set(flashesIndex, a)
}

// Flashes will return all Flashes that were added during the previous session.
// This is useful right after a redirect
func (vars V) Flashes() []string {
	result := vars[flashesIndex]
	if result != nil {
		if strings, ok := result.([]string); ok {
			return strings
		}
	}
	return []string{}
}
