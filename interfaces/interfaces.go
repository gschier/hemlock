package interfaces

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

// Container is used to bind objects to the application
type Container interface {
	// Bind binds the type of v as a dependency
	Bind(fn interface{})

	// Singleton binds the type of v as a dependency. Will only get instantiated once
	Singleton(fn interface{})

	// Instance binds an already-created value as a dependency
	Instance(i interface{})
}

// Response is used to retrieve data from an HTTP request
type Request interface {
	// RouteName returns the name of the current route
	RouteName() string

	// Path returns the path of the current URL
	Path() string

	// URL returns the raw URL
	URL() *url.URL

	// Header return a header value by name. If the header is not found
	// an empty string will be returned.
	Header(name string) string

	// Host returns the host of the request.
	Host() string

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

	// WithContext returns a shallow copy of the request with a new context
	WithContext(ctx context.Context) Request
}

// Response is be used to send data to the client
type Response interface {
	// Cookie sets an HTTP cookie on the response
	Cookie(cookie *http.Cookie) Response

	// Status sets the HTTP status code of the response. This can only be called once.
	Status(status int) Response

	// Header adds an HTTP header to the response
	Header(name, value string) Response

	// Data responds with data provided
	//
	// Most types will converted to a string representation except structs,
	// which will be serialized to JSON.
	Data(data interface{}) Result

	// Error sends the default 500 response and logs the error message
	Error(error) Result

	// Sprintf builds a response using `fmt.Sprintf`
	Sprintf(format string, a ...interface{}) Result

	// View renders a view for the response with a provided layout and data
	View(name, layout string, data interface{}) Result

	// Redirect redirects the client to a URL
	Redirect(uri string, code int) Result

	// RedirectRoute is the same as Redirect but looks up a URL from a route name
	RedirectRoute(name string, params map[string]string, code int) Result

	// End ends the response chain
	End() Result
}

type Result interface {
	Data(data interface{}) Result
	Error(error) Result
	Sprintf(format string, a ...interface{}) Result
	View(name, layout string, data interface{}) Result
	Redirect(uri string, code int) Result
	RedirectRoute(name string, params map[string]string, code int) Result
}

// Router provides the ability to define HTTP routes
type Router interface {
	Callback(callback Callback) Route
	Get(uri string, callback Callback) Route
	Head(uri string, callback Callback) Route
	Post(uri string, callback Callback) Route
	Put(uri string, callback Callback) Route
	Patch(uri string, callback Callback) Route
	Delete(uri string, callback Callback) Route
	Options(uri string, callback Callback) Route
	Any(uri string, callback Callback) Route
	Match(methods []string, uri string, callback Callback) Route
	Methods(methods ...string) Route

	// Redirect creates a static redirect route for the provided URI
	Redirect(uri, to string, code int) Route

	// View creates a route
	View(uri, view, layout string, data interface{}) Route

	// Prefix resolves new Route instance for the provided URI
	Prefix(uri string) Route

	// Host resolves a new Route instance for the provided Host
	//
	// For example:
	//
	//     r.Host("www.example.com", func (r *hemlock.Router) { ... })
	//     r.Host("{subdomain}.domain.com", func (r *hemlock.Router) { ... }))
	//     r.Host("{subdomain:[a-z]+}.domain.com", func (r *hemlock.Router) { ... }))
	Host(uri string) Route

	// With returns a new Router with middleware applied
	//
	// For example:
	//
	//     admin := r.With(adminMiddleware)
	//	   admin.Get("/admin", adminDashboard)
	With(...Middleware) Route
	WithG(...func(http.Handler) http.Handler) Route

	Use(...Middleware)
	UseG(...func(http.Handler) http.Handler)

	// URL returns a URL based on an assigned route name
	URL(path string) string

	// Route resolves a URL from a route name combined with parameters
	Route(name string, params RouteParams) string

	// TODO: Make this private
	Handler() http.Handler
}

type RouteParams map[string]string

// Route represents an HTTP route
type Route interface {
	Callback(callback Callback) Route
	Get(uri string, callback Callback) Route
	Post(uri string, callback Callback) Route
	Put(uri string, callback Callback) Route
	Patch(uri string, callback Callback) Route
	Delete(uri string, callback Callback) Route
	Options(uri string, callback Callback) Route
	Head(uri string, callback Callback) Route
	Connect(uri string, callback Callback) Route
	Trace(uri string, callback Callback) Route
	Any(uri string, callback Callback) Route
	Match(methods []string, uri string, callback Callback) Route
	Methods(methods ...string) Route

	Redirect(uri, to string, code int) Route

	// Name assigns a name to the route for referencing
	Name(name string) Route
	Host(uri string) Route
	Prefix(uri string) Route
	Group(func(Router))
	With(...Middleware) Route
	WithG(...func(http.Handler) http.Handler) Route
	Use(...Middleware)
	UseG(...func(http.Handler) http.Handler)
}

// Callback is a function that takes injected arguments
type Callback interface{}

// Next is called to continue the chain of  middleware
type Next func(req Request, res Response) Result

// Middleware is an interface for adding middleware to a Router instance
type Middleware func(req Request, res Response, next Next) Result
