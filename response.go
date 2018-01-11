package hemlock

import (
	"fmt"
	"github.com/gschier/hemlock/interfaces"
	"net/http"
)

type Response struct {
	w http.ResponseWriter
	status int
	data interface{}
	cookies map[string]string
}

func NewResponse(w http.ResponseWriter) interfaces.Response {
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

func (res *Response) Data(data interface{}) interfaces.Response {
	res.data = data
	return res
}

func (res *Response) View() interfaces.View {
	return &View{
		Status: res.status,
		Bytes: []byte(fmt.Sprintf("%v", res.data)),
	}
}
