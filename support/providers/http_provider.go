package providers

import (
	"github.com/gschier/hemlock"
	"net/http"
)

type HttpProvider struct{}

func (p *HttpProvider) Register(c *hemlock.Container) {
	//p.registerServer(c)
}

func (p *HttpProvider) Boot(app *hemlock.Application) {
	//srv := &http.Server{
	//
	//}
}

func (p *HttpProvider) registerServer(c *hemlock.Container) {
	//c.Singleton(func(app *hemlock.Application) (facades.Router, error) {
	//	return chi.NewRouter(), nil
	//})
}
