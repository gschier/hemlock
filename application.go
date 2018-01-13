package hemlock

import (
	"context"
	"github.com/gschier/hemlock/internal"
	"reflect"
)

type Application struct {
	Config    *Config
	container *internal.Container
	ctx       context.Context
}

func NewApplication(config *Config) *Application {
	app := &Application{
		Config: config,
	}

	// Ensure all service constructors take in *Application as an
	// argument
	serviceConstructorArgs := []interface{}{app}
	app.container = internal.NewContainer(serviceConstructorArgs)

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
	iType := reflect.TypeOf(i)
	internal.AssertPtrType(iType, "Cannot make non-pointer")

	var sw *internal.ServiceWrapper
	if iType.Elem().Kind() == reflect.Interface {
		sw = a.container.FindServiceWrapperByInterface(iType.Elem())
	} else {
		sw = a.container.FindServiceWrapperByPtr(iType)
	}

	return sw.Make()
}

func (a *Application) Resolve(v interface{}) {
	vType, vValue := internal.TypeAndValue(v)

	instance := a.Make(v)
	instanceValue := reflect.ValueOf(instance)

	if vType.Elem().Kind() == reflect.Interface {
		vValue.Elem().Set(instanceValue)
	} else {
		vValue.Elem().Set(instanceValue.Elem())
	}
}

func (a *Application) Env(name string) string {
	return Env(name)
}

func (a *Application) EnvOr(name, fallback string) string {
	return EnvOr(name, fallback)
}
