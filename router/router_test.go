package router

import (
	"fmt"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

func GetFrom(loc string) {
	resp, err := http.Get(loc)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if resp.StatusCode != 200 {
		fmt.Print(resp.StatusCode)
	}

	fmt.Println(string(body))
}

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

var CoinSignaler SignalFunc = func(vars map[string]interface{}) bool {
	coinName := (V)(vars).CurrentSection()
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

func init() {
	MemberRoute := Route(CoinSignaler, CoinMemberWare)

	CollectionRoute := BasicRoute("coins", CoinCollectionWare)
	CollectionRoute.AddRoute(MemberRoute)

	Root := BasicRoute("", RootWare)
	Root.AddRoute(CollectionRoute)

	conn := httper.HttpConnection(":4011")
	go conn.Go(Root)
}

func Example_Root() {
	GetFrom("http://localhost:4011")

	//output: <html>Check out <a href='/coins'>My Coins</a></html>
}

func Example_CoinCollection() {
	GetFrom("http://localhost:4011/coins")

	//output: <html><ul><li><a href='/coins/dime'>Dime</a></li><li><a href='/coins/nickel'>Nickel</a></li><li><a href='/coins/penny'>Penny</a></li><li><a href='/coins/quarter'>Quarter</a></li></ul></html>
}

func Example_CoinMembers() {
	GetFrom("http://localhost:4011/coins/penny")
	GetFrom("http://localhost:4011/coins/nickel")
	GetFrom("http://localhost:4011/coins/dime")
	GetFrom("http://localhost:4011/coins/quarter")


//Nota Buena: go fmt messes this next section up - it puts tabs in, which then makes the output incorrect

/* output:
useless
heavy and annoying
light and annoying
not obsolete quite yet
*/
}

func Example_Missing() {
	GetFrom("http://localhost:4011/coins/halfdollar")
	//output: 404
}

func Example_CoinCollectionExtraSlash() {
	GetFrom("http://localhost:4011/coins/")

	//output: <html><ul><li><a href='/coins/dime'>Dime</a></li><li><a href='/coins/nickel'>Nickel</a></li><li><a href='/coins/penny'>Penny</a></li><li><a href='/coins/quarter'>Quarter</a></li></ul></html>
}

func Example_CoinCollectionMiscapitalized() {
	GetFrom("http://localhost:4011/CoInS")

	//output: <html><ul><li><a href='/coins/dime'>Dime</a></li><li><a href='/coins/nickel'>Nickel</a></li><li><a href='/coins/penny'>Penny</a></li><li><a href='/coins/quarter'>Quarter</a></li></ul></html>
}
