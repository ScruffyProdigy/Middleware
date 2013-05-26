package redirecter

import (
	"fmt"
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

func Example_Redirect() {
	rackup := rack.New()
	rackup.Add(rack.Func(func(vars map[string]interface{}, next func()) {
		h := (httper.V)(vars)
		path := h.GetRequest().URL.Path
		if path != "/" {
			h.SetMessageString(path)
		} else {
			next()
		}
	}))
	rackup.Add(Redirecter{"/redirected/"})

	conn := httper.HttpConnection(":3000")
	go conn.Go(rackup)

	GetFrom("http://localhost:3000/")
	//output: /redirected/
}
