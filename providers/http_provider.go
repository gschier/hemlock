package providers

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"net/http"
)

type HttpProvider struct{}

func (p *HttpProvider) Register(c interfaces.Container) {
	c.Singleton(func(app *hemlock.Application) (*http.Server, error) {
		return &http.Server{
			Addr: app.Config.HTTP.Host + ":" + app.Config.HTTP.Port,
		}, nil
	})
}

func (p *HttpProvider) Boot(app *hemlock.Application) error {
	return nil
}
