package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type Router struct {
	handler     http.Handler
	app         *hemlock.Application
	mux         *mux.Router
	middlewares []interfaces.Middleware
}

func NewRouter(app *hemlock.Application) *Router {
	router := &Router{app: app, mux: mux.NewRouter()}

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
			p := r.URL.Path
			cwd, _ := os.Getwd()
			fullPath := filepath.Join(cwd, app.Config.PublicDirectory, p)
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
	return router.newRoute().Redirect(uri, to, code)
}

func (router *Router) Get(uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Get(uri, callback)
}

func (router *Router) Post(uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Post(uri, callback)
}

func (router *Router) Put(uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Put(uri, callback)
}

func (router *Router) Patch(uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Patch(uri, callback)
}

func (router *Router) Delete(uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Delete(uri, callback)
}

func (router *Router) Options(uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Options(uri, callback)
}

func (router *Router) Any(uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Any(uri, callback)
}

func (router *Router) Match(methods []string, uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Match(methods, uri, callback)
}

func (router *Router) Prefix(uri string) interfaces.Route {
	return router.newRoute().Prefix(uri)
}

func (router *Router) Host(hostname string) interfaces.Route {
	return router.newRoute().Host(hostname)
}

func (router *Router) With(m ...interfaces.Middleware) interfaces.Route {
	return router.newRoute().With(m...)
}

func (router *Router) Use(m ...interfaces.Middleware) {
	router.useMiddleware(m...)
}

// Handler returns the HTTP handler
func (router *Router) Handler() http.Handler {
	return router.handler
}

func (router *Router) Route(name string, params interfaces.RouteParams) string {
	r := router.mux.Get(name)
	if r == nil {
		log.Panicf("Failed to find route by name '%s'", name)
	}

	args := make([]string, 0)
	for k, v := range params {
		args = append(args, k, v)
	}

	u, err := r.URL(args...)
	if err != nil {
		log.Panicf("Failed to get URL name=%s args=%v Error: %s", name, args, err)
	}

	return router.URL(u.Path)
}

func (router *Router) URL(p string) string {
	base := router.app.Config.URL
	u, err := url.Parse(base)
	if err != nil {
		log.Panicf("Invalid App URL: %s", base)
	}

	u.Path = path.Join(u.Path, p)

	return u.String()
}

func (router *Router) fork(mux *mux.Router) *Router {
	return &Router{mux: mux, app: router.app}
}

func (router *Router) newRoute() *Route {
	return NewRoute(router, router.mux.NewRoute())
}

func (router *Router) useMiddleware(m ...interfaces.Middleware) {
	router.middlewares = append(router.middlewares, m...)
}
