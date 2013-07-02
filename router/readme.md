#	Routes
This provides a branching based router that goes through every directory specified in the requested url and finds a subrouter until it gets to the end, then calls whichever middleware is associated with the last subrouter found

## 	Installation
`go get github.com/ScruffyProdigy/Middleware/router`

##  Usage

* Create a bunch of routes
	* You can use router.BasicRoute() if you have a static string to match
		* The first parameter is the static string defining the route
		* The second parameter is the middleware that gets ran if this is the end route
	* You can use router.Get(), router.Post(), router.Put(), or router.Delete() if you want to match on both a static string and a request type
		* The parameters are the same as router.BasicRoute()
	* You can call router.Route() to create a custom route
		* The first parameter is a Signaler
			* The signaler is an object which has a method Run()
				* Run takes in a set of vars, and returns whether or not this is the correct route to follow
				* Generally, you will need to call (router.V)(vars).CurrentSection() and use that to figure out if we're on the right track
				* You can also check whether or not a user has permissions to go down this branch
				* You can also store variables in the vars for later use
		* The second parameter is a Middleware
			* The Middleware only gets called if this is the last section of the route
			* If it is only part of the route, but not the very end, the Middleware does not get called
		* The third parameter is a Namer
			* The Namer is used to figure out what the proper route to get to this router is
				* It takes in the vars, and uses it to return a string that corresponds to the branch string that it corresponds to
	* You can call router.New() to create a blank custom route
		* You will then have to fill in Routing, Action, and Name similarly to the way you would pass in the parameters to the router.Route function
	* You can create your own type that uses one of the previous methods to create a route
		* It will need to implement the router.HasRouter interface
			* It needs to return a router on command - Router()
			* It needs to be able to store it's parent - SetParent()
			* It needs to be able to return it's name - Route()
* Add child routes to parent ones with the AddRoute() method
* Once you have one base route, and the rest are descendants of that route, add it to your rack

## 	Example

	package main

	import (
		"github.com/ScruffyProdigy/Middleware/router"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
		"strings"
	)

	var coins = map[string]string{"penny": "useless", "nickel": "heavy and annoying", "dime": "light and annoying", "quarter": "not obsolete quite yet"}

	var RootWare rack.Func = func(vars map[string]interface{}, next func()) {
		(httper.V)(vars).SetMessageString("<html>Check out <a href='/coins'>My Coins</a></html>")
	}

	var CoinCollectionWare rack.Func = func(vars map[string]interface{}, next func()) {
		v := httper.V(vars)
		v.SetMessageString("<html><ul>")
		coinnames := []string{}
		for coin, _ := range coins {
			coinnames = append(coinnames, coin)
		}
		sort.Strings(coinnames)
		for _, coin := range coinnames {
			v.AppendMessageString("<li><a href='/coins/" + coin + "'>")
			v.AppendMessageString(strings.ToUpper(coin[:1]) + coin[1:])
			v.AppendMessageString("</a></li>")
		}
		v.AppendMessageString("</ul></html>")
	}

	func init() {
		var MemberRoute *Router = router.New()

		MemberRoute.Routing = SignalFunc(func(vars map[string]interface{}) bool {
			coinName := (V)(vars).CurrentSection()
			coinInfo, exists := coins[coinName]
			if !exists {
				return false
			}
			vars["Name"] = coinName
			vars["Info"] = coinInfo
			return true
		})

		MemberRoute.Action = rack.Func(func(vars map[string]interface{}, next func()) {
			name := vars["Name"].(string)
			info := vars["Info"].(string)
			(httper.V)(vars).SetMessageString(name + " - " + info + " - " + MemberRoute.Route(vars))
		})

		MemberRoute.Name = NamerFunc(func(vars map[string]interface{}, prev func() string) string {
			name := "(coin)"
			if coin, ok := vars["Name"].(string); ok {
				name = coin
			}
			return prev() + name + "/"
		})

		CollectionRoute := router.BasicRoute("coins", CoinCollectionWare)
		CollectionRoute.AddRoute(MemberRoute)

		Root := router.BasicRoute("", RootWare)
		Root.AddRoute(CollectionRoute)

		conn := httper.HttpConnection(":4011")
		go conn.Go(Root)
	}
	

running this provides appropriate pages at localhost:3000 for "/","/coins","/coins/penny","/coins/nickel","/coins/dime","/coins/quarter" (also, since router defaults to being not case sensitive, the same routes will work regardless of how you capitalize them)
