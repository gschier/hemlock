package router

import (
	"github.com/gorilla/mux"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
	"net/http"
)

type Route struct {
	route  *mux.Route
	router *Router
}

func NewRoute(router *Router, route *mux.Route) *Route {
	return &Route{router: router, route: route}
}

func (r *Route) Redirect(uri, to string, code int) interfaces.Route {
	return r.assignCallback(nil, uri, func(res interfaces.Response) interfaces.Result {
		return res.Redirect(to, code)
	})
}

func (r *Route) View(uri, view, layout string, data interface{}) interfaces.Route {
	return r.Get(uri, func(res interfaces.Response) interfaces.Result {
		return res.View(view, layout, data)
	})
}

func (r *Route) Callback(callback interfaces.Callback) interfaces.Route {
	r.route.HandlerFunc(r.wrap(callback))
	return r
}

func (r *Route) Get(uri string, callback interfaces.Callback) interfaces.Route {
	return r.assignCallback([]string{http.MethodGet}, uri, callback)
}

func (r *Route) Post(uri string, callback interfaces.Callback) interfaces.Route {
	return r.assignCallback([]string{http.MethodPost}, uri, callback)
}

func (r *Route) Put(uri string, callback interfaces.Callback) interfaces.Route {
	return r.assignCallback([]string{http.MethodPut}, uri, callback)
}

func (r *Route) Patch(uri string, callback interfaces.Callback) interfaces.Route {
	return r.assignCallback([]string{http.MethodPatch}, uri, callback)
}

func (r *Route) Delete(uri string, callback interfaces.Callback) interfaces.Route {
	return r.assignCallback([]string{http.MethodDelete}, uri, callback)
}

func (r *Route) Options(uri string, callback interfaces.Callback) interfaces.Route {
	return r.assignCallback([]string{http.MethodOptions}, uri, callback)
}

func (r *Route) Head(uri string, callback interfaces.Callback) interfaces.Route {
	return r.assignCallback([]string{http.MethodHead}, uri, callback)
}

func (r *Route) Connect(uri string, callback interfaces.Callback) interfaces.Route {
	return r.assignCallback([]string{http.MethodConnect}, uri, callback)
}

func (r *Route) Trace(uri string, callback interfaces.Callback) interfaces.Route {
	return r.assignCallback([]string{http.MethodTrace}, uri, callback)
}

func (r *Route) Any(uri string, callback interfaces.Callback) interfaces.Route {
	return r.assignCallback(nil, uri, callback)
}

func (r *Route) Match(methods []string, uri string, callback interfaces.Callback) interfaces.Route {
	return r.assignCallback(methods, uri, callback)
}

func (r *Route) Methods(methods ...string) interfaces.Route {
	r.route.Methods(methods...)
	return r
}

func (r *Route) Host(uri string) interfaces.Route {
	r.route.Host(uri)
	return r
}

func (r *Route) Prefix(uri string) interfaces.Route {
	r.route.PathPrefix(uri)
	return r
}

func (r *Route) Name(n string) interfaces.Route {
	r.route.Name(n)
	return r
}

func (r *Route) Use(m ...interfaces.Middleware) {
	r.router.Use(m...)
}

func (r *Route) UseG(m ...func(http.Handler) http.Handler) {
	r.router.UseG(m...)
}

func (r *Route) With(m ...interfaces.Middleware) interfaces.Route {
	return r.router.With(m...)
}

func (r *Route) WithG(m ...func(http.Handler) http.Handler) interfaces.Route {
	return r.router.WithG(m...)
}

func (r *Route) Group(fn func(router interfaces.Router)) {
	fn(r.router.fork())
}

func (r *Route) wrap(callback interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r2 *http.Request) {
		var renderer templates.Renderer
		r.router.app.Resolve(&renderer)

		req := newRequest(r2)
		res := newResponse(w, req, &renderer, r.router)

		newApp := hemlock.CloneApplication(r.router.app)
		newApp.Instance(req)
		newApp.Instance(res)

		extraArgs := make([]interface{}, 0)
		for _, v := range mux.Vars(r2) {
			extraArgs = append(extraArgs, v)
		}

		results := newApp.ResolveInto(callback, extraArgs...)
		if len(results) != 1 {
			panic("Route did not return a value. Got " + string(len(results)))
		}
	}
}

func (r *Route) assignCallback(methods []string, uri string, callback interface{}) interfaces.Route {
	if methods != nil {
		r.Methods(methods...)
	}
	r.route.Path(uri).HandlerFunc(r.wrap(callback))
	return r
}
