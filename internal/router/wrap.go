package router

import (
	"github.com/gschier/hemlock/interfaces"
	"net/http"
)

type mHandler struct {
	CalledNext bool
	next       interfaces.Next
	req        interfaces.Request
	res        interfaces.Response
}

func (h *mHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Update request/response because they may be different
	h.req.(*Request).R = r
	h.res.(*Response).W = w
	h.CalledNext = true
}

func wrapMiddleware(h func(next http.Handler) http.Handler) interfaces.Middleware {
	return func(req interfaces.Request, res interfaces.Response, next interfaces.Next) interfaces.Result {
		handler := &mHandler{next: next, req: req, res: res}
		h(handler).ServeHTTP(res.(*Response).W, req.(*Request).R)

		if handler.CalledNext {
			return next(req, res)
		}

		// Middleware returned a result, so end it!
		return res.End()
	}
}

func nextMiddleware(middlewares []interfaces.Middleware, i int, req interfaces.Request, res interfaces.Response) interfaces.Result {
	if i == len(middlewares) {
		return nil
	}

	next := func(newReq interfaces.Request, newRes interfaces.Response) interfaces.Result {
		return nextMiddleware(middlewares, i+1, newReq, newRes)
	}

	return middlewares[i](req, res, next)
}
