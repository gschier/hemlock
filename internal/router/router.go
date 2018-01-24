package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Router struct {
	handler     http.Handler
	app         *hemlock.Application
	mux         *mux.Router
	middlewares []interfaces.Middleware
	routes      []*Route
}

func NewRouter(app *hemlock.Application) *Router {
	router := &Router{mux: mux.NewRouter(), app: app}

	//m.Use(middleware.Recoverer)
	//m.Use(middleware.DefaultCompress)
	//m.Use(middleware.CloseNotify)
	//m.Use(middleware.RedirectSlashes)

	if app.Config.Env == "development" {
		router.Use(wrapMiddleware(func(next http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, next)
		}))
	}

	if app.Config.Env == "production" {
		router.Use(wrapMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ext := filepath.Ext(r.URL.Path)
				if ext == ".css" || ext == ".js" {
					w.Header().Add("Cache-Control", "public, max-age=7200")
				}
				next.ServeHTTP(w, r)
			})
		}))
		router.Use(wrapMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-Forwarded-Proto") == "http" {
					newUrl := "https://" + r.Host + r.URL.String()
					http.Redirect(w, r, newUrl, http.StatusFound)
				} else {
					next.ServeHTTP(w, r)
				}
			})
		}))
	}

	router.Use(wrapMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			cwd, _ := os.Getwd()
			fullPath := filepath.Join(cwd, app.Config.PublicDirectory, path)
			s, err := os.Stat(fullPath)
			if err != nil || s.IsDir() {
				next.ServeHTTP(w, r)
				return
			}

			http.ServeFile(w, r, fullPath)
		})
	}))

	router.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var renderer templates.Renderer
		router.app.Resolve(&renderer)
		req := newRequest(r)
		res := newResponse(w, r, &renderer)
		result := nextMiddleware(router.middlewares, 0, req, res)
		if result == nil {
			router.mux.ServeHTTP(w, r)
		}
	})

	return router
}

func (router *Router) Redirect(uri, to string, code int) interfaces.Route {
	r := router.mux.HandleFunc(uri, func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, to, code)
	})

	return NewRoute(r)
}

func (router *Router) Get(uri string, callback interfaces.Callback) interfaces.Route {
	return router.addRoute([]string{http.MethodGet}, uri, callback)
}

func (router *Router) Post(uri string, callback interfaces.Callback) interfaces.Route {
	return router.addRoute([]string{http.MethodPost}, uri, callback)
}

func (router *Router) Put(uri string, callback interfaces.Callback) interfaces.Route {
	return router.addRoute([]string{http.MethodPut}, uri, callback)
}

func (router *Router) Patch(uri string, callback interfaces.Callback) interfaces.Route {
	return router.addRoute([]string{http.MethodPatch}, uri, callback)
}

func (router *Router) Delete(uri string, callback interfaces.Callback) interfaces.Route {
	return router.addRoute([]string{http.MethodDelete}, uri, callback)
}

func (router *Router) Options(uri string, callback interfaces.Callback) interfaces.Route {
	return router.addRoute([]string{http.MethodOptions}, uri, callback)
}

func (router *Router) Any(uri string, callback interfaces.Callback) interfaces.Route {
	return router.addRoute(nil, uri, callback)
}

func (router *Router) Match(methods []string, uri string, callback interfaces.Callback) interfaces.Route {
	return router.addRoute(methods, uri, callback)
}

func (router *Router) Use(m ...interfaces.Middleware) {
	router.middlewares = append(router.middlewares, m...)
}

func (router *Router) Prefix(uri string, fn func(interfaces.Router)) {
	fn(&Router{
		mux: router.mux.NewRoute().PathPrefix(uri).Subrouter(),
		app: router.app,
	})
}

func (router *Router) With(m ...interfaces.Middleware) interfaces.Router {
	newRouter := &Router{
		mux: router.mux.NewRoute().Subrouter(),
		app: router.app,
	}
	newRouter.Use(m...)
	return newRouter
}

// Handler returns the HTTP handler
func (router *Router) Handler() http.Handler {
	return router.handler
}

func (router *Router) URL(name string, params map[string]string) string {
	r := router.mux.Get(name)
	if r == nil {
		log.Panicf("Failed to find URL name=%s", name)
	}

	args := make([]string, 0)
	for k, v := range params {
		args = append(args, k, v)
	}

	u, err := r.URL(args...)
	if err != nil {
		log.Panicf("Failed to get URL name=%s err=%v", name, err)
	}

	return u.Path
}

func (router *Router) addRoute(methods []string, uri string, callback interface{}) interfaces.Route {
	if len(methods) == 0 {
		router.mux.HandleFunc(uri, router.wrap(callback))
	}

	route := router.mux.Methods(methods...).Path(uri).HandlerFunc(router.wrap(callback))
	r := NewRoute(route)
	router.routes = append(router.routes, r)
	return r
}

func (router *Router) wrap(callback interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var renderer templates.Renderer
		router.app.Resolve(&renderer)

		req := newRequest(r)
		res := newResponse(w, r, &renderer)

		newApp := hemlock.CloneApplication(router.app)
		newApp.Instance(req)
		newApp.Instance(res)

		extraArgs := make([]interface{}, 0)
		for _, v := range mux.Vars(r) {
			extraArgs = append(extraArgs, v)
		}

		results := newApp.ResolveInto(callback, extraArgs...)
		if len(results) != 1 {
			panic("Route did not return a value. Got " + string(len(results)))
		}
	}
}
