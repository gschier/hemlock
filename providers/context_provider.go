package providers

import (
	"context"
	"github.com/gschier/hemlock"
	"net/http"
)

type ContextProvider struct{}

func (p *ContextProvider) Register(c hemlock.Container) {
	c.Bind(func (app *hemlock.Application) (context.Context, error) {
		var r *http.Request = nil
		app.Resolve(r)
		if r == nil {
			panic("BACKGROUND")
			return context.Background(), nil
		}
		return r.Context(), nil
	})
}

func (p *ContextProvider) Boot(app *hemlock.Application) {
	// Nothing yet
}
