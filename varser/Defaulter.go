/*
	varser is a lightweight variable-setter
*/
package varser

// V is a type you can cast your vars into in order to be a Middleware
type V map[string]interface{}

//Run implements the rack.Middleware interface
func (this V) Run(vars map[string]interface{}, next func()) {
	for k, v := range this {
		vars[k] = v
	}
	next()
}
