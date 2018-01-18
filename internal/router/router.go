package router

import (
	"github.com/go-chi/chi"
	"github.com/gschier/hemlock"
	"net/http"
)

type router struct {
	app  *hemlock.Application
	root chi.Router
}

func NewRouter(app *hemlock.Application) *router {
	return &router{root: chi.NewRouter(), app: app}
}

func (router *router) Redirect(uri, to string, code int) {
	router.root.HandleFunc(uri, func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, to, code)
	})
}

func (router *router) Get(uri string, callback interface{}) {
	router.addRoute(http.MethodGet, uri, callback)
}

func (router *router) Post(uri string, callback interface{}) {
	router.addRoute(http.MethodPost, uri, callback)
}

func (router *router) Put(uri string, callback interface{}) {
	router.addRoute(http.MethodPut, uri, callback)
}

func (router *router) Patch(uri string, callback interface{}) {
	router.addRoute(http.MethodPatch, uri, callback)
}

func (router *router) Delete(uri string, callback interface{}) {
	router.addRoute(http.MethodDelete, uri, callback)
}

func (router *router) Options(uri string, callback interface{}) {
	router.addRoute(http.MethodOptions, uri, callback)
}

// Handler returns the HTTP handler
func (router *router) Handler() http.Handler {
	return router.root
}

func (router *router) addRoute(method string, pattern string, callback interface{}) {
	router.root.MethodFunc(method, pattern, func(w http.ResponseWriter, r *http.Request) {
		newApp := hemlock.CloneApplication(router.app)
		newApp.Instance(newRequest(r))
		newApp.Instance(newResponse(w))

		c := chi.RouteContext(r.Context())
		extraArgs := make([]interface{}, len(c.URLParams.Values))
		for i, v := range c.URLParams.Values {
			extraArgs[i] = v
		}

		results := newApp.ResolveInto(callback, extraArgs...)
		if len(results) != 1 {
			panic("Route did not return a value. Got " + string(len(results)))
		}

		view, ok := results[0].(*View)
		if !ok {
			panic("Route did not return View instance")
		}

		w.WriteHeader(view.Status)
		w.Write(view.Bytes)
	})
}
