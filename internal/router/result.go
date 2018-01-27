package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type Result struct {
	w http.ResponseWriter
	r *http.Request

	status   int
	renderer *templates.Renderer
	router   *Router
	error    error

	// hasSentHeaders signifies that data has already been written
	// and headers can no longer be applied
	hasSentHeaders bool
	hasSentData    bool
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
	r.error = err

	r.flushHeaders()

	// Render error view if we haven't sent data
	if !r.hasSentData {
		return r.View("error.html", "", nil)
	}

	return r
}

func (r *Result) Sprintf(format string, a ...interface{}) interfaces.Result {
	return r.Data(fmt.Sprintf(format, a...))
}

func (r *Result) View(name, layout string, data interface{}) interfaces.Result {
	// Set content type based on extension of template
	ext := filepath.Ext(name)
	r.w.Header().Set("Content-Type", mime.TypeByExtension(ext))

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

	r.hasSentData = true

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
	if r.hasSentHeaders {
		return
	}

	// Set status if not set yet
	if r.status == 0 {
		r.status = http.StatusOK
	}

	// Send error response if we have not written anything yet
	if r.error != nil && r.status == 0 {
		r.status = http.StatusInternalServerError
	}

	// Log error if there was one
	// TODO: Make this better
	if r.error != nil {
		fmt.Printf("[router] Error: %v\n", r.error)
	}

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

	// Write the status code and headers
	r.w.WriteHeader(r.status)

	// Mark that we've done it because we're not allowed to write
	// any headers after this point
	r.hasSentHeaders = true
}
