package providers

import (
	"fmt"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"log"
	"net/http"
)

type HttpProvider struct{}

func (p *HttpProvider) Register(c hemlock.Container) {
	// Nothing yet
}

func (p *HttpProvider) Boot(app *hemlock.Application) {
	var router interfaces.Router
	app.Resolve(&router)

	srv := &http.Server{
		Handler: router.Handler(),
		Addr: app.Config.HTTP.Host+":"+app.Config.HTTP.Port,
	}

	fmt.Printf("Starting server at %v\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
