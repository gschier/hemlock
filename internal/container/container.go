package container

import (
	"log"
	"reflect"
)

type Container struct {
	registered             []*serviceWrapper
	serviceConstructorArgs []interface{}
}

func New(serviceConstructorArgs []interface{}) *Container {
	return &Container{
		registered:             make([]*serviceWrapper, 0),
		serviceConstructorArgs: serviceConstructorArgs,
	}
}

func Clone(c *Container, serviceConstructorArgs []interface{}) *Container {
	newContainer := New(serviceConstructorArgs)
	copy(newContainer.registered, c.registered)
	return newContainer
}

// Bind binds the type of v as a dependency
func (c *Container) Bind(fn interface{}) {
	w := newServiceWrapper(fn, false, c.serviceConstructorArgs)
	c.registered = append(c.registered, w)
}

// Singleton binds the type of v as a dependency. Will only get instantiated once
func (c *Container) Singleton(fn interface{}) {
	w := newServiceWrapper(fn, true, c.serviceConstructorArgs)
	c.registered = append(c.registered, w)
}

// Instance binds an already-created value as a dependency
func (c *Container) Instance(i interface{}) {
	w := newServiceWrapperInstance(i, true)
	c.registered = append(c.registered, w)
}

func (c *Container) Make(i interface{}) interface{} {
	iType := reflect.TypeOf(i)
	if iType.Kind() != reflect.Ptr {
		panic("Cannot make non-pointer")
	}

	var sw *serviceWrapper
	if iType.Elem().Kind() == reflect.Interface {
		sw = c.findServiceWrapperByInterface(iType.Elem())
	} else {
		sw = c.findServiceWrapperByPtr(iType)
	}

	return sw.Make()
}

func (c *Container) Resolve(v interface{}) {
	vType, vValue := getTypeAndValue(v)

	instance := c.Make(v)
	instanceValue := reflect.ValueOf(instance)

	if vType.Elem().Kind() == reflect.Interface {
		vValue.Elem().Set(instanceValue)
	} else if vType.Kind() == reflect.Ptr && !vValue.Elem().IsValid() {
		log.Panicf("Cannot resolve into zero-value pointer %#v\n", vValue)
	} else {
		vValue.Elem().Set(instanceValue.Elem())
	}
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
		var sw *serviceWrapper
		switch argType.Kind() {
		case reflect.Interface:
			sw = c.findServiceWrapperByInterface(argType)
		case reflect.Ptr:
			sw = c.findServiceWrapperByPtr(argType)
		default:
			sw = c.FindServiceWrapperByValue(argType)
		}

		if sw == nil {
			log.Panicf("Failed to find correct type %v\n", argType)
		}

		filledArgs[i] = reflect.ValueOf(sw.Make())
	}

	allArgs := append(filledArgs, getValues(extraArgs)...)
	returnValues := fnValue.Call(allArgs)
	returnInstances := make([]interface{}, len(returnValues))
	for i, rv := range returnValues {
		returnInstances[i] = rv.Interface()
	}

	return returnInstances
}

func (c *Container) findServiceWrapperByInterface(iType reflect.Type) *serviceWrapper {
	if iType.Kind() != reflect.Interface {
		panic("Argument type must be an interface")
	}

	var matchedSW *serviceWrapper
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
		log.Panicf("Could not resolve anything for interface %v out of %v\n", iType, c.registered)
	}

	return matchedSW
}

func (c *Container) findServiceWrapperByPtr(ptrType reflect.Type) *serviceWrapper {
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

	log.Panicf("Could not resolve anything for ptr %v out of %v\n", ptrType, c.registered)
	return nil
}

func (c *Container) FindServiceWrapperByValue(valueType reflect.Type) *serviceWrapper {
	for _, sw := range c.registered {
		if sw.instanceType.AssignableTo(valueType) {
			return sw
		}
	}

	return nil
}
