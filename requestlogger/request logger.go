/*
	requestLogger will log all requests to your web server
*/
package requestlogger

import (
	"github.com/ScruffyProdigy/Middleware/logger"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
)

//Logger is a Middleware that will log all requests to your web server
var Logger rack.Func = func(vars map[string]interface{}, next func()) {
	r := (httper.V)(vars).GetRequest()
	(logger.V)(vars).Println(r.Method, r.URL.String())
	next()
}
