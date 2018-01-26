package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Result struct {
	w http.ResponseWriter
	r *http.Request

	status   int
	renderer *templates.Renderer
	router   *Router

	// hasWrittenData signifies that data has already been written
	// and headers can no longer be applied
	hasWrittenData bool
}

func newResult(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	renderer *templates.Renderer,
	router *Router,
) interfaces.Result {
	return &Result{w: w, r: r, status: status, renderer: renderer, router: router}
}

func (r *Result) Redirect(uri string, code int) interfaces.Result {
	http.Redirect(r.w, r.r, uri, code)
	return r
}

func (r *Result) RedirectRoute(name string, params map[string]string, code int) interfaces.Result {
	return r.Redirect(r.router.Route(name, params), code)
}

func (r *Result) Error(err error) interfaces.Result {
	r.status = http.StatusInternalServerError
	// TODO: Make this better
	fmt.Printf("[router] Error: %v\n", err)
	return r.Data("Internal server error")
}

func (r *Result) Sprintf(format string, a ...interface{}) interfaces.Result {
	return r.Data(fmt.Sprintf(format, a...))
}

func (r *Result) View(name, layout string, data interface{}) interfaces.Result {
	r.flushHeaders()

	ctx := r.getRenderContext(data)
	err := r.renderer.RenderTemplate(r.w, name, layout, ctx)
	if err != nil {
		log.Panicf("Failed to render: %v", err)
	}
	return r
}

func (r *Result) Data(data interface{}) interfaces.Result {
	r.flushHeaders()

	if v, ok := data.(string); ok {
		r.w.Write([]byte(v))
	} else if v, ok := data.([]byte); ok {
		r.w.Write(v)
	} else if v, ok := data.(io.Reader); ok {
		io.Copy(r.w, v)
	} else {
		v := []byte(fmt.Sprintf("%v", data))
		r.w.Write(v)
	}

	return r
}

func (r *Result) getRenderContext(data interface{}) interface{} {
	var config hemlock.Config
	var router interfaces.Router
	hemlock.App().Resolve(&config, &router)

	u, _ := url.Parse(config.URL)
	u.Path = r.r.URL.Path

	return map[string]interface{}{
		"App": map[string]string{
			"Name": config.Name,
			"URL":  config.URL,
		},
		"Page":         data,
		"CacheBustKey": hemlock.CacheBustKey,
		"Production":   strings.ToLower(config.Env) == "production",
		"Request": map[string]string{
			"URL":   u.String(),
			"Path":  r.r.URL.Path,
			"Route": mux.CurrentRoute(r.r).GetName(),
		},
	}
}

func (r *Result) flushHeaders() {
	if r.hasWrittenData {
		return
	}

	// Write the status code
	r.w.WriteHeader(r.status)

	// Set any headers that need to be
	if r.w.Header().Get("Content-Type") == "" {
		if strings.HasSuffix(r.r.URL.Path, ".js") {
			r.w.Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(r.r.URL.Path, ".css") {
			r.w.Header().Set("Content-Type", "text/css")
		} else {
			r.w.Header().Set("Content-Type", "text/html")
		}
	}

	// Mark that we've done it
	r.hasWrittenData = true
}
