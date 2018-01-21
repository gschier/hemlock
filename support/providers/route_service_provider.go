package providers

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/router"
)

type RouteServiceProvider struct{}

func (p *RouteServiceProvider) Register(c interfaces.Container) {
	c.Singleton(func(app *hemlock.Application) (interfaces.Router, error) {
		return router.NewRouter(app), nil
	})
}

func (p *RouteServiceProvider) Boot(*hemlock.Application) error {
	return nil
}
