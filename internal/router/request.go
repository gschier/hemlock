package router

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/gschier/hemlock/interfaces"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Request struct {
	R *http.Request
}

func newRequest(r *http.Request) *Request {
	return &Request{R: r}
}

func (req *Request) URL() *url.URL {
	return req.R.URL
}

func (req *Request) Header(name string) string {
	return req.R.Header.Get(name)
}

func (req *Request) Host() string {
	return req.R.Host
}

func (req *Request) Path() string {
	return req.R.URL.Path
}

func (req *Request) Method() string {
	return req.R.Method
}

func (req *Request) RouteName() string {
	r := mux.CurrentRoute(req.R)
	if r == nil {
		return ""
	}
	return r.GetName()
}

func (req *Request) Query(name string) string {
	return req.R.URL.Query().Get(name)
}

func (req *Request) QueryInt(name string) int {
	value, err := strconv.Atoi(req.Query(name))
	if err != nil {
		return 0
	}
	return value
}

func (req *Request) WithContext(ctx context.Context) interfaces.Request {
	return newRequest(req.R.WithContext(ctx))
}

func (req *Request) Post(name string) string {
	req.R.ParseForm()
	return req.R.Form.Get(name)
}

func (req *Request) Cookie(name string) string {
	panic("Implement me")
}

func (req *Request) File(name string) io.Reader {
	panic("Implement me")
}

func (req *Request) Context() context.Context {
	return req.R.Context()
}
