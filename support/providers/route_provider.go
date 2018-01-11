package providers

import (
	"github.com/go-chi/chi"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/facades"
)

type RouteProvider struct{}

func (p *RouteProvider) Register(c *hemlock.Container) {
	p.registerRouter(c)
}

func (p *RouteProvider) Boot(*hemlock.Application) {
	// TODO: Add routes or something?
}

func (p *RouteProvider) registerRouter(c *hemlock.Container) {
	c.Singleton(func(app *hemlock.Application) (facades.Router, error) {
		return chi.NewRouter(), nil
	})
}
