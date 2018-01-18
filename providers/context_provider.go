package providers

import (
	"context"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
)

type ContextProvider struct{}

func (p *ContextProvider) Register(c interfaces.Container) {
	c.Bind(func (app *hemlock.Application) (context.Context, error) {
		var r interfaces.Request
		app.Resolve(&r)
		return r.Context(), nil
	})
}

func (p *ContextProvider) Boot(app *hemlock.Application) {
	// Nothing yet
}
