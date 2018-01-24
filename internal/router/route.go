package router

import (
	"github.com/gorilla/mux"
	"github.com/gschier/hemlock/interfaces"
)


type Route struct {
	route *mux.Route
}

func NewRoute(route *mux.Route) *Route {
	return &Route{route: route}
}

func (r *Route) Name(n string) interfaces.Route {
	r.route.Name(n)
	return r
}
