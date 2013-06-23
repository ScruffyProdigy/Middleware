package controller

import (
	"bytes"
	"github.com/ScruffyProdigy/Middleware/redirecter"
	"github.com/ScruffyProdigy/Middleware/sessioner"
	"github.com/ScruffyProdigy/Middleware/templater"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"sort"
	"testing"
)

var coins = map[string]string{
	"penny":   "useless",
	"nickel":  "heavy and annoying",
	"dime":    "light and annoying",
	"quarter": "not obsolete quite yet",
}

type Coins struct {
	Heart
}

type CoinModel struct {
	Name, Description string
}

func (this CoinModel) Url() string {
	return "/coins/" + this.Name
}

func (this Coins) Find(name string, vars map[string]interface{}) (interface{}, bool) {
	description, ok := coins[name]
	if ok {
		return &CoinModel{Name: name, Description: description}, true
	}
	return nil, false
}

func (this Coins) Index() {
	coinNames := []string{}
	for coinName, _ := range coins {
		coinNames = append(coinNames, coinName)
	}
	sort.Strings(coinNames)
	this.Set("Coins", coinNames)
}

func (this Coins) Create() {
	name := this.FormValue("Coin[Name]")
	description := this.FormValue("Coin[Description]")

	if name == "" || description == "" {
		this.AddFlash("invalid coin")
		this.RedirectTo("/coins/")
		return
	}

	if _, existing := coins[name]; existing {
		this.AddFlash("coin already created")
		this.RedirectTo("/coins/")
		return
	}

	coins[name] = description

	this.RespondWith(CoinModel{Name: name, Description: description})
}

func (this Coins) Show() {
}

func (this Coins) Update() {
	old := this.GetVar("Coin").(*CoinModel)

	name := this.FormValue("Coin[Name]")
	_, prexisting := coins[name]
	if prexisting && name != old.Name {
		this.AddFlash("coin collision; aborting update")
		this.RedirectTo("/coins")
		return
	}
	if name == "" {
		name = old.Name
	}

	description := this.FormValue("Coin[Description]")
	if description == "" {
		description = old.Description
	}

	delete(coins, old.Name)
	coins[name] = description

	this.RespondWith(CoinModel{Name: name, Description: description})
}

func (this Coins) Destroy() {
	coin := this.GetVar("Coin")
	delete(coins, coin.(*CoinModel).Name)
}

var client *http.Client

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic("No Cookie Jar!")
	}
	client = &http.Client{Jar: jar}
	root := NewRoot(redirecter.Redirecter{"http://localhost:5001/coins"})
	coinRoute := NewResource(&Coins{})
	root.AddRoute(coinRoute)

	rackup := rack.New()
	rackup.Add(templater.GetTemplates("./test_templates"))
	rackup.Add(sessioner.Middleware)
	rackup.Add(root)

	conn := httper.HttpConnection(":5001")
	go conn.Go(rackup)
}

func ResponseFrom(t *testing.T, method, url string, body io.Reader) (string, *http.Response) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err.Error())
	}

	if method == "POST" || method == "PUT" {
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	return string(bytes), res
}

func TextFrom(t *testing.T, method, url string, body io.Reader) string {
	text, _ := ResponseFrom(t, method, url, body)
	return text
}

func StatusFrom(t *testing.T, method, url string, body io.Reader) int {
	_, res := ResponseFrom(t, method, url, body)
	return res.StatusCode
}

func CompareStrings(t *testing.T, a, b string) {
	if a != b {
		t.Error("Should get \"" + b + "\", not \"" + a + "\"")
	}
}

func CompareStatus(t *testing.T, a, b int) {
	if a != b {
		t.Errorf("Status should be %d, not %d\n", b, a)
	}
}

func CompareDestination(t *testing.T, a *http.Response, b string) {
	dest := a.Request.URL.String()
	if dest != b {
		t.Errorf("Final URL should be %s, not %s", b, dest)
	}
}

func Test_Index(t *testing.T) {
	CompareStrings(t, TextFrom(t, "GET", "http://localhost:5001/coins/", nil), "These are my coins: dime nickel penny quarter ")
}

func Test_Show(t *testing.T) {
	CompareStrings(t, TextFrom(t, "GET", "http://localhost:5001/coins/penny", nil), "penny - useless")
	CompareStrings(t, TextFrom(t, "GET", "http://localhost:5001/coins/nickel", nil), "nickel - heavy and annoying")
	CompareStrings(t, TextFrom(t, "GET", "http://localhost:5001/coins/dime", nil), "dime - light and annoying")
	CompareStrings(t, TextFrom(t, "GET", "http://localhost:5001/coins/quarter", nil), "quarter - not obsolete quite yet")
	CompareStatus(t, StatusFrom(t, "GET", "http://localhost:5001/coins/qwutzl", nil), 404)
}

func Test_CreateAndDelete(t *testing.T) {
	text, res := ResponseFrom(t, "POST", "http://localhost:5001/coins/", bytes.NewBufferString("Coin[Name]=half-dollar&Coin[Description]=what%20is%20it%3f"))
	CompareDestination(t, res, "http://localhost:5001/coins/half-dollar")
	CompareStrings(t, text, "half-dollar - what is it?")
	CompareStrings(t, TextFrom(t, "GET", "http://localhost:5001/coins/", nil), "These are my coins: dime half-dollar nickel penny quarter ")

	text, res = ResponseFrom(t, "DELETE", "http://localhost:5001/coins/half-dollar", nil)
	CompareStatus(t, res.StatusCode, 302)
	CompareStrings(t, res.Header.Get("Location"), "/coins")

	CompareStrings(t, TextFrom(t, "GET", "http://localhost:5001/coins/", nil), "These are my coins: dime nickel penny quarter ")
}

func Test_Update(t *testing.T) {
	text, res := ResponseFrom(t, "PUT", "http://localhost:5001/coins/penny", bytes.NewBufferString("Coin[Description]=why%20are%20these%20still%20made%3f"))
	CompareDestination(t, res, "http://localhost:5001/coins/penny")
	CompareStrings(t, res.Request.Method, "GET")
	CompareStrings(t, text, "penny - why are these still made?")

	text, res = ResponseFrom(t, "PUT", "http://localhost:5001/coins/penny", bytes.NewBufferString("Coin[Name]=worthless"))
	CompareDestination(t, res, "http://localhost:5001/coins/worthless")
	CompareStrings(t, res.Request.Method, "GET")
	CompareStrings(t, text, "worthless - why are these still made?")

	text, res = ResponseFrom(t, "PUT", "http://localhost:5001/coins/worthless", bytes.NewBufferString("Coin[Name]=nickel"))
	CompareDestination(t, res, "http://localhost:5001/coins")
	CompareStrings(t, res.Request.Method, "GET")
	CompareStrings(t, text, "coin collision; aborting update These are my coins: dime nickel quarter worthless ")

	text, res = ResponseFrom(t, "PUT", "http://localhost:5001/coins/worthless", bytes.NewBufferString("Coin[Name]=penny&Coin[Description]=useless"))
	CompareDestination(t, res, "http://localhost:5001/coins/penny")
	CompareStrings(t, res.Request.Method, "GET")
	CompareStrings(t, text, "penny - useless")
}
