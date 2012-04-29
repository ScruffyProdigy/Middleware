package router

import (
	"github.com/HairyMezican/TheRack/rack"
	"net/http"
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

func (this *Router) Run(r *http.Request, vars rack.Vars, next rack.Next) (status int, header http.Header, message []byte) {
	if vars.Apply(CurrentSection) == "/" {
		return this.Action.Run(r, vars, next)
	}
	for _, subroute := range this.subroutes {
		if subroute.Routing.Run(r, vars) {
			vars.Apply(nextSection)
			return subroute.Run(r, vars, next)
		}
	}
	return next()
}

func (this *Router) AddRoute(r ...*Router) {
	this.subroutes = append(this.subroutes, r...)
}

var Root *Router = NewRouter()

type Signaler interface {
	Run(r *http.Request, vars rack.Vars) bool
}

type SignalFunc func(*http.Request, rack.Vars) bool

func (this SignalFunc) Run(r *http.Request, vars rack.Vars) bool {
	return this(r, vars)
}
