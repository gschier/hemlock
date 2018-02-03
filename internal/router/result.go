package router

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
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
	// Only do something if there is an actual error
	if err == nil {
		return r
	}

	r.defaultStatus(http.StatusInternalServerError)
	r.w.Header().Set("Content-Type", "text/plain")

	// Log error if there was one
	// TODO: Make this better
	fmt.Printf("[router] Error: %v\n", err)

	// Respond with error page
	// TODO: Use user-provided error page if available
	return r.Data("Internal Server Error")
}

func (r *Result) Sprintf(format string, a ...interface{}) interfaces.Result {
	return r.Data(fmt.Sprintf(format, a...))
}

func (r *Result) View(name, layout string, data interface{}) interfaces.Result {
	// Set content type based on extension of template
	ext := filepath.Ext(name)
	r.w.Header().Set("Content-Type", mime.TypeByExtension(ext))

	r.flushHeaders()
	r.hasSentData = true

	ctx := r.getRenderContext(data)
	err := r.renderer.RenderTemplate(r.w, name, layout, ctx)
	if err != nil {
		return r.Error(err)
	}
	return r
}

func (r *Result) Data(data interface{}) interfaces.Result {
	// If it's an error, call Result.Error instead
	if err, ok := data.(error); ok {
		return r.Error(err)
	}

	r.flushHeaders()

	r.hasSentData = true

	if v, ok := data.([]byte); ok {
		r.w.Write(v)
	} else if v, ok := data.(io.Reader); ok {
		io.Copy(r.w, v)
	} else {
		dataType := reflect.TypeOf(data)
		if dataType.Kind() == reflect.Ptr {
			dataType = dataType.Elem()
		}

		if dataType.Kind() == reflect.Struct {
			e := json.NewEncoder(r.w)
			e.Encode(data)
		} else {
			v := []byte(fmt.Sprintf("%v", data))
			r.w.Write(v)
		}
	}

	return r
}

func (r *Result) getRenderContext(data interface{}) interface{} {
	var config hemlock.Config
	var router interfaces.Router
	r.router.app.Resolve(&config, &router)

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

	r.applyDefaultHeaders()

	// Write the status code and headers
	r.w.WriteHeader(r.status)

	// Mark that we've done it because we're not allowed to write
	// any headers after this point
	r.hasSentHeaders = true
}

func (r *Result) defaultStatus(code int) {
	if r.status == 0 {
		r.status = code
	}
}

func (r *Result) defaultHeader(name, value string) {
	if r.w.Header().Get(name) == "" {
		r.w.Header().Set(name, value)
	}
}

func (r *Result) applyDefaultHeaders() {
	r.defaultHeader("Server", "Hemlock/"+hemlock.Version())
	r.defaultHeader("Content-Type", mime.TypeByExtension(filepath.Ext(r.r.URL.Path)))
}
