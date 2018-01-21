package router

import (
	"fmt"
	"html/template"
	"io"
)

type View struct {
	status       int
	data         interface{}
	templateData interface{}
	template     *template.Template
}

func (v *View) Status() int {
	if v.status == 0 {
		return 200
	}

	return v.status
}

func (v *View) Write(w io.Writer) {
	if v.template != nil {
		v.template.Execute(w, v.templateData)
		return
	}

	if b, ok := v.data.([]byte); ok {
		w.Write(b)
	} else if src, ok := v.data.(io.Reader); ok {
		io.Copy(w, src)
	} else {
		w.Write([]byte(fmt.Sprintf("%v", v.data)))
	}
}
