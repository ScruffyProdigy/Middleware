package router

import (
	"github.com/ScruffyProdigy/TheRack/rack"
)

type Router struct {
	subroutes []*Router
	Action    rack.Middleware
	Routing   Signaler
}

func NewRouter() *Router {
	this := new(Router)
	this.subroutes = make([]*Router, 0)
	return this
}

func Route(routing Signaler, action rack.Middleware) *Router {
	this := NewRouter()
	this.Routing = routing
	this.Action = action
	return this
}

func (this *Router) Run(vars map[string]interface{}, next func()) {
	if V(vars).CurrentSection() == RouteEnd {
		this.Action.Run(vars, next)
	} else {
		for _, subroute := range this.subroutes {
			if subroute.Routing.Run(vars) {
				V(vars).nextSection()
				subroute.Run(vars, next)
				return
			}
		}
		next()
	}
}

func (this *Router) AddRoute(r ...*Router) {
	this.subroutes = append(this.subroutes, r...)
}

type Signaler interface {
	Run(vars map[string]interface{}) bool
}

type SignalFunc func(map[string]interface{}) bool

func (this SignalFunc) Run(vars map[string]interface{}) bool {
	return this(vars)
}
