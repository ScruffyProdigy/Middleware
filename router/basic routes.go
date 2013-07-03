package router

import (
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"strings"
)

type basicRoute struct {
	methodMatters bool
	method        string
	name          string
}

func (this *basicRoute) Run(vars map[string]interface{}) bool {
	sec := V(vars).CurrentSection()

	name := this.name
	if !V(vars).IsCaseSensitive() {
		name = strings.ToLower(name)
	}

	if sec != name {
		return false
	}

	if this.methodMatters {
		req := (httper.V)(vars).GetRequest()
		if req.Method != this.method {
			return false
		}
	}
	return true
}

func createMethodRoute(method, name string, m rack.Middleware) (result *Router) {
	result = New()
	result.Routing = &basicRoute{method: method, name: name, methodMatters: true}
	result.Action = m
	result.Name = NameString(name)
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

//BasicRoute returns a Router that acts with a static string
func BasicRoute(name string, m rack.Middleware) (result *Router) {
	result = New()
	result.Routing = &basicRoute{name: name, method: "", methodMatters: false}
	result.Action = m
	result.Name = NameString(name)
	return
}

var All = BasicRoute // a shortcut added for convenience to go along with Get,Put,Post, and Delete
