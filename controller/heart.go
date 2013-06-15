package controller

import (
	"github.com/ScruffyProdigy/Middleware/redirecter"
	"github.com/ScruffyProdigy/Middleware/renderer"
	"github.com/ScruffyProdigy/Middleware/sessioner"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"net/http"
	"strings"
)

const (
	finishedIndex = "controllerheartisfinished"
)

// When Creating a Controller, you MUST put an anonymous controller.Heart into your controller (unless you really know what you're doing)
// Not only do some of the functions require some a couple of the default methods
type Heart struct {
	descriptor
	Vars   map[string]interface{}
	finish func()
}

type ModelMap interface {
	SetRackVars(descriptor, map[string]interface{}, func())
	IsFinished() bool
	SetFinished()
	SetUnfinished()
}

// this is how we hide the rack variables from the controllers who don't really care so much about these
// they are later accessible in case you need them, but for the most part, you can just ignore these
func (this *Heart) SetRackVars(t descriptor, vars map[string]interface{}, next func()) {
	this.descriptor = t
	this.Vars = vars
	this.finish = next
}

// if you are calling another Rack Middleware, you should call this function to get the variables it will need to run
func (this Heart) getRackFuncVars() (vars map[string]interface{}, next func()) {
	return this.Vars, this.finish
}

func (this Heart) IsFinished() bool {
	isFinished, isValid := this.Vars[finishedIndex].(bool)
	return isValid && isFinished
}

func (this Heart) SetFinished() {
	this.Vars[finishedIndex] = true
}

func (this Heart) SetUnfinished() {
	this.Vars[finishedIndex] = true
}

func (this Heart) finishingFunc(action func()) {
	if this.IsFinished() {
		panic("called a second finishing function")
	}
	action()
	this.SetFinished()
}

func (this Heart) Finish() {
	this.finishingFunc(func() {
		this.finish()
	})
}

// this should be the return value for most of your Create control functions
// you should pass it the resource you created
// since the default variable wasn't set because we didn't get into a specific resource
// this will set the default variable to the resource you just created
func (this Heart) RespondWith(object interface{}) {
	this.finishingFunc(func() {
		this.Vars[this.varName] = object
		this.Finish()
	})
}

// if you have a piece of middleware that you want to respond with
// return this instead of Finish along with the middleware you want to run
func (this Heart) FinishWithMiddleware(m rack.Middleware) {
	this.finishingFunc(func() {
		m.Run(this.getRackFuncVars())
	})
}

// if things don't go according to plan, you can redirect somewhere else
// return this instead of Finish along with where you want to redirect to
func (this Heart) RedirectTo(url string) {
	this.finishingFunc(func() {
		(redirecter.V)(this.Vars).Redirect(url)
	})
}

// use this if you want to render something other than the default template
// return this instead of Finish along with the template you want to render
func (this Heart) Render(tmpl string) {
	this.finishingFunc(func() {
		if !strings.Contains(tmpl, "/") {
			tmpl = this.routeName + "/" + tmpl
		}
		(renderer.V)(this.Vars).Render(tmpl)
	})
}

func (this Heart) AddFlash(flash string) {
	(sessioner.V)(this.Vars).AddFlash(flash)
}

func (this Heart) Session() sessioner.V {
	return (sessioner.V)(this.Vars)
}

// this will get the form value from the form that was passed in
func (this Heart) GetFormValue(value string) string {
	return (httper.V)(this.Vars).GetRequest().FormValue(value)
}

func (this Heart) NotAuthorized() {
	this.finishingFunc(func() {})

	(httper.V)(this.Vars).Status(http.StatusUnauthorized)
}

// this is used to set variables
// for the most part, this is used so that the template will have access to more variables when rendering
func (this Heart) Set(k string, v interface{}) {
	this.Vars[k] = v
}

// this is used to get previously set variables
// the most common variable to get is the one we stored for you for all member methods
func (this Heart) GetVar(k string) interface{} {
	return this.Vars[k]
}
