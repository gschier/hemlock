package hemlock

import (
	"github.com/go-chi/chi"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/request"
	"github.com/gschier/hemlock/internal/response"
	"net/http"
)

type Router struct {
	app  *Application
	root chi.Router
}

func NewRouter(app *Application) interfaces.Router {
	return &Router{root: chi.NewRouter(), app: app}
}

func (router *Router) Redirect(uri, to string, code int) {
	router.root.HandleFunc(uri, func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, to, code)
	})
}

func (router *Router) Get(uri string, callback interfaces.Callback) {
	router.addRoute(http.MethodGet, uri, callback)
}

func (router *Router) Post(uri string, callback interfaces.Callback) {
	router.addRoute(http.MethodPost, uri, callback)
}

func (router *Router) Put(uri string, callback interfaces.Callback) {
	router.addRoute(http.MethodPut, uri, callback)
}

func (router *Router) Patch(uri string, callback interfaces.Callback) {
	router.addRoute(http.MethodPatch, uri, callback)
}

func (router *Router) Delete(uri string, callback interfaces.Callback) {
	router.addRoute(http.MethodDelete, uri, callback)
}

func (router *Router) Options(uri string, callback interfaces.Callback) {
	router.addRoute(http.MethodOptions, uri, callback)
}

// Handler returns the HTTP handler
func (router *Router) Handler() http.Handler {
	return router.root
}

func (router *Router) addRoute(method string, pattern string, callback interfaces.Callback) {
	router.root.MethodFunc(method, pattern, func(w http.ResponseWriter, r *http.Request) {
		//newApp := CloneApplication(router.app)
		newApp := router.app
		newApp.Instance(request.New(r))
		newApp.Instance(response.New(w))

		c := chi.RouteContext(r.Context())
		extraArgs := make([]interface{}, len(c.URLParams.Values))
		for i, v := range c.URLParams.Values {
			extraArgs[i] = v
		}

		results := newApp.ResolveInto(callback, extraArgs...)
		if len(results) != 1 {
			panic("Route did not return a value. Got " + string(len(results)))
		}

		view, ok := results[0].(*response.View)
		if !ok {
			panic("Route did not return View instance")
		}

		w.WriteHeader(view.Status)
		w.Write(view.Bytes)
	})
}

// Should be a function
type Callback interface{}

