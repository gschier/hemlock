package providers

import (
	"fmt"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"log"
	"net/http"
)

type HttpProvider struct{}

func (p *HttpProvider) Register(c interfaces.Container) {
	// Nothing yet
}

func (p *HttpProvider) Boot(app *hemlock.Application) error {
	var router interfaces.Router
	app.Resolve(&router)

	srv := &http.Server{
		Handler: router.Handler(),
		Addr:    app.Config.HTTP.Host + ":" + app.Config.HTTP.Port,
	}

	go func () {
		fmt.Printf("Started server at %v\n", srv.Addr)
	}()

	log.Fatal(srv.ListenAndServe())

	return nil
}
