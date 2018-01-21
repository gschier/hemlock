package hemlock

import (
	"context"
	"fmt"
	"github.com/gschier/hemlock/internal/container"
	"log"
	"os"
	"path/filepath"
	"reflect"
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

	// Add config Extra
	for _, c := range app.Config.Extra {
		app.Instance(c)
	}

	// Add providers from config
	for _, p := range providers {
		name := reflect.TypeOf(p).Elem().Name()
		p.Register(app.container)
		fmt.Printf("[app] Registered %s\n", name)
	}

	// Boot all providers
	for _, p := range providers {
		err := p.Boot(app)
		name := reflect.TypeOf(p).Elem().Name()
		if err != nil {
			log.Panicf("Failed to boot %s: %v\n", name, err)
		}
		fmt.Printf("[app] Booted %s\n", name)
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

func (a *Application) ResolveDir(elem ...string) string {
	cwd, _ := os.Getwd()
	newElem := make([]string, len(elem) + 1)
	newElem[0] = cwd
	for i, e := range elem {
		newElem[i+1] = e
	}
	return filepath.Join(newElem...)
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
