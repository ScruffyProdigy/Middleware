package renderer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
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
	rackup.Add((""))
	rackup.Add(Renderer{"test"})

	conn := httper.HttpConnection(":3000")
	go conn.Go(rackup)

	GetFrom("http://localhost:3000/")
	//output: /redirected/
}