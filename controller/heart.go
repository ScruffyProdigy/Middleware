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

type Model interface {
	ID() string
}

//An Indexer should be able to take an index and return a resource
//s is the index of the resource
//vars is the list of variables that have been assembled so far
//resource is the found resource
//found is whether or not anything was actually found
type Indexer interface {
	Find(s string, vars map[string]interface{}) (resource Model, found bool)
}

//HasHeart is an interface that exposes the parts of a controller that the mapper needs that can be completed by having a Heart in your class
type HasHeart interface {
	SetRackVars(descriptor, map[string]interface{}, func())
	IsFinished() bool
}

//A ResourceController generally needs a Heart
// and needs to be able to index resources.
// and should be able to describe the route that leads to it
//You should also add in control functions.
//The following Restful routes will automatically be created for you if you create the appropriate functions:
// • Index() (Collection - GET /)
// • Create() (Collection - POST /)
// • New() (Collection - GET /new)
// • Show() (Member - GET /)
// • Update() (Member - PUT /)
// • Edit() (Member - GET /edit)
// • Destroy() (Member - DELETE /)
//
// To add other routes, the following naming scheme is used:
// MemberFoo() - Adds a /foo route for all members of the resource
// CollectionBar() - Adds a /bar route for the collection
// GetMemberBaz() - Adds a /baz route for members of the resource, but only if it is a GET request
// PostCollectionQux() - Adds a /qux route for the collection, but only if it is a POST request
// PutMemberSoul() - Adds a /soul route for members of the resource, but only if it is a PUT request
// DeleteCollection() - Responds to DELETE requests for the collection
type ResourceController interface {
	HasHeart
	Indexer
}

// this is how we hide the rack variables from the controllers who don't really care so much about these
// they are later accessible in case you need them, but for the most part, you can just ignore these
func (this *Heart) SetRackVars(t descriptor, vars map[string]interface{}, next func()) {
	this.descriptor = t
	this.Vars = vars
	this.finish = next
}

func (this Heart) getRackFuncVars() (vars map[string]interface{}, next func()) {
	return this.Vars, this.finish
}

//returns whether or not a control function has called any of the finishing functions
func (this Heart) IsFinished() bool {
	isFinished, isValid := this.Vars[finishedIndex].(bool)
	return isValid && isFinished
}

//this makes sure that only one finishing functions gets called per control function.
//if this has been called before within a control function, it panics
func (this Heart) finishingFunc(action func()) {
	if this.IsFinished() {
		panic("called a second finishing function")
	}
	defer func() {
		this.Vars[finishedIndex] = true
	}()
	action()
}

/**************

Finishing Functions

**************/

//Finish() is the default finishing functions.
// It just declares that the control is finished, and goes along with the default finish
// Most of the time, this function is optional.
// If you don't call it, it will be called for you after your control function finishes
func (this Heart) Finish() {
	this.finishingFunc(func() {
		this.finish()
	})
}

//RespondWith() sets the base variable before continuing with the default finish.
// this should be the return value for most of your Create control functions;
// you should pass it the resource you created.
// Since the default variable wasn't set because we didn't get into a specific resource,
// this will set the default variable to the resource you just created
func (this Heart) RespondWith(object interface{}) {
	this.finishingFunc(func() {
		this.Vars[this.varName] = object
		this.Finish()
	})
}

//FinishWithMiddleware() replaces the default finish with any other middleware you might have on hand.
// if you have a piece of middleware that you want to respond with
// call this instead of Finish along with the middleware you want to run
func (this Heart) FinishWithMiddleware(m rack.Middleware) {
	this.finishingFunc(func() {
		m.Run(this.getRackFuncVars())
	})
}

//RedirectTo() replaces the default finish with a redirection.
// If things don't go according to plan, you can redirect the user somewhere else.
// Call this instead of Finish along with where you want to redirect to
func (this Heart) RedirectTo(url string) {
	this.finishingFunc(func() {
		(redirecter.V)(this.Vars).Redirect(url)
	})
}

//Render() replaces the default finish with a rendering.
// use this if you want to render something other than the default template.
// Call this instead of Finish along with the template you want to render
func (this Heart) Render(tmpl string) {
	this.finishingFunc(func() {
		if !strings.Contains(tmpl, "/") {
			tmpl = this.routeName + "/" + tmpl
		}
		(renderer.V)(this.Vars).Render(tmpl)
	})
}

//AddFlash() adds a flash to the session.
// Useful before a redirect, as all flashes will be retrievable then,
// but will be erased from the session so as to be gone before the user does anything else
func (this Heart) AddFlash(flash string) {
	this.Session().AddFlash(flash)
}

//Session() gives access to the session variables
// both for setting and getting
func (this Heart) Session() sessioner.V {
	return (sessioner.V)(this.Vars)
}

// FormValue() will get the form value from the form that was passed in
func (this Heart) FormValue(value string) string {
	return (httper.V)(this.Vars).GetRequest().FormValue(value)
}

// NotAuthorized() replaces the default finish, and simply tells the user that they cannot access the current page
// use this if they need to be an admin or something
func (this Heart) NotAuthorized() {
	this.finishingFunc(func() {})

	(httper.V)(this.Vars).Status(http.StatusUnauthorized)
}

// this is used to set rack variables
// for the most part, this is used so that the template will have access to more variables when rendering
func (this Heart) Set(k string, v interface{}) {
	this.Vars[k] = v
}

// this is used to get previously set rack variables
// the most common variable to get is the one we stored for you for all member methods
func (this Heart) GetVar(k string) interface{} {
	return this.Vars[k]
}
