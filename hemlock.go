package hemlock

import (
	"log"
	"os"
	"reflect"
)

// ~~~~~~~~~~~ //
// Application //
// ~~~~~~~~~~~ //

type Application struct {
	Router    Router
	Config    *Config
	container *Container
}

func NewApplication(config *Config) *Application {
	app := &Application{
		Config:    config,
		container: NewContainer(),
	}

	// Add providers from config
	for _, p := range app.Config.Providers {
		p.Register(app.container)
	}

	// Boot all providers
	for _, p := range app.Config.Providers {
		p.Boot(app)
	}

	return app
}

func (a *Application) Bind(fn ServiceConstructor) {
	a.container.Bind(fn)
}

func (a *Application) Singleton(fn ServiceConstructor) {
	a.container.Singleton(fn)
}

func (a *Application) Instance(v interface{}) {
	a.container.Instance(v)
}

func (a *Application) Call(fn interface{}) {
	a.container.Call(fn, a)
}

func (a *Application) With(fn interface{}) {
	a.container.Call(fn, a)
}

func (a *Application) Make(i interface{}) interface{} {
	iType := reflect.TypeOf(i)
	if iType.Kind() != reflect.Ptr {
		panic("Cannot make non-pointer")
	}

	var sw *ServiceWrapper
	if iType.Elem().Kind() == reflect.Interface {
		sw = a.container.findServiceWrapperByInterface(iType.Elem())
	} else {
		sw = a.container.findServiceWrapperByPtr(iType)
	}

	return sw.Make(a)
}

func (a *Application) MakeInto(v interface{}) {
	vValue := reflect.ValueOf(v)
	vType := reflect.TypeOf(v)

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

type ServiceWrapper struct {
	singleton      bool
	cachedInstance interface{}
	constructor    ServiceConstructor
	instanceType   reflect.Type
}

//type ServiceConstructor func(*Application) (interface{}, error)
type ServiceConstructor interface{}

func newServiceWrapper(fn ServiceConstructor, singleton bool) *ServiceWrapper {
	fnType := reflect.TypeOf(fn)
	instanceType := fnType.Out(0)

	isInterface := instanceType.Kind() == reflect.Interface
	if !isInterface && instanceType.Kind() != reflect.Ptr {
		panic("Must be pointer")
	}

	// Add the dependency to the graph
	return &ServiceWrapper{
		singleton:      singleton,
		cachedInstance: nil,
		constructor:    fn,
		instanceType:   instanceType,
	}
}

func newServiceWrapperInstance(instance interface{}, singleton bool) *ServiceWrapper {
	instanceType := reflect.TypeOf(instance)
	if instanceType.Kind() != reflect.Ptr {
		panic("Must be pointer")
	}

	// Add the dependency to the graph
	return &ServiceWrapper{
		singleton:      singleton,
		cachedInstance: instance,
		constructor:    nil,
		instanceType:   instanceType,
	}
}

func (sw *ServiceWrapper) Make(app *Application) interface{} {
	if sw.singleton && sw.cachedInstance != nil {
		return sw.cachedInstance
	}

	constructorType := reflect.TypeOf(sw.constructor)
	if constructorType.Kind() != reflect.Func {
		panic("Should be func")
	}

	numIn := constructorType.NumIn()
	if numIn != 1 {
		panic("Fn must take Application as argument")
	}

	inType := constructorType.In(0)
	if inType.Kind() != reflect.Ptr {
		panic("First arg must be pointer to Application")
	}

	if inType.Elem() != reflect.TypeOf(Application{}) {
		panic("First arg must be Application")
	}

	numOut := constructorType.NumOut()
	if numOut != 2 {
		panic("Fn must return a value and error")
	}

	out2Type := constructorType.Out(1)
	errType := reflect.TypeOf((*error)(nil)).Elem()
	if !out2Type.Implements(errType) {
		log.Panicf("Fn second return value must be error type. Got %v != %v\n", out2Type, errType)
	}

	constructorValue := reflect.ValueOf(sw.constructor)
	iocValue := reflect.ValueOf(app)
	returns := constructorValue.Call([]reflect.Value{iocValue})
	instance := returns[0].Interface()
	if !returns[1].IsNil() {
		err := returns[1].Interface().(error)
		panic("Failed to initialize service err=" + err.Error())
	}

	sw.cachedInstance = instance

	return instance
}

type Container struct {
	registered map[reflect.Type]*ServiceWrapper
}

func NewContainer() *Container {
	return &Container{
		registered: make(map[reflect.Type]*ServiceWrapper),
	}
}

// Bind binds the type of v as a dependency
func (c *Container) Bind(fn ServiceConstructor) {
	w := newServiceWrapper(fn, false)
	c.registered[w.instanceType] = w
}

// Singleton binds the type of v as a dependency. Will only get instantiated once
func (c *Container) Singleton(fn ServiceConstructor) {
	w := newServiceWrapper(fn, true)
	c.registered[w.instanceType] = w
}

// Instance binds an already-created value as a dependency
func (c *Container) Instance(i interface{}) {
	w := newServiceWrapperInstance(i, true)
	c.registered[w.instanceType] = w
}

func (c *Container) Call(fn interface{}, app *Application) {
	fnType := reflect.TypeOf(fn)
	fnValue := reflect.ValueOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("Cannot provide to non-function")
	}

	// Create a Container to hold the arguments we'll call the fn with
	numArgs := fnType.NumIn()
	args := make([]reflect.Value, numArgs)

	// Build argument values one-by-one
	for i := 0; i < fnType.NumIn(); i++ {
		argType := fnType.In(i)
		var sw *ServiceWrapper
		switch argType.Kind() {
		case reflect.Interface:
			sw = c.findServiceWrapperByInterface(argType)
		case reflect.Ptr:
			sw = c.findServiceWrapperByPtr(argType)
		default:
			panic("Function argument was not pointer nor interface")
		}

		args[i] = reflect.ValueOf(sw.Make(app))
	}

	fnValue.Call(args)
}

func (c *Container) findServiceWrapperByInterface(iType reflect.Type) *ServiceWrapper {
	if iType.Kind() != reflect.Interface {
		panic("Argument type must be an interface")
	}

	var matchedSW *ServiceWrapper
	leastMethods := -1
	for _, sw := range c.registered {
		//fmt.Printf("Checking %v =? %v\n", iType, sw.instanceType)
		numMethods := sw.instanceType.NumMethod()

		if leastMethods != -1 && numMethods > leastMethods {
			continue
		}

		if !sw.instanceType.Implements(iType) {
			continue
		}

		leastMethods = numMethods
		matchedSW = sw
	}

	if matchedSW == nil {
		log.Panicf("Type not found for arg: %v", iType)
	}

	return matchedSW
}

func (c *Container) findServiceWrapperByPtr(ptrType reflect.Type) *ServiceWrapper {
	if ptrType.Kind() != reflect.Ptr {
		panic("Argument type must be an pointer")
	}

	for _, sw := range c.registered {
		//fmt.Printf("Checking Ptr %v =? %v\n", ptrType, sw.instanceType)
		// TODO: Find best match interface
		if sw.instanceType.Kind() == reflect.Interface && ptrType.Implements(sw.instanceType){
			return sw
		}

		if sw.instanceType.AssignableTo(ptrType) {
			return sw
		}

	}

	log.Panicf("Type not found for arg: %v", ptrType)
	return nil
}

func Env(name string) string {
	return os.Getenv(name)
}

func EnvOr(name, fallback string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}

	return value
}

type Providers []Provider

type Provider interface {
	// Register registers a new provider. Any setup should happen here
	Register(*Container)

	// Boot is called after all service providers have been registered
	Boot(*Application)
}
