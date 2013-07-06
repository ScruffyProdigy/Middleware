/*
	router creates a branch based routing system
*/
package router

import (
	"github.com/ScruffyProdigy/TheRack/rack"
)

//Router is a Middleware that represents a piece of a URL
type Router struct {
	parent    *Router         // the piece of the URL that came before
	subroutes []HasRouter     // the pieces of the URL that might come after
	Action    rack.Middleware // what happens when this is the end of the URL
	Routing   Signaler        // how to tell whether or not this is the correct section of URL
	Name      Namer           // how to tell what the string that represented this section of URL was
}

// In order to allow non-Router's be added as subroutes, they just need to fulfill this interface
type HasRouter interface {
	Router() *Router                     // we need to be able to get an actual router sometimes
	SetParent(*Router)                   // when we add you as a subroute, you need to know who your parent is
	Route(map[string]interface{}) string // we need to know what route to get to you
}

//New() will return a blank router
func New() *Router {
	this := new(Router)
	this.subroutes = make([]HasRouter, 0)
	return this
}

//Route() will return a router filled with the basic information
func Route(routing Signaler, action rack.Middleware, name Namer) *Router {
	this := New()
	this.Routing = routing
	this.Action = action
	this.Name = name
	return this
}

//Router() Just returns itself to fulfill the HasRouter Interface
func (this *Router) Router() *Router {
	return this
}

//SetParent() sets the parent to help with figuring out the name
func (this *Router) SetParent(parent *Router) {
	this.parent = parent
}

// Route() will recursively figure out what route leads here
func (this *Router) Route(vars map[string]interface{}) (route string) {
	if this.parent != nil {
		route = this.parent.Route(vars)
	}
	route += this.Name.Run(vars)
	route += "/"
	return
}

//Run implements the rack.Middleware interface
func (this *Router) Run(vars map[string]interface{}, next func()) {
	if V(vars).CurrentSection() == RouteEnd {
		this.Action.Run(vars, next)
	} else {
		for _, subroute := range this.subroutes {
			router := subroute.Router()
			if router.Routing.Run(vars) {
				V(vars).nextSection()
				router.Run(vars, next)
				return
			}
		}
		next()
	}
}

//AddRoute adds one or more subroutes to this route
func (this *Router) AddRoute(routes ...HasRouter) {
	for _, route := range routes {
		this.subroutes = append(this.subroutes, route)
		route.SetParent(this)
	}
}

//a Signaler can tell whether the Router should continue down this subroute
type Signaler interface {
	Run(vars map[string]interface{}) bool
}

//a SignalFunc is the simplest possible Signaler, any func of this form can be cast into a Signaler with this
type SignalFunc func(map[string]interface{}) bool

//Run implements the Signaler interface
func (this SignalFunc) Run(vars map[string]interface{}) bool {
	return this(vars)
}

//a Namer tells what the piece of the URL should look like
type Namer interface {
	Run(vars map[string]interface{}) string
}

//a NamerFunc allows any func of this form can be cast into a Namer
type NamerFunc func(vars map[string]interface{}) string

// Run implements the Namer interface
func (this NamerFunc) Run(vars map[string]interface{}) string {
	return this(vars)
}

//a NameString is the simplest possible Namer, it allows any string to act as a static Namer
type NameString string

//Run implements the Namer interface
func (this NameString) Run(vars map[string]interface{}) string {
	return string(this)
}
