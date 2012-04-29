#	Routes
This provides a branching based router that goes through every directory specified in the requested url and finds a subrouter until it gets to the end, then calls whichever middleware is associated with the last subrouter found

## 	Dependencies
None

## 	Installation
`go get github.com/HairyMezican/Middleware/routes`

## 	Example

	package main

	import (
		"github.com/HairyMezican/Middleware/router"
		"github.com/HairyMezican/TheRack/rack"
		"net/http"
		"strings"
	)

	var coins = map[string]string{"penny": "useless", "nickel": "heavy and annoying", "dime": "light and annoying", "quarter": "not obsolete quite yet"}

	var RootWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		return http.StatusOK, rack.NewHeader(), []byte("<html>Check out <a href='/coins'>My Coins</a></html>")
	}

	var CoinCollectionWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		message := "<html><ul>"
		for coin, _ := range coins {
			message += "<li><a href='/coins/" + coin + "'>" + strings.ToUpper(coin[:1]) + coin[1:] + "</a></li>"
		}
		message += "</ul></html>"
		return http.StatusOK, rack.NewHeader(), []byte(message)
	}

	var CoinSignaler router.SignalFunc = func(r *http.Request, vars rack.Vars) bool {
		coinName := vars.Apply(router.CurrentSection).(string)
		coinInfo, exists := coins[coinName]
		if !exists {
			return false
		}
		vars["Info"] = []byte(coinInfo)
		return true
	}

	var CoinMemberWare rack.Func = func(r *http.Request, vars rack.Vars, next rack.Next) (int, http.Header, []byte) {
		return http.StatusOK, rack.NewHeader(), vars["Info"].([]byte)
	}

	func main() {
		CollectionRoute := router.BasicRoute("coins", CoinCollectionWare)
		MemberRoute := router.Route(CoinSignaler, CoinMemberWare)
		router.Root.Action = RootWare
		router.Root.AddRoute(CollectionRoute)
		CollectionRoute.AddRoute(MemberRoute)
		conn := rack.HttpConnection(":3000")
		rack.Up.Add(router.Parser)
		rack.Up.Add(router.Root)
		rack.Run(conn, rack.Up)
	}
	

running this provides appropriate pages at localhost:3000 for "/","/coins","/coins/penny","/coins/nickel","/coins/dime","/coins/quarter" (also, since router is not case sensitive, the same routes will work regardless of how you capitalize them)
