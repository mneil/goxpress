# GoXpress

An [express](https://expressjs.com/) "like" server implementation for go. Being a Node developer for a few years I've grown fond of the standard server and packages built on top of it like express.

This is an active work in progress that will get more Express like as time allows. Currently implemented features are:

 - Router with pattern matching [thanks to httprouter](https://github.com/julienschmidt/httprouter)
 - Middleware used on every request
 - Request specific context object

# Routes

See [httprouter](https://github.com/julienschmidt/httprouter) for pattern matching on routes. While GoXpress uses the same pattern matching
tree from httprouter it is not a direct extension of that project. The pattern matching was brought in.

### Example

```
package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"fmt"

	"github.com/mneil/goxpress"
)

func main() {

	var App = &goxpress.Router{}

	App.Use(ContentType)
	App.Use(GetUser)

	App.GET("/", func(w http.ResponseWriter, r *http.Request, c *goxpress.Context) {
		ctype, _ := c.Get("ctype")
		fmt.Println("The requst content type was", ctype)
	})
	App.GET("/users/:id", func(w http.ResponseWriter, r *http.Request, c *goxpress.Context) {
		u, _ := c.Get("user")
		user := u.(*User)
		fmt.Println(fmt.Sprintf("user id %d is named %s with param %s", user.ID, user.Name, c.Params.Get("id")))
	})

	s, err := App.Listen(4000)
	log.Fatal(err)
	s.ReadTimeout = 30 * time.Second
	s.WriteTimeout = 30 * time.Second
	s.MaxHeaderBytes = 1 << 20
}

// ContentType determines the request content type and stores it in the request context
func ContentType(w http.ResponseWriter, r *http.Request, c *goxpress.Context) (string, error) {
	var ctype = r.Header.Get("Content-Type")
	c.Set("ctype", strings.ToLower(ctype))
	return "", nil
}

// User is an authorized user of the application
type User struct {
	ID   int
	Name string
}

// GetUser stores a mock user in the request context
func GetUser(w http.ResponseWriter, r *http.Request, c *goxpress.Context) (string, error) {
	user := &User{1, "Michael"}
	c.Set("user", user)
	return "", nil
}

```

## Performance

Performance is anecdotal but in case you're wondering here are some benchmarks using the example script above to request the user/:id route.

```
siege -c96 -t60s http://192.168.1.181:4000/users/1
** SIEGE 3.0.5
** Preparing 96 concurrent users for battle.
The server is now under siege...
Lifting the server siege...      done.

Transactions:                  11243 hits
Availability:                 100.00 %
Elapsed time:                  59.26 secs
Data transferred:               0.00 MB
Response time:                  0.01 secs
Transaction rate:             189.72 trans/sec
Throughput:                     0.00 MB/sec
Concurrency:                    1.60
Successful transactions:       11243
Failed transactions:               0
Longest transaction:            0.36
Shortest transaction:           0.00
```

**ROADMAP**

 [x] pattern matching routes  
 [x] route based on http method  
 [x] ability to pass context between requests  
 [] cors settings  
 [] ability to add catch-all route  
 [] cookie management  
 [] global response header management  
 [] simplify response writer interface  
 [] static file server  
