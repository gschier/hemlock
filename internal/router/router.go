package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
	"net/http"
)

type router struct {
	app         *hemlock.Application
	root        chi.Router
	middlewares []interfaces.Middleware
}

func NewRouter(app *hemlock.Application) *router {
	root := chi.NewRouter()

	// TODO: Make middleware configurable
	root.Use(middleware.Recoverer)
	root.Use(middleware.DefaultCompress)
	root.Use(middleware.CloseNotify)
	root.Use(middleware.Logger)
	root.Use(middleware.RedirectSlashes)
	root.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Forwarded-Proto") == "http" {
				newUrl := "https://" + r.Host + r.URL.String()
				http.Redirect(w, r, newUrl, http.StatusFound)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	})

	router := &router{root: root, app: app}
	router.root.NotFound(router.serve(func(req interfaces.Request, res interfaces.Response) interfaces.Result {
		return res.Data("Not Found").Status(404).End()
	}))
	return router
}

func (router *router) Redirect(uri, to string, code int) {
	router.root.HandleFunc(uri, func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, to, code)
	})
}

func (router *router) Get(uri string, callback interfaces.Callback) {
	router.addRoute(http.MethodGet, uri, callback)
}

func (router *router) Post(uri string, callback interfaces.Callback) {
	router.addRoute(http.MethodPost, uri, callback)
}

func (router *router) Put(uri string, callback interfaces.Callback) {
	router.addRoute(http.MethodPut, uri, callback)
}

func (router *router) Patch(uri string, callback interfaces.Callback) {
	router.addRoute(http.MethodPatch, uri, callback)
}

func (router *router) Delete(uri string, callback interfaces.Callback) {
	router.addRoute(http.MethodDelete, uri, callback)
}

func (router *router) Options(uri string, callback interfaces.Callback) {
	router.addRoute(http.MethodOptions, uri, callback)
}

func (router *router) Use(m interfaces.Middleware) {
	router.middlewares = append(router.middlewares, m)
}

// Handler returns the HTTP handler
func (router *router) Handler() http.Handler {
	return router.root
}
func (router *router) callNext(i int, req interfaces.Request, res interfaces.Response) interfaces.Result {
	if i == len(router.middlewares) {
		return nil
	}

	fn := router.middlewares[i]
	view := fn(req, res, func(newReq interfaces.Request, newRes interfaces.Response) interfaces.Result {
		return router.callNext(i+1, newReq, newRes)
	})

	return view
}

func (router *router) addRoute(method string, pattern string, callback interface{}) {
	router.root.MethodFunc(method, pattern, router.serve(callback))
}

func (router *router) serve(callback interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var renderer templates.Renderer
		router.app.Resolve(&renderer)

		req := newRequest(r)
		res := newResponse(w, r, &renderer)

		result := router.callNext(0, req, res)

		// The middleware sent a response so we're done
		if result != nil {
			return
		}

		newApp := hemlock.CloneApplication(router.app)
		newApp.Instance(req)
		newApp.Instance(res)

		c := chi.RouteContext(r.Context())
		extraArgs := make([]interface{}, len(c.URLParams.Values))
		for i, v := range c.URLParams.Values {
			extraArgs[i] = v
		}

		results := newApp.ResolveInto(callback, extraArgs...)
		if len(results) != 1 {
			panic("Route did not return a value. Got " + string(len(results)))
		}

		var ok bool
		result, ok = results[0].(interfaces.Result)
		if !ok {
			panic("Route did not return View instance")
		}
	}
}
