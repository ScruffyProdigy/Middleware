package methoder

import (
	"fmt"
	"github.com/ScruffyProdigy/Middleware/parser"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"io/ioutil"
	"net/http"
	"net/url"
)

var HttpWare rack.Func = func(vars map[string]interface{}, next func()) {
	v := (httper.V)(vars)
	request := v.GetRequest()
	v.SetMessageString("You used " + request.Method)
}

func PostTo(loc string, vals url.Values) {
	resp, err := http.PostForm(loc, vals)
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

	fmt.Println(string(body))
}

func init() {
	rackup := rack.New()
	rackup.Add(parser.Form)
	rackup.Add(Override)
	rackup.Add(HttpWare)

	conn := httper.HttpConnection(":3000")
	go conn.Go(rackup)
}

func Example_Default() {
	PostTo("http://localhost:3000", url.Values{})
	//output: You used POST
}

func Example_Post() {
	PostTo("http://localhost:3000", url.Values{"_method": {"post"}})
	//output: You used POST
}

func Example_Get() {
	PostTo("http://localhost:3000", url.Values{"_method": {"get"}})
	//output: You used GET
}

func Example_Delete() {
	PostTo("http://localhost:3000", url.Values{"_method": {"delete"}})
	//output: You used DELETE
}

func Example_Put() {
	PostTo("http://localhost:3000", url.Values{"_method": {"put"}})
	//output: You used PUT
}

func Example_Invalid() {
	PostTo("http://localhost:3000", url.Values{"_method": {"invalid"}})
	//output: You used POST
}
