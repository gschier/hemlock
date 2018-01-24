package interfaces

import (
	"context"
	"io"
	"net/http"
	"net/url"
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
	// URL returns the URL of the request
	URL() *url.URL

	// Input grabs input from the query string by name
	Query(name string) string

	// Input grabs input from the query string by name
	QueryInt(name string) int

	// Cookie grabs input from cookies by name
	Cookie(name string) string

	// File grabs the file
	File(name string) io.Reader

	// Context returns the context.Context of the current request
	Context() context.Context
}

type Response interface {
	Cookie(cookie *http.Cookie) Response
	Status(status int) Response
	Header(name, value string) Response
	Render(name, base string, data interface{}) Response
	Data(data interface{}) Response
	Dataf(format string, a ...interface{}) Response

	Redirect(uri string, code int) Result
	RedirectRoute(name string, code int) Result
	End() Result
}

type Result interface {
	Response() Response
}

type Router interface {
	Redirect(uri, to string, code int) Route
	Get(uri string, callback Callback) Route
	Post(uri string, callback Callback) Route
	Put(uri string, callback Callback) Route
	Patch(uri string, callback Callback) Route
	Delete(uri string, callback Callback) Route
	Options(uri string, callback Callback) Route
	Any(uri string, callback Callback) Route
	Match(methods []string, uri string, callback Callback) Route

	// Use appends middleware to the chain
	Use(...Middleware)

	// Prefix returns a new router instance for the provided URI
	Prefix(uri string, fn func(Router))

	// Group
	With(...Middleware) Router

	// URL Returns a URL based on an assigned route name
	URL(name string, params map[string]string) string

	// TODO: Make this private
	Handler() http.Handler
}

type Route interface {
	// Name assigns a name to the route for referencing
	Name(name string) Route
}

// Callback is a function that takes injected arguments
type Callback interface{}

type Next func(req Request, res Response) Result
type Middleware func(req Request, res Response, next Next) Result
