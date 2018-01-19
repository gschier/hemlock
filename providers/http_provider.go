package providers

import (
	"fmt"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"log"
	"net/http"
	"time"
)

type HttpProvider struct{}

func (p *HttpProvider) Register(c interfaces.Container) {
	// Nothing yet
}

func (p *HttpProvider) Boot(app *hemlock.Application) {
	var router interfaces.Router
	app.Resolve(&router)

	srv := &http.Server{
		Handler: router.Handler(),
		Addr:    app.Config.HTTP.Host + ":" + app.Config.HTTP.Port,
	}

	go func() {
		<- time.NewTimer(time.Second * 1).C
		fmt.Printf("Started server at %v\n", srv.Addr)
	}()

	log.Fatal(srv.ListenAndServe())
}
