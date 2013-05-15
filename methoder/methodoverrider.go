package methoder

import (
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheRack/rack"
	"strings"
)

var legal = map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}

func isLegal(s string) bool {
	return legal[s]
}

var Override rack.Func = func(vars map[string]interface{}, next func()) {
	r := httper.V(vars).GetRequest()
	method := strings.ToUpper(r.Form.Get("_method"))
	if isLegal(method) {
		r.Method = method
	}
	next()
}
