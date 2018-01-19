package router

import (
	"fmt"
	"io"
)

type View struct {
	status int
	data   interface{}
}

func (v *View) Status() int {
	return v.status
}

func (v *View) Write(w io.Writer) {
	if b, ok := v.data.([]byte); ok {
		w.Write(b)
	} else if src, ok := v.data.(io.Reader); ok {
		io.Copy(w, src)
	} else {
		w.Write([]byte(fmt.Sprintf("%v", v.data)))
	}
}
