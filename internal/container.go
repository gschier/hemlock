package internal

import (
	"fmt"
	"reflect"
)

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
func (c *Container) Bind(fn interface{}) {
	w := newServiceWrapper(fn, false, c.serviceConstructorArgs)
	c.registered[w.instanceType] = w
}

// Singleton binds the type of v as a dependency. Will only get instantiated once
func (c *Container) Singleton(fn interface{}) {
	w := newServiceWrapper(fn, true, c.serviceConstructorArgs)
	c.registered[w.instanceType] = w
}

// Instance binds an already-created value as a dependency
func (c *Container) Instance(i interface{}) {
	w := newServiceWrapperInstance(i, true)
	c.registered[w.instanceType] = w
}

func (c *Container) Call(fn interface{}, extraArgs []interface{}) []interface{} {
	fnType := reflect.TypeOf(fn)
	fnValue := reflect.ValueOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("Cannot provide to non-function")
	}

	// Create a Container to hold the arguments we'll call the fn with
	numArgsToFill := fnType.NumIn() - len(extraArgs)
	filledArgs := make([]reflect.Value, numArgsToFill)

	// Build argument values one-by-one
	for i := 0; i < numArgsToFill; i++ {
		argType := fnType.In(i)
		var sw *ServiceWrapper
		switch argType.Kind() {
		case reflect.Interface:
			sw = c.FindServiceWrapperByInterface(argType)
		case reflect.Ptr:
			sw = c.FindServiceWrapperByPtr(argType)
		default:
			sw = c.FindServiceWrapperByValue(argType)
		}

		filledArgs[i] = reflect.ValueOf(sw.Make())
	}

	allArgs := append(filledArgs, Values(extraArgs)...)
	returnValues := fnValue.Call(allArgs)
	returnInstances := make([]interface{}, len(returnValues))
	for i, rv := range returnValues {
		returnInstances[i] = rv.Interface()
	}

	return returnInstances
}

func (c *Container) FindServiceWrapperByInterface(iType reflect.Type) *ServiceWrapper {
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

	if matchedSW != nil {
		return matchedSW
	}

	return nil
}

func (c *Container) FindServiceWrapperByPtr(ptrType reflect.Type) *ServiceWrapper {
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

	return nil
}

func (c *Container) FindServiceWrapperByValue(valueType reflect.Type) *ServiceWrapper {
	for _, sw := range c.registered {
		if sw.instanceType.AssignableTo(valueType) {
			return sw
		}
	}

	return nil
}
