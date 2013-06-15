package logger

import (
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"net/http"
	"os"
)

var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
	lg := V(vars).Get()
	if lg != nil {
		lg.Println("Hello World!")
	}
	(httper.V)(vars).SetMessageString("Hello World!")
}

func Example_Basic() {
	rackup := rack.New()
	rackup.Add(Set(os.Stdout, "Log Test - ", 0))
	rackup.Add(HelloWorldWare)

	conn := httper.HttpConnection(":4003")
	go conn.Go(rackup)
	http.Get("http://localhost:4003")
	//output: Log Test - Hello World!
}

func Example_NoLogger() {
	rackup := rack.New()
	rackup.Add(HelloWorldWare)

	conn := httper.HttpConnection(":4004")
	go conn.Go(rackup)
	http.Get("http://localhost:4004")
	//output:
}
