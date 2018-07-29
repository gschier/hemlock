package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type middlewareContainer struct {
	hemlock interfaces.Middleware
	native  func(http.Handler) http.Handler
}

type Router struct {
	app          *hemlock.Application
	mux          *mux.Router
	middlewares  []*middlewareContainer
	didSetupURLs bool
}

func NewRouter(app *hemlock.Application) *Router {
	router := &Router{app: app, mux: mux.NewRouter()}

	// Redirect slashes
	router.mux.StrictSlash(true)

	// This needs to be first
	if !app.IsDev() {
		router.UseG(handlers.RecoveryHandler())
	}

	router.UseG(handlers.CompressHandler)

	// Add logging middleware
	if app.IsDev() {
		router.UseG(func(next http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, next)
		})
	}

	if !app.IsDev() {
		// Add caching middleware
		router.Use(func(req interfaces.Request, res interfaces.Response, next interfaces.Next) interfaces.Result {
			ext := filepath.Ext(req.Path())
			if ext == ".css" || ext == ".js" {
				res.Header("Cache-Control", "public, max-age=7200")
			}
			return next(req, res)
		})
		// Add HTTP redirect middleware
		router.Use(func(req interfaces.Request, res interfaces.Response, next interfaces.Next) interfaces.Result {
			if req.Header("X-Forwarded-Proto") == "http" {
				newUrl := "https://" + req.Host() + req.URL().String()
				return res.Redirect(newUrl, http.StatusFound)
			} else {
				return next(req, res)
			}
		})
	}

	// Add main handler to call middleware
	router.mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			router.nextCombinedMiddleware(0, w, r, func(w2 http.ResponseWriter, r2 *http.Request) {
				next.ServeHTTP(w2, r2)
			})
		})
	})

	// Add static handler
	u, err := url.Parse(app.Config.PublicPrefix)
	if err == nil { // Make sure PublicPrefix is a valid path or URL
		publicPrefixPath := u.Path
		router.Prefix(publicPrefixPath).Methods(http.MethodGet).Callback(
			func(req interfaces.Request, res interfaces.Response) interfaces.Result {
				p := req.Path()
				p = strings.TrimPrefix(p, publicPrefixPath)
				cwd, _ := os.Getwd()
				fullPath := filepath.Join(cwd, app.Config.PublicDirectory, p)
				s, err := os.Stat(fullPath)
				if err != nil || s.IsDir() {
					return res.Status(404).Data("Resource not found")
				}
				f, err := os.Open(fullPath)

				ext := filepath.Ext(fullPath)
				return res.Header("Content-Type", mime.TypeByExtension(ext)).Data(f)
			},
		)
	}

	return router
}

func (router *Router) Redirect(uri, to string, code int) interfaces.Route {
	return router.newRoute().Redirect(uri, to, code)
}

func (router *Router) View(uri, view, layout string, data map[string]interface{}) interfaces.Route {
	return router.newRoute().View(uri, view, layout, data)
}

func (router *Router) Callback(callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Callback(callback)
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

func (router *Router) Head(uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Head(uri, callback)
}

func (router *Router) Connect(uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Connect(uri, callback)
}

func (router *Router) Trace(uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Trace(uri, callback)
}

func (router *Router) Any(uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Any(uri, callback)
}

func (router *Router) Match(methods []string, uri string, callback interfaces.Callback) interfaces.Route {
	return router.newRoute().Match(methods, uri, callback)
}

func (router *Router) Methods(methods ...string) interfaces.Route {
	return router.newRoute().Methods(methods...)
}

func (router *Router) Prefix(uri string) interfaces.Route {
	return router.newRoute().Prefix(uri)
}

func (router *Router) Host(hostname string) interfaces.Route {
	return router.newRoute().Host(hostname)
}

func (router *Router) With(m ...interfaces.Middleware) interfaces.Route {
	newRouter := router.fork()
	newRouter.Use(m...)
	return newRouter.newRoute()
}

func (router *Router) WithG(m ...func(http.Handler) http.Handler) interfaces.Route {
	newRouter := router.fork()
	newRouter.UseG(m...)
	return newRouter.newRoute()
}

func (router *Router) Use(m ...interfaces.Middleware) {
	for _, m := range m {
		router.middlewares = append(router.middlewares, &middlewareContainer{
			hemlock: m,
		})
	}
}

func (router *Router) UseG(m ...func(http.Handler) http.Handler) {
	for _, m := range m {
		router.middlewares = append(router.middlewares, &middlewareContainer{
			native: m,
		})
	}
}

// Handler returns the HTTP handler
func (router *Router) Handler() http.Handler {
	return router.mux
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

func (router *Router) fork() *Router {
	return &Router{mux: router.mux.NewRoute().Subrouter(), app: router.app}
}

func (router *Router) newRoute() *Route {
	return NewRoute(router, router.mux.NewRoute())
}

func (router *Router) nextCombinedMiddleware(
	i int,
	w http.ResponseWriter,
	r *http.Request,
	fn func(w http.ResponseWriter, r *http.Request),
) {
	if i == len(router.middlewares) {
		fn(w, r)
		return
	}

	m := router.middlewares[i]
	if m.hemlock != nil {
		next := func(newReq interfaces.Request, newRes interfaces.Response) interfaces.Result {
			router.nextCombinedMiddleware(i+1, newRes.(*Response).W, newReq.(*Request).R, fn)
			return nil
		}

		var renderer templates.Renderer
		router.app.Resolve(&renderer)
		req := newRequest(r)
		res := newResponse(w, req, &renderer, router)
		m.hemlock(req, res, next)
	} else {
		m.native(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			router.nextCombinedMiddleware(i+1, w, r, fn)
		})).ServeHTTP(w, r)
	}
}
