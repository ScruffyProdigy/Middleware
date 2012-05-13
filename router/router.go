package router

import (
	"github.com/HairyMezican/TheRack/rack"
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

func (this *Router) Run(vars rack.Vars, next func()) {
	if CurrentSection(vars) == RouteEnd {
		this.Action.Run(vars, next)
	} else {
		for _, subroute := range this.subroutes {
			if subroute.Routing.Run(vars) {
				nextSection(vars)
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
	Run(vars rack.Vars) bool
}

type SignalFunc func(rack.Vars) bool

func (this SignalFunc) Run(vars rack.Vars) bool {
	return this(vars)
}
