package request

import (
	"github.com/gschier/hemlock/interfaces"
	"io"
	"net/http"
)

type Request struct {
	r *http.Request
}

func New(r *http.Request) interfaces.Request {
	return &Request{r: r}
}

// Input grabs input from the request body by name
func (req *Request) Input(name string) string {
	panic("Implement me")
}

// Input grabs input from the query string by name
func (req *Request) Query(name string) string {
	return req.r.URL.Query().Get(name)
}

// Cookie grabs input from cookies by name
func (req *Request) Cookie(name string) string {
	panic("Implement me")
}

// Cookie grabs file name
func (req *Request) File(name string) io.Reader {
	panic("Implement me")
}

