package statuser

import (
	"fmt"
	"github.com/ScruffyProdigy/Middleware/encapsulator"
	"github.com/ScruffyProdigy/Middleware/statuser"
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

var ErrorWare rack.Func = func(vars map[string]interface{}, next func()) {
	httper.V(vars).StatusError()
}

var ErrorWare2 rack.Func = func(vars map[string]interface{}, next func()) {
	(httper.V)(vars).Status(501)
}

func Example_General() {
	rackup := rack.New()
	rackup.Add(templater.GetTemplates("test_templates"))
	rackup.Add(encapsulator.AddLayout)
	rackup.Add(statuser.SetErrorLayout)
	rackup.Add(ErrorWare)

	conn := httper.HttpConnection(":3000")
	go conn.Go(rackup)

	GetFrom("http://localhost:3000/")

	//output: Error - 500
}

func Example_Specific() {
	rackup := rack.New()
	rackup.Add(templater.GetTemplates("test_templates"))
	rackup.Add(encapsulator.AddLayout)
	rackup.Add(statuser.SetErrorLayout)

	conn := httper.HttpConnection(":3001")
	go conn.Go(rackup)

	GetFrom("http://localhost:3001/")

	//output: Not Found
}

func Example_SpecificOverride() {
	rackup := rack.New()
	rackup.Add(templater.GetTemplates("test_templates"))
	rackup.Add(encapsulator.AddLayout)
	rackup.Add(statuser.SetErrorLayout)
	rackup.Add(ErrorWare2)

	conn := httper.HttpConnection(":3002")
	go conn.Go(rackup)

	GetFrom("http://localhost:3002/")

	//output: Not Implemented!
}
