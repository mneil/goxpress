package goxpress

// Context is an empty interface that allows us to pass any struct through as the context of a request
type Context struct {
	data map[string]interface{}
	Params
}

// Init initializes data so that context can be set
func (r *Context) Init() {
	r.data = make(map[string]interface{})
}

// Get returns the value of interface of some request context
func (r *Context) Get(s string) (interface{}, bool) {
	v, ok := r.data[s]
	return v, ok
}

// Set sets some value by string name during the request context
func (r *Context) Set(s string, i interface{}) {
	r.data[s] = i
}
