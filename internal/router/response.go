package router

import (
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
	"net/http"
)

type Response struct {
	W              http.ResponseWriter
	req            *Request
	renderer       *templates.renderer
	hasWrittenData bool
	router         *Router

	status  int
	headers *http.Header
}

func newResponse(
	w http.ResponseWriter,
	req *Request,
	renderer *templates.renderer,
	router *Router,
) *Response {
	return &Response{
		W:        w,
		req:      req,
		renderer: renderer,
		router:   router,
	}
}

func (res *Response) Cookie(cookie *http.Cookie) interfaces.Response {
	http.SetCookie(res.W, cookie)
	return res
}

func (res *Response) Status(status int) interfaces.Response {
	// In a regular Go server, headers cannot be added after calling the
	// WriteHeader(statusCode) method. To make it more user-friendly, we cache
	// status and write it lazily when the body is written or the response is
	// ended.
	res.status = status
	return res
}

func (res *Response) Header(name, value string) interfaces.Response {
	res.W.Header().Set(name, value)
	return res
}

func (res *Response) View(name, layout string, data interface{}) interfaces.Result {
	return res.newResult().View(name, layout, data)
}

func (res *Response) Data(data interface{}) interfaces.Result {
	return res.newResult().Data(data)
}

func (res *Response) Error(err error) interfaces.Result {
	return res.newResult().Error(err)
}

func (res *Response) Sprintf(format string, a ...interface{}) interfaces.Result {
	return res.newResult().Sprintf(format, a...)
}

func (res *Response) Redirect(uri string, code int) interfaces.Result {
	return res.newResult().Redirect(uri, code)
}

func (res *Response) RedirectRoute(name string, params map[string]string, code int) interfaces.Result {
	return res.newResult().RedirectRoute(name, params, code)
}

func (res *Response) End() interfaces.Result {
	return res.newResult()
}

func (res *Response) newResult() interfaces.Result {
	return newResult(res.W, res.req.R, res.status, res.renderer, res.router)
}
