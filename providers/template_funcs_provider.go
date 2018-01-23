package providers

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"html/template"
	"log"
	"net/url"
	"path"
	"strings"
)

type TemplateFuncsProvider struct{}

func (p *TemplateFuncsProvider) Register(c interfaces.Container) {
	c.Bind(func(app *hemlock.Application) (*template.FuncMap, error) {
		fm := template.FuncMap{
			"Asset": func(name string) template.URL {
				base := app.Config.URL
				fullURL := name
				if strings.Contains(base, "://") {
					u, err := url.Parse(app.Config.URL)
					if err != nil {
						log.Panicf("Invalid App URL: %s", base)
					}
					u.Path = path.Join(u.Path, name)
					fullURL = u.String()
				} else {
					fullURL = path.Join(base, name)
				}

				return template.URL(fullURL)
			},
		}
		return &fm, nil
	})
}

func (p *TemplateFuncsProvider) Boot(app *hemlock.Application) error {
	return nil
}
