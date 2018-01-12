package hemlock

import (
	"github.com/gschier/hemlock/internal"
	"log"
	"os"
	"reflect"
)

// ~~~~~~~~~~~ //
// Application //
// ~~~~~~~~~~~ //

type Application struct {
	Config    *Config
	container *Container
}

func NewApplication(config *Config) *Application {
	app := &Application{
		Config: config,
	}

	// Ensure all service constructors take in *Application as an
	// argument
	serviceConstructorArgs := []interface{}{app}
	app.container = NewContainer(serviceConstructorArgs)

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

func (a *Application) ResolveInto(fn interface{}) []interface{} {
	return a.container.Call(fn, a)
}

func (a *Application) With(fn interface{}) []interface{} {
	return a.container.Call(fn, a)
}

func (a *Application) Make(i interface{}) interface{} {
	iType := reflect.TypeOf(i)
	internal.AssertPtrType(iType, "Cannot make non-pointer")

	var sw *ServiceWrapper
	if iType.Elem().Kind() == reflect.Interface {
		sw = a.container.findServiceWrapperByInterface(iType.Elem())
	} else {
		sw = a.container.findServiceWrapperByPtr(iType)
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

type ServiceWrapper struct {
	singleton       bool
	cachedInstance  interface{}
	constructor     ServiceConstructor
	constructorArgs []interface{}
	instanceType    reflect.Type
}

//type ServiceConstructor func(*Application) (interface{}, error)
type ServiceConstructor interface{}

func newServiceWrapper(fn ServiceConstructor, singleton bool, constructorArgs []interface{}) *ServiceWrapper {
	fnType := reflect.TypeOf(fn)
	instanceType := fnType.Out(0)

	// Make sure func has correct arguments (*Application)
	inTypes := internal.Types(constructorArgs)
	internal.AssertInTypes(fnType, inTypes, "Func arg mismatch")

	// Make sure func has correct return values (value, error)
	outTypes := []reflect.Type{internal.AnyType, internal.ErrType}
	internal.AssertOutTypes(fnType, outTypes, "Func return mismatch")

	// Add the dependency to the graph
	return &ServiceWrapper{
		singleton:       singleton,
		cachedInstance:  nil,
		constructor:     fn,
		constructorArgs: constructorArgs,
		instanceType:    instanceType,
	}
}

func newServiceWrapperInstance(instance interface{}, singleton bool) *ServiceWrapper {
	instanceType := reflect.TypeOf(instance)
	internal.AssertPtrType(instanceType, "Cannot wrap non-pointer")

	// Add the dependency to the graph
	return &ServiceWrapper{
		singleton:       singleton,
		cachedInstance:  instance,
		constructor:     nil,
		constructorArgs: nil,
		instanceType:    instanceType,
	}
}

func (sw *ServiceWrapper) Make() interface{} {
	// Return cached instance if it's a singleton
	if sw.singleton && sw.cachedInstance != nil {
		return sw.cachedInstance
	}

	// Create a new instance by calling the constructor
	outValues := internal.CallFunc(
		reflect.ValueOf(sw.constructor),
		internal.Values(sw.constructorArgs)...,
	)

	err := outValues[1].Interface()
	if err != nil {
		log.Panicf("Failed to initialize service err=%v", err)
	}

	instance := outValues[0].Interface()

	// Cache it for next time
	sw.cachedInstance = instance

	return instance
}

type Container struct {
	registered             map[reflect.Type]*ServiceWrapper
	serviceConstructorArgs []interface{}
}

func NewContainer(serviceConstructorArgs []interface{}) *Container {
	return &Container{
		registered:             make(map[reflect.Type]*ServiceWrapper),
		serviceConstructorArgs: serviceConstructorArgs,
	}
}

// Bind binds the type of v as a dependency
func (c *Container) Bind(fn ServiceConstructor) {
	w := newServiceWrapper(fn, false, c.serviceConstructorArgs)
	c.registered[w.instanceType] = w
}

// Singleton binds the type of v as a dependency. Will only get instantiated once
func (c *Container) Singleton(fn ServiceConstructor) {
	w := newServiceWrapper(fn, true, c.serviceConstructorArgs)
	c.registered[w.instanceType] = w
}

// Instance binds an already-created value as a dependency
func (c *Container) Instance(i interface{}) {
	w := newServiceWrapperInstance(i, true)
	c.registered[w.instanceType] = w
}

func (c *Container) Call(fn interface{}, app *Application) []interface{} {
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

		args[i] = reflect.ValueOf(sw.Make())
	}

	returnValues := fnValue.Call(args)
	returnInstances := make([]interface{}, len(returnValues))
	for i, rv := range returnValues {
		returnInstances[i] = rv.Interface()
	}

	return returnInstances
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
		if sw.instanceType.Kind() == reflect.Interface && ptrType.Implements(sw.instanceType) {
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
