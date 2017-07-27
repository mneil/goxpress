package goxpress

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// Get returns a parameter by name
func (p Params) Get(name string) string {
	for i := range p {
		if p[i].Key == name {
			return p[i].Value
		}
	}
	return ""
}
