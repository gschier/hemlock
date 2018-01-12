package hemlock

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gschier/hemlock/interfaces"
	"net/http"
)

type Router struct {
	app  *Application
	root chi.Router
}

func NewRouter(app *Application) interfaces.Router {
	return &Router{root: chi.NewRouter(), app: app}
}

func (r *Router) Redirect(uri, to string, code int) {
	r.root.HandleFunc(uri, func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, to, code)
	})
}

func (r *Router) Get(uri string, callback interfaces.Callback) {
	r.addRoute(http.MethodGet, uri, callback)
}

func (r *Router) Post(uri string, callback interfaces.Callback) {
	r.addRoute(http.MethodPost, uri, callback)
}

func (r *Router) Put(uri string, callback interfaces.Callback) {
	r.addRoute(http.MethodPut, uri, callback)
}

func (r *Router) Patch(uri string, callback interfaces.Callback) {
	r.addRoute(http.MethodPatch, uri, callback)
}

func (r *Router) Delete(uri string, callback interfaces.Callback) {
	r.addRoute(http.MethodDelete, uri, callback)
}

func (r *Router) Options(uri string, callback interfaces.Callback) {
	r.addRoute(http.MethodOptions, uri, callback)
}

// Handler returns the HTTP handler
func (r *Router) Handler() http.Handler {
	return r.root
}

func (r *Router) addRoute(method string, uri string, callback interfaces.Callback) {
	fmt.Printf("ADDED ROUTE %v %v\n", method, uri)
	r.root.MethodFunc(method, uri, func(w http.ResponseWriter, req *http.Request) {
		r.app.Instance(w)
		r.app.Instance(req)
		r.app.Instance(NewResponse(w))

		c := chi.RouteContext(req.Context())
		extraArgs := make([]interface{}, len(c.URLParams.Values))
		for i, v := range c.URLParams.Values {
			extraArgs[i] = v
		}

		results := r.app.ResolveInto(callback, extraArgs...)
		if len(results) != 1 {
			panic("Route did not return a value. Got "+string(len(results)))
		}

		view, ok := results[0].(*View)
		if !ok {
			panic("Route did not return View instance")
		}

		w.WriteHeader(view.Status)
		w.Write(view.Bytes)
	})
}

// Should be a function
type Callback interface{}

type Request struct {
}

// Input grabs input from the request body by name
func (req *Request) Input(name string) string {
	return "some input"
}

// Input grabs input from the query string by name
func (req *Request) Query(name string) string {
	panic("Implement me")
}

// Cookie grabs input from cookies by name
func (req *Request) Cookie(name string) string {
	panic("Implement me")
}

// Cookie grabs file name
func (req *Request) File(name string) []byte {
	panic("Implement me")
}
