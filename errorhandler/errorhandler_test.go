package errorhandler

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

func Example_BasicError() {
	rackup := rack.New()
	rackup.Add(ErrorHandler)
	rackup.Add(rack.Func(func(vars map[string]interface{}, next func()) {
		httper.V(vars).SetMessageString("Just Fine!")
		array := make([]byte, 0)
		array[1] = 0 //this action results in a runtime error; we are indexing past the range of the slice
	}))

	conn := httper.HttpConnection(":3000")
	go conn.Go(rackup)

	GetFrom("http://localhost:3000/")
	//output: runtime error: index out of range
}

func Example_StringError() {
	rackup := rack.New()
	rackup.Add(ErrorHandler)
	rackup.Add(rack.Func(func(vars map[string]interface{}, next func()) {
		httper.V(vars).SetMessageString("Just Fine!")
		panic("Error!")
	}))

	conn := httper.HttpConnection(":3001")
	go conn.Go(rackup)

	GetFrom("http://localhost:3001/")
	//output: Error!
}

func Example_NoError() {
	rackup := rack.New()
	rackup.Add(ErrorHandler)
	rackup.Add(rack.Func(func(vars map[string]interface{}, next func()) {
		httper.V(vars).SetMessageString("Just Fine!")
	}))

	conn := httper.HttpConnection(":3002")
	go conn.Go(rackup)

	GetFrom("http://localhost:3002/")
	//output: Just Fine!
}
