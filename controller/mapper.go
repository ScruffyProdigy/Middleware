package controller

import (
	"bitbucket.org/pkg/inflect"
	"github.com/ScruffyProdigy/Middleware/redirecter"
	"github.com/ScruffyProdigy/Middleware/renderer"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"reflect"
	"strings"
)

type Urler interface {
	Url() string
}

type dispatchAction struct {
	descriptor
	name   string
	action rack.Middleware
}

type descriptor struct {
	t                  reflect.Type
	varName, routeName string
}

func createDescriptor(iface interface{}) descriptor {
	pv := reflect.ValueOf(iface)
	v := reflect.Indirect(pv)
	t := v.Type()

	var varName string
	if nameGetter, hasVarName := iface.(interface {
		VarName() string
	}); hasVarName {
		varName = nameGetter.VarName()
	} else {
		varName = inflect.Singularize(t.Name())
	}

	var routeName string
	if nameGetter, hasRouteName := iface.(interface {
		RouteName() string
	}); hasRouteName {
		routeName = nameGetter.RouteName()
	} else {
		routeName = strings.ToLower(t.Name())
	}

	return descriptor{
		t:         t,
		varName:   varName,
		routeName: routeName,
	}
}

func (this descriptor) addDispatchAction(funcs map[string]rack.Middleware, name string, method reflect.Method) {
	name = strings.ToLower(name)
	d := new(dispatchAction)
	d.descriptor = this
	d.name = name
	d.action = rack.Func(func(vars map[string]interface{}, next func()) {
		copy := reflect.New(this.t)
		mapper := copy.Interface().(ResourceController)
		mapper.SetRackVars(this, vars, next)
		method.Func.Call([]reflect.Value{reflect.Indirect(copy)})
		if !mapper.IsFinished() {
			next()
		}
	})

	funcs[name] = d
}

func (this dispatchAction) Run(vars map[string]interface{}, next func()) {
	actions := rack.New()
	actions.Add(this.action)
	switch (httper.V)(vars).GetRequest().Method {
	case "GET":
		//if it was a get, the default action should be to render the template corresponding with the action
		actions.Add(renderer.Renderer{this.routeName + "/" + this.name})
	case "POST", "PUT":
		//if it was a put or a post, we the default action should be to redirect to the affected item
		actions.Add(rack.Func(func(vars map[string]interface{}, next func()) {
			urler, isUrler := vars[this.varName].(Urler)
			if !isUrler {
				panic("Object doesn't have an URL to direct to")
			}
			(redirecter.V)(vars).Redirect(urler.Url())
		}))
	case "DELETE":
		urler, isUrler := vars[this.varName].(Urler)
		if !isUrler {
			panic("Object doesn't have an URL to direct to")
		}
		url := urler.Url()

		if url[len(url)-1:] == "/" {
			url = url[:len(url)-1]
		}

		i := strings.LastIndex(url, "/")
		if i != -1 {
			url = url[:i]
		}
		actions.Add(rack.Func(func(vars map[string]interface{}, next func()) {
			next()
		}))
		actions.Add(redirecter.Redirecter{url})
	default:
		panic("Unknown method")
	}
	actions.Run(vars, next)
}

func isControlFunc(m reflect.Method) bool {
	t := m.Type
	if t.Kind() != reflect.Func { //it should be a function
		return false
	}
	if t.NumIn() != 1 { //it should have one input parameter (the 'this' controller)
		return false
	}
	if t.NumOut() != 0 { //it should have no output parameters
		return false
	}
	return true
}

func (this descriptor) GetRestMap() (restfuncs map[string]rack.Middleware) {
	restfuncs = make(map[string]rack.Middleware)

	for _, funcName := range []string{"Index", "Create", "New", "Show", "Edit", "Update", "Destroy"} {
		method, methodExists := this.t.MethodByName(funcName)
		if methodExists && isControlFunc(method) {
			this.addDispatchAction(restfuncs, funcName, method)
		}
	}

	return
}

type mapList struct {
	all, get, put, post, delete map[string]rack.Middleware
}

func (this descriptor) GetGenericMapList(functype string) (funcs mapList) {
	funcs.all = this.GetGenericMap(functype)
	funcs.get = this.GetGenericMap("Get" + functype)
	funcs.put = this.GetGenericMap("Put" + functype)
	funcs.post = this.GetGenericMap("Post" + functype)
	funcs.delete = this.GetGenericMap("Delete" + functype)
	return
}

func (this descriptor) GetGenericMap(functype string) (funcs map[string]rack.Middleware) {
	funcs = make(map[string]rack.Middleware)
	typelen := len(functype)

	for i, c := 0, this.t.NumMethod(); i < c; i = i + 1 {
		method := this.t.Method(i)
		//if the first part of the name is whatever we're looking for, and it's a control function
		if len(method.Name) >= typelen && method.Name[:typelen] == functype && isControlFunc(method) {
			this.addDispatchAction(funcs, method.Name[typelen:], method)
		}
	}
	return
}
