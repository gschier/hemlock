package router

import (
	"fmt"
	"github.com/gschier/hemlock/interfaces"
	"html/template"
	"net/http"
)

type Response struct {
	w            http.ResponseWriter
	status       int
	data         interface{}
	templateData interface{}
	templateName string
	cookies      map[string]string
}

func newResponse(w http.ResponseWriter) interfaces.Response {
	return &Response{w: w}
}

func (res *Response) Cookie(name, value string) interfaces.Response {
	res.cookies[name] = value
	return res
}

func (res *Response) Status(status int) interfaces.Response {
	res.status = status
	return res
}

func (res *Response) Template(name string, data interface{}) interfaces.Response {
	res.templateName = name
	res.templateData = data
	return res
}

func (res *Response) Data(data interface{}) interfaces.Response {
	res.data = data
	return res
}

func (res *Response) Dataf(format string, a ...interface{}) interfaces.Response {
	res.data = fmt.Sprintf(format, a...)
	return res
}

func (res *Response) View() interfaces.View {
	v := &View{
		status:       res.status,
		data:         res.data,
	}

	if res.templateName != "" {
		t, err := template.ParseFiles("resources/templates/" + res.templateName + ".html")
		v.template = t
		v.templateData = res.templateData
		if err != nil {
			panic("Failed to parse template: " + err.Error())
		}
	}

	return v
}
