package router

import (
	"github.com/gschier/hemlock/interfaces"
)

type Result struct {
	res *Response
}

func (v *Result) Response() interfaces.Response {
	return v.res
}
