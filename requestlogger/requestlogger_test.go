package requestlogger

import (
	"github.com/ScruffyProdigy/Middleware/logger"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"net/http"
	"os"
)

func Example_Basic() {
	rackup := rack.New()
	rackup.Add(logger.Set(os.Stdout, "", 0))
	rackup.Add(Logger)

	conn := httper.HttpConnection(":3000")
	go conn.Go(rackup)
	http.Get("http://localhost:3000/location")
	//output: GET /location
}
