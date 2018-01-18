package providers

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/internal/router"
	"github.com/gschier/hemlock/interfaces"
)

type RouteServiceProvider struct{}

func (p *RouteServiceProvider) Register(c interfaces.Container) {
	p.registerRouter(c)
}

func (p *RouteServiceProvider) Boot(*hemlock.Application) {
	// Nothing
}

func (p *RouteServiceProvider) registerRouter(c interfaces.Container) {
	c.Singleton(func(app *hemlock.Application) (interfaces.Router, error) {
		return router.NewRouter(app), nil
	})
}
