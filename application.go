package hemlock

import (
	"context"
	"github.com/gschier/hemlock/internal/container"
)

type Application struct {
	Config    *Config
	container *container.Container
	ctx       context.Context
}

func NewApplication(config *Config) *Application {
	app := &Application{
		Config: config,
	}

	// Ensure all service constructors take in *Application as an
	// argument
	serviceConstructorArgs := []interface{}{app}
	app.container = container.New(serviceConstructorArgs)

	// Add providers from config
	for _, p := range app.Config.Providers {
		p.Register(app)
	}

	// Boot all providers
	for _, p := range app.Config.Providers {
		p.Boot(app)
	}

	return app
}

func CloneApplication(app *Application) *Application {
	return &Application{
		Config:    app.Config,
		container: app.container.Clone(),
	}
}

func (a *Application) Bind(fn interface{}) {
	a.container.Bind(fn)
}

func (a *Application) Singleton(fn interface{}) {
	a.container.Singleton(fn)
}

func (a *Application) Instance(v interface{}) {
	a.container.Instance(v)
}

func (a *Application) ResolveInto(fn interface{}, extraArgs ...interface{}) []interface{} {
	return a.container.Call(fn, extraArgs)
}

func (a *Application) Make(i interface{}) interface{} {
	return a.container.Make(i)
}

func (a *Application) Resolve(v interface{}) {
	a.container.Resolve(v)
}

func (a *Application) Env(name string) string {
	return Env(name)
}

func (a *Application) EnvOr(name, fallback string) string {
	return EnvOr(name, fallback)
}
