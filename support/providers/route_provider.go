package providers

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/router"
)

type RouteProvider struct{}

func (p *RouteProvider) Register(c interfaces.Container) {
	c.Singleton(func(app *hemlock.Application) (interfaces.Router, error) {
		return router.NewRouter(app), nil
	})
}

func (p *RouteProvider) Boot(*hemlock.Application) error {
	return nil
}
