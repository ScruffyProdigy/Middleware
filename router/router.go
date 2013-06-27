package router

import (
	"github.com/ScruffyProdigy/TheRack/rack"
)

type Router struct {
	parent    *Router
	subroutes []HasRouter
	Action    rack.Middleware
	Routing   Signaler
	Name      Namer
}

type HasRouter interface {
	Router() *Router
	SetParent(*Router)
	Route(map[string]interface{}) string
}

func New() *Router {
	this := new(Router)
	this.subroutes = make([]HasRouter, 0)
	return this
}

func Route(routing Signaler, action rack.Middleware, name Namer) *Router {
	this := New()
	this.Routing = routing
	this.Action = action
	this.Name = name
	return this
}

func (this *Router) Router() *Router {
	return this
}

func (this *Router) SetParent(parent *Router) {
	this.parent = parent
}

func (this *Router) Route(vars map[string]interface{}) string {
	return this.Name.Run(vars, func() string {
		if this.parent != nil {
			return this.parent.Route(vars)
		}
		return ""
	})
}

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

func (this *Router) AddRoute(routes ...HasRouter) {
	for _, route := range routes {
		this.subroutes = append(this.subroutes, route)
		route.SetParent(this)
	}
}

type Signaler interface {
	Run(vars map[string]interface{}) bool
}

type SignalFunc func(map[string]interface{}) bool

func (this SignalFunc) Run(vars map[string]interface{}) bool {
	return this(vars)
}

type Namer interface {
	Run(vars map[string]interface{}, prev func() string) string
}

type NamerFunc func(vars map[string]interface{}, prev func() string) string

func (this NamerFunc) Run(vars map[string]interface{}, prev func() string) string {
	return this(vars, prev)
}

type NameString string

func (this NameString) Run(vars map[string]interface{}, prev func() string) string {
	return prev() + string(this) + "/"
}
