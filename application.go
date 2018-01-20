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

func NewApplication(config *Config, providers []Provider) *Application {
	app := &Application{
		Config: config,
	}

	// Ensure all constructors take in *Application as an argument
	serviceConstructorArgs := []interface{}{app}
	app.container = container.New(serviceConstructorArgs)

	// Bind some useful things to container
	app.Instance(app)
	app.Instance(app.Config)

	// Add providers from config
	for _, p := range providers {
		p.Register(app.container)
	}

	// Boot all providers
	for _, p := range providers {
		p.Boot(app)
	}

	return app
}

func CloneApplication(app *Application) *Application {
	newApp := &Application{Config: app.Config}

	// Ensure all constructors take in *Application as an argument
	serviceConstructorArgs := []interface{}{newApp}
	newApp.container = container.Clone(app.container, serviceConstructorArgs)

	return app
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
