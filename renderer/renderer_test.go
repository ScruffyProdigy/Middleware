package renderer

import (
	"fmt"
	"github.com/ScruffyProdigy/Middleware/logger"
	"github.com/ScruffyProdigy/Middleware/templater"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"io/ioutil"
	"net/http"
	"os"
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

func Example_Render() {
	rackup := rack.New()
	rackup.Add(logger.Set(os.Stdout, "", 0))
	rackup.Add(templater.GetTemplates("./test_templates"))
	rackup.Add(rack.Func(func(vars map[string]interface{}, next func()) {
		vars["Object"] = "World"
		next()
	}))
	rackup.Add(Renderer{"test"})

	conn := httper.HttpConnection(":3000")
	go conn.Go(rackup)

	GetFrom("http://localhost:3000/")
	//output: Hello World
}
