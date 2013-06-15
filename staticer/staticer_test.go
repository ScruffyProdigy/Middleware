package staticer

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

func Example_Basic() {
	rackup := rack.New()
	rackup.Add(New("/static/", "./test_files"))

	conn := httper.HttpConnection(":3000")
	go conn.Go(rackup)

	GetFrom("http://localhost:3000/static/test.txt")
	//output: Hello World!
}
