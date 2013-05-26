package encapsulator

import (
	"fmt"
	"github.com/ScruffyProdigy/Middleware/templater"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"io/ioutil"
	"net/http"
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

	fmt.Println(string(body))
}

func Example_Basic() {
	rackup := rack.New()
	rackup.Add(templater.GetTemplates("./test_templates"))
	rackup.Add(AddLayout)
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
	rackup.Add(templater.GetTemplates("./test_templates"))
	rackup.Add(AddLayout)
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
	rackup.Add(templater.GetTemplates("./test_templates"))
	rackup.Add(AddLayout)
	rackup.Add(rack.Func(func(vars map[string]interface{}, next func()) {
		vars["Layout"] = "invalid"
		vars["Title"] = "Hello World"
		(httper.V)(vars).AppendMessageString("Hello World!")
	}))

	conn := httper.HttpConnection(":3002")
	go conn.Go(rackup)

	GetFrom("http://localhost:3002/")
	//output: Hello World!
}
