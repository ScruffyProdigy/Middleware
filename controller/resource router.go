package controller

import (
	"github.com/HairyMezican/Middleware/router"
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheRack/rack"
	"net/http"
)

/*
a ResourceRouter assumes that it represents a RESTful resource, and will process it as such
it also allows you to add non-RESTful member and collection routes by exposing a route branch for each
*/

type ControllerShell struct {
	Collection *router.Router //you can add non-RESTful collection-level routes here
	Member     *router.Router //you can add non-RESTful member-level routes here
}

type splitter struct {
	get, post, put, delete rack.Middleware
}

func (this splitter) Run(vars map[string]interface{}, next func()) {
	var result rack.Middleware
	switch (httper.V)(vars).GetRequest().Method {
	case "GET":
		result = this.get
	case "POST":
		result = this.post
	case "PUT":
		result = this.put
	case "DELETE":
		result = this.delete
	default:
		(httper.V)(vars).Status(http.StatusBadRequest)
	}
	if result == nil {
		//that particular method wasn't set, but perhaps a later middleware will take care of it
		next()
	} else {
		result.Run(vars, next)
	}
}

type memberSignaler struct {
	varName string
	indexer func(string, map[string]interface{}) (interface{}, bool)
}

func (this memberSignaler) Run(vars map[string]interface{}) bool {
	id := (router.V)(vars).CurrentSection()
	result, found := this.indexer(id, vars)
	if !found {
		return false
	}

	vars[this.varName] = result

	return true
}

type collectionSignaler struct {
	name string
}

func (this collectionSignaler) Run(vars map[string]interface{}) bool {
	section := (router.V)(vars).CurrentSection()
	if section == this.name {
		return true
	}
	return false
}

/*
	Resource will return a RESTful Resource Router
	it expects
	name: a string that represents the name of the resource.  This is used in the routing process
	restfuncs: the RESTful routes that this resource expects to handle
		the usable keys in the map are: "index","new","create","show","edit","update", and "delete"
	variablename: If we are drilling down into a member of the resource, we will add a variable to the rack variables, and this will be the name that it will set
	getter:	if we need to get a member resource, you'll have to help us;  we'll give you the string representing the ID, you give us the resource
*/

func AddMapRoutes(superroute *router.Router, routemap map[string]rack.Middleware, methodfinder func(string, rack.Middleware) *router.Router) {
	for name, action := range routemap {
		superroute.AddRoute(methodfinder(name, action))
	}
}

func AddMapListRoutes(superroute *router.Router, maplist mapList) {
	AddMapRoutes(superroute, maplist.get, router.Get)
	AddMapRoutes(superroute, maplist.put, router.Put)
	AddMapRoutes(superroute, maplist.post, router.Post)
	AddMapRoutes(superroute, maplist.delete, router.Delete)
	AddMapRoutes(superroute, maplist.all, router.All)
}

func RegisterController(m ModelMap, routeName, varName string, indexer func(s string, vars map[string]interface{}) (interface{}, bool)) *ControllerShell {
	resource := new(ControllerShell)

	descriptor := createDescriptor(m, routeName, varName)

	restfuncs := descriptor.GetRestMap()
	memberfuncs := descriptor.GetGenericMapList("Member")
	collectionfuncs := descriptor.GetGenericMapList("Collection")

	resource.Member = router.NewRouter()
	resource.Member.Routing = memberSignaler{varName: varName, indexer: indexer}
	resource.Member.Action = splitter{get: restfuncs["show"], put: restfuncs["update"], delete: restfuncs["destroy"]}

	if restfuncs["edit"] != nil {
		memberfuncs.get["edit"] = restfuncs["edit"]
	}
	AddMapListRoutes(resource.Member, memberfuncs)

	resource.Collection = router.NewRouter()
	resource.Collection.Routing = collectionSignaler{name: routeName}
	resource.Collection.Action = splitter{get: restfuncs["index"], post: restfuncs["create"]}

	if restfuncs["new"] != nil {
		collectionfuncs.get["new"] = restfuncs["new"]
	}
	AddMapListRoutes(resource.Collection, collectionfuncs)

	resource.Collection.AddRoute(resource.Member)

	return resource
}

func (this ControllerShell) AddTo(superroute *router.Router) {
	superroute.AddRoute(this.Collection)
}

func (this ControllerShell) AddAsSubresource(parent *ControllerShell) {
	parent.Member.AddRoute(this.Collection)
}

func (this ControllerShell) AddAsSubmethod(parent *ControllerShell) {
	parent.Collection.AddRoute(this.Collection)
}
