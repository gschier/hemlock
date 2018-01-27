package providers

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"github.com/gschier/hemlock/internal/templates/funcs"
	"html/template"
)

type TemplateFuncsProvider struct{}

func (p *TemplateFuncsProvider) Register(c interfaces.Container) {
	c.Bind(func(app *hemlock.Application) (*template.FuncMap, error) {
		return funcs.Funcs(app), nil
	})
}

func (p *TemplateFuncsProvider) Boot(app *hemlock.Application) error {
	return nil
}
