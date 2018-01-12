package internal

import (
	"log"
	"reflect"
)

var ErrType = reflect.TypeOf((*error)(nil)).Elem()
var AnyType = reflect.TypeOf(nil)

func FnArgTypes(fn reflect.Type) []reflect.Type {
	AssertFuncType(fn, "Cannot get arguments of non-function")
	fnType := UnPtrType(fn)
	numArgs := fnType.NumIn()
	types := make([]reflect.Type, numArgs)
	for i := 0; i < numArgs; i++ {
		types[i] = fnType.In(i)
	}
	return types
}

func AssertNumIn(fn reflect.Type, n int, panicMsg string) {
	AssertFuncType(fn, "Cannot assert num in of non-function")
	numIn := fn.NumIn()
	if numIn != n {
		log.Panicf("%v num in %d != %d\n", panicMsg, numIn, n)
	}
}

func AssertInTypes(fn reflect.Type, types []reflect.Type, panicMsg string) {
	AssertFuncType(fn, "Cannot assert in of non-function")
	AssertNumIn(fn, len(types), panicMsg)
	fnTypes := FnArgTypes(fn)
	for i, t := range fnTypes {
		if types[i] == AnyType {
			continue
		}
		if !t.AssignableTo(types[i]) {
			log.Panicf("%v %v != %v\n", panicMsg, t, types[i])
		}
	}
}

func AssertNumOut(fn reflect.Type, n int, panicMsg string) {
	AssertFuncType(fn, "Cannot assert num out of non-function")
	numOut := fn.NumOut()
	if numOut != n {
		log.Panicf("%v num out %d != %d\n", panicMsg, numOut, n)
	}
}

func AssertOutTypes(fn reflect.Type, types []reflect.Type, panicMsg string) {
	AssertFuncType(fn, "Cannot assert out of non-function")
	AssertNumOut(fn, len(types), panicMsg)
	fnTypes := FnArgTypes(fn)
	for i, t := range fnTypes {
		if types[i] == AnyType {
			continue
		}
		if !t.AssignableTo(types[i]) {
			log.Panicf("%v %v != %v\n", panicMsg, t, types[i])
		}
	}
}

func CallFunc(fn reflect.Value, args ...reflect.Value) []reflect.Value {
	AssertFuncType(fn.Type(), "Cannot call non-func")
	return fn.Call(args)
}

func FnReturnTypes(fn reflect.Type) []reflect.Type {
	AssertFuncType(fn, "Cannot get return types of non-function")
	fnType := UnPtrType(fn)
	numReturn := fnType.NumOut()
	types := make([]reflect.Type, numReturn)
	for i := 0; i < numReturn; i++ {
		types[i] = fnType.Out(i)
	}
	return types
}

func AssertFuncType(fn reflect.Type, panicMsg string) {
	if fn.Kind() != reflect.Func {
		panic(panicMsg)
	}
}

func AssertNotInterfaceType(t reflect.Type, panicMsg string) {
	if t.Kind() == reflect.Interface {
		panic(panicMsg)
	}
}

func Types(v []interface{}) []reflect.Type {
	types := make([]reflect.Type, len(v))
	for i, v := range v {
		types[i] = reflect.TypeOf(v)
	}
	return types
}

func Values(v []interface{}) []reflect.Value {
	values := make([]reflect.Value, len(v))
	for i, v := range v {
		values[i] = reflect.ValueOf(v)
	}
	return values
}

func AssertPtrType(t reflect.Type, panicMsg string) {
	if t.Kind() != reflect.Ptr {
		panic(panicMsg)
	}
}

// UnPtr recursively returns the value that type points to (if any)
func UnPtrType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func UnPtrValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

func TypeAndValue(v interface{}) (reflect.Type, reflect.Value) {
	return reflect.TypeOf(v), reflect.ValueOf(v)
}
