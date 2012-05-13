#	Routes
This provides a branching based router that goes through every directory specified in the requested url and finds a subrouter until it gets to the end, then calls whichever middleware is associated with the last subrouter found

## 	Installation
`go get github.com/HairyMezican/Middleware/routes`

## 	Example

	package main

	import (
		"github.com/HairyMezican/Middleware/router"
		"github.com/HairyMezican/TheRack/rack"
		"strings"
	)

	var coins = map[string]string{"penny": "useless", "nickel": "heavy and annoying", "dime": "light and annoying", "quarter": "not obsolete quite yet"}

	var RootWare rack.Func = func(vars rack.Vars, next func()) {
		rack.SetMessageString(vars, "<html>Check out <a href='/coins'>My Coins</a></html>")
	}

	var CoinCollectionWare rack.Func = func(vars rack.Vars, next func()) {
		rack.SetMessageString(vars, "<html><ul>")
		for coin, _ := range coins {
			rack.AppendMessageString(vars, "<li><a href='/coins/"+coin+"'>")
			rack.AppendMessageString(vars, strings.ToUpper(coin[:1])+coin[1:])
			rack.AppendMessageString(vars, "</a></li>")
		}
		rack.AppendMessageString(vars, "</ul></html>")
	}

	var CoinSignaler router.SignalFunc = func(vars rack.Vars) bool {
		coinName := router.CurrentSection(vars)
		coinInfo, exists := coins[coinName]
		if !exists {
			return false
		}
		vars["Info"] = []byte(coinInfo)
		return true
	}

	var CoinMemberWare rack.Func = func(vars rack.Vars, next func()) {
		rack.SetMessage(vars, vars["Info"].([]byte))
	}

	func main() {
		MemberRoute := router.Route(CoinSignaler, CoinMemberWare)

		CollectionRoute := router.BasicRoute("coins", CoinCollectionWare)
		CollectionRoute.AddRoute(MemberRoute)

		Root := router.BasicRoute("", RootWare)
		Root.AddRoute(CollectionRoute)

		conn := rack.HttpConnection(":3000")
		conn.Go(Root)
	}
	

running this provides appropriate pages at localhost:3000 for "/","/coins","/coins/penny","/coins/nickel","/coins/dime","/coins/quarter" (also, since router defaults to being not case sensitive, the same routes will work regardless of how you capitalize them)
