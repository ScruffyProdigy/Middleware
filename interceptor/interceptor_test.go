package interceptor

import (
	"fmt"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"io/ioutil"
	"net/http"
)

var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
	(httper.V)(vars).SetMessageString("Hello World")
}

var RootWare rack.Func = func(vars map[string]interface{}, next func()) {
	(httper.V)(vars).SetMessageString("<html>Check out my <a href=\"helloworld\">Hello World</a></html>")
}

func init() {
	cept := New()
	cept.Intercept("/helloworld", HelloWorldWare)

	rackup := rack.New()
	rackup.Add(cept)
	rackup.Add(RootWare)

	conn := httper.HttpConnection(":3000")
	go conn.Go(rackup)
}

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

	fmt.Println(string(body))
}

func Example_1() {
	GetFrom("http://localhost:3000/")
	//output: <html>Check out my <a href="helloworld">Hello World</a></html>
}

func Example_2() {
	GetFrom("http://localhost:3000/helloworld")
	//output: Hello World
}
