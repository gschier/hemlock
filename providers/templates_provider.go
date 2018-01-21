package providers

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
)

type TemplatesProvider struct{}

func (p *TemplatesProvider) Register(c interfaces.Container) {
	c.Singleton(func(app *hemlock.Application) (*templates.Renderer, error) {
		dir := app.ResolveDir(app.Config.TemplatesDirectory)
		return templates.NewRenderer(dir), nil
	})
}

func (p *TemplatesProvider) Boot(app *hemlock.Application) error {
	var renderer templates.Renderer
	app.Resolve(&renderer)

	err := renderer.Load()
	if err != nil {
		return err
	}

	return nil
}
