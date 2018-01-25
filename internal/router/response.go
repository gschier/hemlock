package router

import (
	"fmt"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
	"io"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	W              http.ResponseWriter
	r              *http.Request
	renderer       *templates.Renderer
	hasWrittenData bool
}

func newResponse(w http.ResponseWriter, r *http.Request, renderer *templates.Renderer) interfaces.Response {
	return &Response{
		W:        w,
		r:        r,
		renderer: renderer,
	}
}

func (res *Response) Cookie(cookie *http.Cookie) interfaces.Response {
	http.SetCookie(res.W, cookie)
	return res
}

func (res *Response) Status(status int) interfaces.Response {
	res.W.WriteHeader(status)
	return res
}

func (res *Response) View(view, layout string, data interface{}) interfaces.Result {
	ctx := res.getRenderContext(data)
	err := res.renderer.RenderTemplate(res.W, view, layout, ctx)
	if err != nil {
		log.Panicf("Failed to render: %v", err)
	}
	return res.End()
}

func (res *Response) Data(data interface{}) interfaces.Response {
	if !res.hasWrittenData {
		res.setContentTypeHeader()
	}

	if v, ok := data.(string); ok {
		res.W.Write([]byte(v))
	} else if v, ok := data.([]byte); ok {
		res.W.Write(v)
	} else if v, ok := data.(io.Reader); ok {
		io.Copy(res.W, v)
	} else if v, ok := data.(error); ok {
		fmt.Printf("Error: %v\n", v)
		// TODO: Check if status already written
		res.W.WriteHeader(http.StatusInternalServerError)
		res.W.Write([]byte("Internal Server Error"))
	} else {
		v := []byte(fmt.Sprintf("%v", data))
		res.W.Write(v)
	}

	// Remember this the next time
	res.hasWrittenData = true

	return res
}

func (res *Response) Sprintf(format string, a ...interface{}) interfaces.Response {
	return res.Data(fmt.Sprintf(format, a...))
}

func (res *Response) Header(name, value string) interfaces.Response {
	res.W.Header().Set(name, value)
	return res
}

func (res *Response) Redirect(uri string, code int) interfaces.Result {
	http.Redirect(res.W, res.r, uri, code)
	return res.End()
}

func (res *Response) RedirectRoute(name string, code int) interfaces.Result {
	var r interfaces.Router
	hemlock.App().Resolve(&r)
	return res.Redirect(r.URL(name, nil), code)
}

func (res *Response) End() interfaces.Result {
	return &Result{res: res}
}

func (res *Response) setContentTypeHeader() {
	w := res.W
	r := res.r
	if strings.HasSuffix(r.URL.Path, ".js") {
		w.Header().Add("Content-Type", "application/javascript")
	} else if strings.HasSuffix(r.URL.Path, ".css") {
		w.Header().Add("Content-Type", "text/css")
	} else {
		w.Header().Add("Content-Type", "text/html")
	}
}

func (res *Response) getRenderContext (data interface{}) interface{} {
	var config hemlock.Config
	hemlock.App().Resolve(&config)

	return map[string]interface{}{
		"App": map[string]string{
			"Name": config.Name,
			"URL": config.URL,
		},
		"Page": data,
		"URL": res.r.URL.Path,
		"CacheBustKey": hemlock.CacheBustKey,
		"Production": strings.ToLower(config.Env) == "production",
	}
}
