package varser

type V map[string]interface{}

func (this V) Run(vars map[string]interface{}, next func()) {
	for k, v := range this {
		vars[k] = v
	}
	next()
}
