package router

import (
	"fmt"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
	"io"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	w              http.ResponseWriter
	r              *http.Request
	renderer       *templates.Renderer
	hasWrittenData bool
}

func newResponse(w http.ResponseWriter, r *http.Request, renderer *templates.Renderer) interfaces.Response {
	return &Response{
		w:        w,
		r:        r,
		renderer: renderer,
	}
}

func (res *Response) Cookie(name, value string) interfaces.Response {
	c := &http.Cookie{Name: name, Value: value}
	http.SetCookie(res.w, c)
	return res
}

func (res *Response) Status(status int) interfaces.Response {
	res.w.WriteHeader(status)
	return res
}

func (res *Response) Render(name, base string, data interface{}) interfaces.Response {
	ctx := res.getRenderContext(data)
	err := res.renderer.RenderTemplate(res.w, name, base, ctx)
	if err != nil {
		panic("Failed to render: " + err.Error())
	}
	return res
}

func (res *Response) Data(data interface{}) interfaces.Response {
	if !res.hasWrittenData {
		res.setContentTypeHeader()
	}

	if v, ok := data.(string); ok {
		res.w.Write([]byte(v))
	} else if v, ok := data.([]byte); ok {
		res.w.Write(v)
	} else if v, ok := data.(io.Reader); ok {
		io.Copy(res.w, v)
	} else {
		log.Panicf("Unknown response type: %#v\n", data)
	}

	// Remember this the next time
	res.hasWrittenData = true

	return res
}

func (res *Response) Dataf(format string, a ...interface{}) interfaces.Response {
	return res.Data(fmt.Sprintf(format, a...))
}

func (res *Response) End() interfaces.Result {
	return &Result{res: res}
}

func (res *Response) setContentTypeHeader() {
	w := res.w
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
	return map[string]interface{}{
		"URL": res.r.URL.Path,
		"Page": data,
	}
}
