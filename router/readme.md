#	Routes
This provides a branching based router that goes through every directory specified in the requested url and finds a subrouter until it gets to the end, then calls whichever middleware is associated with the last subrouter found

## 	Installation
`go get github.com/ScruffyProdigy/Middleware/routes`

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
		for coin, _ := range coins {
			v.AppendMessageString("<li><a href='/coins/" + coin + "'>")
			v.AppendMessageString(strings.ToUpper(coin[:1]) + coin[1:])
			v.AppendMessageString("</a></li>")
		}
		v.AppendMessageString("</ul></html>")
	}

	var CoinSignaler router.SignalFunc = func(vars map[string]interface{}) bool {
		coinName := (router.V)(vars).CurrentSection()
		coinInfo, exists := coins[coinName]
		if !exists {
			return false
		}
		vars["Info"] = []byte(coinInfo)
		return true
	}

	var CoinMemberWare rack.Func = func(vars map[string]interface{}, next func()) {
		(httper.V)(vars).SetMessage(vars["Info"].([]byte))
	}

	func main() {
		MemberRoute := router.Route(CoinSignaler, CoinMemberWare)

		CollectionRoute := router.BasicRoute("coins", CoinCollectionWare)
		CollectionRoute.AddRoute(MemberRoute)

		Root := router.BasicRoute("", RootWare)
		Root.AddRoute(CollectionRoute)

		conn := httper.HttpConnection(":3000")
		conn.Go(Root)
	}
	

running this provides appropriate pages at localhost:3000 for "/","/coins","/coins/penny","/coins/nickel","/coins/dime","/coins/quarter" (also, since router defaults to being not case sensitive, the same routes will work regardless of how you capitalize them)
