package interfaces

import "net/http"

type Router interface {
	Redirect(uri, to string, code int)
	Get(uri string, callback Callback)
	Post(uri string, callback Callback)
	Put(uri string, callback Callback)
	Patch(uri string, callback Callback)
	Delete(uri string, callback Callback)
	Options(uri string, callback Callback)
	Handler() http.Handler
}

// Callback is a function that takes injected arguments
type Callback interface{}
