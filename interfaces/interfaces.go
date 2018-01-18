package interfaces

import (
	"context"
	"io"
	"net/http"
)

type Container interface {
	// Bind binds the type of v as a dependency
	Bind(fn interface{})

	// Singleton binds the type of v as a dependency. Will only get instantiated once
	Singleton(fn interface{})

	// Instance binds an already-created value as a dependency
	Instance(i interface{})
}

type Request interface {
	// Input grabs input from the request body by name
	Input(name string) string

	// Input grabs input from the query string by name
	Query(name string) string

	// Cookie grabs input from cookies by name
	Cookie(name string) string

	// File grabs the file
	File(name string) io.Reader

	// Context returns the context.Context of the current request
	Context() context.Context
}

type Response interface {
	Cookie(name, value string) Response
	Status(status int) Response
	Data(data interface{}) Response
	Dataf(format string, a ...interface{}) Response
	View() View
}

type View interface {
	// Nothing yet
}

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
type Callback = interface{}
