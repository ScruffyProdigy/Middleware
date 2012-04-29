package router

import (
	"github.com/HairyMezican/TheRack/rack"
	"net/http"
	"strings"
)

type methodRoute struct {
	method string
	name   string
}

func (this *methodRoute) Run(req *http.Request, vars rack.Vars) bool {
	sec := vars.Apply(CurrentSection).(string)
	if sec == this.name && req.Method == this.method {
		return true
	}
	return false
}

func createMethodRoute(method, name string, m rack.Middleware) (result *Router) {
	result = NewRouter()
	result.Routing = &methodRoute{method: method, name: strings.ToLower(name)}
	result.Action = m
	return
}

//Get provides a RouteTerminal that will direct a GET request to a specified handler
func Get(name string, m rack.Middleware) (result *Router) {
	return createMethodRoute("GET", name, m)
}

//Post provides a RouteTermianl that will direct a POST request to a specified handler
func Post(name string, m rack.Middleware) (result *Router) {
	return createMethodRoute("POST", name, m)
}

//Put provides a RouteTerminal that will direct a PUT request to a specified handler
func Put(name string, m rack.Middleware) (result *Router) {
	return createMethodRoute("PUT", name, m)
}

//Delete provides a RouteTerminal that will direct a DELETE request to specified handler
func Delete(name string, m rack.Middleware) (result *Router) {
	return createMethodRoute("DELETE", name, m)
}

type simpleRoute struct {
	name string
}

func (this *simpleRoute) Run(req *http.Request, vars rack.Vars) bool {
	sec := strings.ToLower(vars.Apply(CurrentSection).(string))
	if sec == this.name {
		return true
	}
	return false
}

func BasicRoute(name string, m rack.Middleware) (result *Router) {
	result = NewRouter()
	result.Routing = &simpleRoute{name: strings.ToLower(name)}
	result.Action = m
	return
}

var All = BasicRoute // a shortcut added for convenience to go along with Get,Put,Post, and Delete
