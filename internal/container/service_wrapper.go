package container

import (
	"log"
	"reflect"
)

type serviceWrapper struct {
	singleton       bool
	cachedInstance  interface{}
	constructor     interface{}
	constructorArgs []interface{}
	instanceType    reflect.Type
}

func newServiceWrapper(fn interface{}, singleton bool, constructorArgs []interface{}) *serviceWrapper {
	fnType := reflect.TypeOf(fn)
	instanceType := fnType.Out(0)

	// Make sure func has correct arguments (*Application)
	inTypes := getTypes(constructorArgs)
	assertInTypes(fnType, inTypes, "Func arg mismatch")

	// Make sure func has correct return values (value, error)
	outTypes := []reflect.Type{AnyType, ErrType}
	assertOutTypes(fnType, outTypes, "Func return mismatch")

	// Add the dependency to the graph
	return &serviceWrapper{
		singleton:       singleton,
		cachedInstance:  nil,
		constructor:     fn,
		constructorArgs: constructorArgs,
		instanceType:    instanceType,
	}
}

func newServiceWrapperInstance(instance interface{}, singleton bool) *serviceWrapper {
	instanceType := reflect.TypeOf(instance)

	// Add the dependency to the graph
	return &serviceWrapper{
		singleton:       singleton,
		cachedInstance:  instance,
		constructor:     nil,
		constructorArgs: nil,
		instanceType:    instanceType,
	}
}

func (sw *serviceWrapper) Make() interface{} {
	// Return cached instance if it's a singleton
	if sw.singleton && sw.cachedInstance != nil {
		return sw.cachedInstance
	}

	// Create a new instance by calling the constructor
	outValues := callFunc(
		reflect.ValueOf(sw.constructor),
		getValues(sw.constructorArgs)...,
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
