package encapsulator

import (
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"github.com/ScruffyProdigy/TheTemplater/templater"
	"io/ioutil"
	"net/http"
	"fmt"
)

var templates = templater.New("./test_templates")

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


func Example_Basic() {
	rackup := rack.New()
	rackup.Add(AddLayout(templates))
	rackup.Add(rack.Func(func(vars map[string]interface{}, next func()) {
		vars["Layout"] = "test"
		vars["Title"] = "Hello World"
		(httper.V)(vars).AppendMessageString("Hello World!")
	}))

	conn := httper.HttpConnection(":3000")
	go conn.Go(rackup)
	
	GetFrom("http://localhost:3000/")
	//output: <html><head><title>Hello World</title></head><body>Hello World!</body></html>
}

func Example_NoLayout() {
	rackup := rack.New()
	rackup.Add(AddLayout(templates))
	rackup.Add(rack.Func(func(vars map[string]interface{}, next func()) {
		vars["Title"] = "Hello World"
		(httper.V)(vars).AppendMessageString("Hello World!")
	}))
	
	conn := httper.HttpConnection(":3001")
	go conn.Go(rackup)
	
	GetFrom("http://localhost:3001/")
	//output: Hello World!
}

func Example_BadLayout() {
	rackup := rack.New()
	rackup.Add(AddLayout(templates))
	rackup.Add(rack.Func(func(vars map[string]interface{}, next func()) {
		vars["Layout"] = "invalid"
		vars["Title"] = "Hello World"
		(httper.V)(vars).AppendMessageString("Hello World!")
	}))
	
	conn := httper.HttpConnection(":3001")
	go conn.Go(rackup)
	
	GetFrom("http://localhost:3001/")
	//output: Hello World!
}