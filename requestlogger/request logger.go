package requestlogger

import (
	"github.com/ScruffyProdigy/Middleware/logger"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
)

var Logger rack.Func = func(vars map[string]interface{}, next func()) {
	r := (httper.V)(vars).GetRequest()
	(logger.V)(vars).Println(r.Method, r.URL.String())
	next()
}
