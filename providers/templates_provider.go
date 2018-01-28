package providers

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates"
	"html/template"
)

type TemplatesProvider struct{}

func (p *TemplatesProvider) Register(c interfaces.Container) {
	c.Singleton(func(app *hemlock.Application) (*templates.renderer, error) {
		dir := app.Path(app.Config.TemplatesDirectory)

		var fm template.FuncMap
		app.Resolve(&fm)

		r := templates.NewRenderer(dir, fm)

		err := r.Init()
		if err != nil {
			return nil, err
		}

		return r, nil
	})
}

func (p *TemplatesProvider) Boot(a *hemlock.Application) error {
	return nil
}
