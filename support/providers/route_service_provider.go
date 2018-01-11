package providers

import (
	"github.com/go-chi/chi"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
)

type RouteServiceProvider struct{}

func (p *RouteServiceProvider) Register(c *hemlock.Container) {
	p.registerRouter(c)
}

func (p *RouteServiceProvider) Boot(*hemlock.Application) {
	// Nothing
}

func (p *RouteServiceProvider) registerRouter(c *hemlock.Container) {
	c.Singleton(func(app *hemlock.Application) (interfaces.Router, error) {
		return chi.NewRouter(), nil
	})
}
