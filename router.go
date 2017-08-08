package goxpress

import (
	"fmt"
	"net/http"
	"time"
)

// Handle defines the signature of a request handler for a route
type Handle func(http.ResponseWriter, *http.Request, *Context)

// Router impliments ServeMux and will allow us to attach route handlers
type Router struct {
	trees map[string]*node
	Context
}

// Listen starts a web Router running on port
func (r *Router) Listen(port int) (*http.Server, error) {

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        r,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s, s.ListenAndServe()
}

// ServeHTTP is the main Router implimentation. Runs all middleware then does route matching
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	ctx := NewContext()

	path := req.URL.Path
	var handle Handle
	var ps Params

	if root := r.trees[req.Method]; root != nil {
		if handle, ps, _ = root.getValue(path); handle != nil {
			ctx.Params = ps
		}
	}

	// Iterate over available middlewares
	if len(middlewares) != 0 {
		for _, middleware := range middlewares {
			err := middleware(w, req, ctx)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}
	}

	// Handle the route if it exists
	if handle != nil {
		handle(w, req, ctx)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

// Handler helps implement the ServeMux interface. Not recommended for use however because it
// bypasses the passing of request context which is the entire purpose of this library.
// But this may satisfy edge cases where context is not wanted
func (r *Router) Handler(method, path string, handler http.Handler) {
	r.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request, _ *Context) {
			handler.ServeHTTP(w, req)
		},
	)
}

// HandleFunc satisfies the ServeMux interface. Not recommended for use but may
// satisfy edge cases
func (r *Router) HandleFunc(pattern string, handler http.Handler) {
	r.Handler("*", pattern, handler)
}

// Handle is a generic method for adding a request handle on a pattern route and method
func (r *Router) Handle(method, path string, handle Handle) {
	if path[0] != '/' {
		panic(fmt.Sprintf("path must begin with '/' in path '%s'", path))
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}
	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}
	root.addRoute(path, handle)
}

// Use middleware to route request/response cycle
func (*Router) Use(middleware func(w http.ResponseWriter, r *http.Request, c *Context) error) {
	middlewares = append(middlewares, middleware)
}

// middlewares keeps track of middleware that was added
var middlewares = make([]func(w http.ResponseWriter, r *http.Request, c *Context) error, 0, 5)

// GET registers a path to a GET request
func (r *Router) GET(path string, handle Handle) {
	r.Handle("GET", path, handle)
}

// HEAD registers a path to a HEAD request
func (r *Router) HEAD(path string, handle Handle) {
	r.Handle("HEAD", path, handle)
}

// POST registers a path to a POST request
func (r *Router) POST(path string, handle Handle) {
	r.Handle("POST", path, handle)
}

// PUT registers a path to a PUT request
func (r *Router) PUT(path string, handle Handle) {
	r.Handle("PUT", path, handle)
}

// DELETE registers a path to DELETE request
func (r *Router) DELETE(path string, handle Handle) {
	r.Handle("DELETE", path, handle)
}

// CONNECT registers a path to CONNECT request
func (r *Router) CONNECT(path string, handle Handle) {
	r.Handle("CONNECT", path, handle)
}

// OPTIONS registers a path to a OPTIONS request
func (r *Router) OPTIONS(path string, handle Handle) {
	r.Handle("OPTIONS", path, handle)
}

// PATCH registers a path to a PATCH request
func (r *Router) PATCH(path string, handle Handle) {
	r.Handle("PATCH", path, handle)
}
