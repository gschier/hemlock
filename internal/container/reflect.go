package container

import (
	"log"
	"reflect"
)

var ErrType = reflect.TypeOf((*error)(nil)).Elem()
var AnyType = reflect.TypeOf(nil)

func fnArgTypes(fn reflect.Type) []reflect.Type {
	assertFuncType(fn, "Cannot get arguments of non-function")
	fnType := unPtrType(fn)
	numArgs := fnType.NumIn()
	types := make([]reflect.Type, numArgs)
	for i := 0; i < numArgs; i++ {
		types[i] = fnType.In(i)
	}
	return types
}

func assertNumIn(fn reflect.Type, n int, panicMsg string) {
	assertFuncType(fn, "Cannot assert num in of non-function")
	numIn := fn.NumIn()
	if numIn != n {
		log.Panicf("%v num in %d != %d\n", panicMsg, numIn, n)
	}
}

func assertInTypes(fn reflect.Type, types []reflect.Type, panicMsg string) {
	assertFuncType(fn, "Cannot assert in of non-function")
	assertNumIn(fn, len(types), panicMsg)
	fnTypes := fnArgTypes(fn)
	for i, t := range fnTypes {
		if types[i] == AnyType {
			continue
		}
		if !t.AssignableTo(types[i]) {
			log.Panicf("%v %v != %v\n", panicMsg, t, types[i])
		}
	}
}

func assertNumOut(fn reflect.Type, n int, panicMsg string) {
	assertFuncType(fn, "Cannot assert num out of non-function")
	numOut := fn.NumOut()
	if numOut != n {
		log.Panicf("%v num out %d != %d\n", panicMsg, numOut, n)
	}
}

func assertOutTypes(fn reflect.Type, types []reflect.Type, panicMsg string) {
	assertFuncType(fn, "Cannot assert out of non-function")
	assertNumOut(fn, len(types), panicMsg)
	fnTypes := fnArgTypes(fn)
	for i, t := range fnTypes {
		if types[i] == AnyType {
			continue
		}
		if !t.AssignableTo(types[i]) {
			log.Panicf("%v %v != %v\n", panicMsg, t, types[i])
		}
	}
}

func callFunc(fn reflect.Value, args ...reflect.Value) []reflect.Value {
	assertFuncType(fn.Type(), "Cannot call non-func")
	return fn.Call(args)
}

func funcReturnTypes(fn reflect.Type) []reflect.Type {
	assertFuncType(fn, "Cannot get return types of non-function")
	fnType := unPtrType(fn)
	numReturn := fnType.NumOut()
	types := make([]reflect.Type, numReturn)
	for i := 0; i < numReturn; i++ {
		types[i] = fnType.Out(i)
	}
	return types
}

func assertFuncType(fn reflect.Type, panicMsg string) {
	if fn.Kind() != reflect.Func {
		panic(panicMsg)
	}
}

func getTypes(v []interface{}) []reflect.Type {
	types := make([]reflect.Type, len(v))
	for i, v := range v {
		types[i] = reflect.TypeOf(v)
	}
	return types
}

func getValues(v []interface{}) []reflect.Value {
	values := make([]reflect.Value, len(v))
	for i, v := range v {
		values[i] = reflect.ValueOf(v)
	}
	return values
}

func assertPtrType(t reflect.Type, panicMsg string) {
	if t.Kind() != reflect.Ptr {
		panic(panicMsg)
	}
}

// UnPtr recursively returns the value that type points to (if any)
func unPtrType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func getTypeAndValue(v interface{}) (reflect.Type, reflect.Value) {
	return reflect.TypeOf(v), reflect.ValueOf(v)
}
