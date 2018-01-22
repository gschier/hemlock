package router

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Request struct {
	r *http.Request
}

func newRequest(r *http.Request) *Request {
	return &Request{r: r}
}

func (req *Request) URL() *url.URL {
	return req.r.URL
}

func (req *Request) Query(name string) string {
	return req.r.URL.Query().Get(name)
}

func (req *Request) QueryInt(name string) int {
	value, err := strconv.Atoi(req.Query(name))
	if err != nil {
		return 0
	}
	return value
}

func (req *Request) Cookie(name string) string {
	panic("Implement me")
}

func (req *Request) File(name string) io.Reader {
	panic("Implement me")
}

func (req *Request) Context() context.Context {
	return req.r.Context()
}
