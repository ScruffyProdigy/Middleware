package controller

import (
	"github.com/ScruffyProdigy/Middleware/router"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"net/http"
	"strings"
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
	indexer Indexer
}

func (this memberSignaler) Run(vars map[string]interface{}) bool {
	id := (router.V)(vars).CurrentSection()
	result, found := this.indexer.Find(id, vars)

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
		the usable keys in the map are: "index","new","create","show","edit","update", and "destroy"
	variablename: If we are drilling down into a member of the resource, we will add a variable to the rack variables, and this will be the name that it will set
	getter:	if we need to get a member resource, you'll have to help us;  we'll give you the string representing the ID, you give us the resource
*/

func AddMapRoutes(superroute *router.Router, routemap map[string]rack.Middleware, methodfinder func(string, rack.Middleware) *router.Router) {
	for name, action := range routemap {
		if name != "" {
			superroute.AddRoute(methodfinder(name, action))
		}
	}
}

func AddMapListRoutes(superroute *router.Router, maplist mapList) {
	AddMapRoutes(superroute, maplist.get, router.Get)
	AddMapRoutes(superroute, maplist.put, router.Put)
	AddMapRoutes(superroute, maplist.post, router.Post)
	AddMapRoutes(superroute, maplist.delete, router.Delete)
	AddMapRoutes(superroute, maplist.all, router.All)
}

func firstNonNilMiddleware(options []rack.Middleware) rack.Middleware {
	for _, option := range options {
		if option != nil {
			return option
		}
	}
	return nil
}

func NewResource(m ResourceController) *ControllerShell {
	resource := new(ControllerShell)

	descriptor := createDescriptor(m)

	restfuncs := descriptor.GetRestMap()
	memberfuncs := descriptor.GetGenericMapList("Member")
	collectionfuncs := descriptor.GetGenericMapList("Collection")

	resource.Member = router.NewRouter()
	resource.Member.Routing = memberSignaler{varName: descriptor.varName, indexer: m}
	memberactions := splitter{}
	memberactions.get = firstNonNilMiddleware([]rack.Middleware{restfuncs["show"], memberfuncs.get[""], memberfuncs.all[""]})
	memberactions.post = firstNonNilMiddleware([]rack.Middleware{memberfuncs.post[""], memberfuncs.all[""]})
	memberactions.put = firstNonNilMiddleware([]rack.Middleware{restfuncs["update"], memberfuncs.put[""], memberfuncs.all[""]})
	memberactions.delete = firstNonNilMiddleware([]rack.Middleware{restfuncs["destroy"], memberfuncs.delete[""], memberfuncs.all[""]})
	resource.Member.Action = memberactions

	if restfuncs["edit"] != nil {
		memberfuncs.get["edit"] = restfuncs["edit"]
	}
	AddMapListRoutes(resource.Member, memberfuncs)

	resource.Collection = router.NewRouter()
	resource.Collection.Routing = collectionSignaler{name: descriptor.routeName}
	collectionactions := splitter{}
	collectionactions.get = firstNonNilMiddleware([]rack.Middleware{restfuncs["index"], collectionfuncs.get[""], collectionfuncs.all[""]})
	collectionactions.post = firstNonNilMiddleware([]rack.Middleware{restfuncs["create"], collectionfuncs.post[""], collectionfuncs.all[""]})
	collectionactions.put = firstNonNilMiddleware([]rack.Middleware{collectionfuncs.put[""], collectionfuncs.all[""]})
	collectionactions.delete = firstNonNilMiddleware([]rack.Middleware{collectionfuncs.delete[""], collectionfuncs.all[""]})
	resource.Collection.Action = collectionactions

	if restfuncs["new"] != nil {
		collectionfuncs.get["new"] = restfuncs["new"]
	}
	AddMapListRoutes(resource.Collection, collectionfuncs)

	resource.Collection.AddRoute(resource.Member)

	return resource
}

func (this *ControllerShell) Router() *router.Router {
	return this.Collection
}

//use this to create the root of your routes (e.g. for http://example.com/)
func NewRoot(m rack.Middleware) *router.Router {
	return router.BasicRoute("", m)
}

//use this to create a route namespace (i.e. "admin").
// name is the string used for the namespace.
// intermediate will be ran if the namespace is used (http://example.com/admin/posts).
// final will be ran if the namespace is the destination (http://example.com/admin)
func NewNamespace(name string, intermediate router.Signaler, final rack.Middleware) *router.Router {
	r := router.NewRouter()
	r.Routing = router.SignalFunc(func(vars map[string]interface{}) bool {
		section := router.V(vars).CurrentSection()

		if !router.V(vars).IsCaseSensitive() {
			name = strings.ToLower(name)
		}

		if section != name {
			return false
		}

		if intermediate != nil {
			return intermediate.Run(vars)
		}

		return true
	})

	r.Action = final
	return r
}
