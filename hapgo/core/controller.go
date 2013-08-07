package core

import (
	"log"
	"reflect"
)

type Controller struct {
}

var controllers map[string]reflect.Type = make(map[string]reflect.Type)

func RegisterController(c interface{}) {
	t := reflect.TypeOf(c)

	controllers[t.Name()] = t
}

func CallMethod(controllerName, methodName string, args ...interface{}) {
	var t reflect.Type
	var ok bool
	if t, ok = controllers[controllerName]; !ok {
		panic("hapgo.controllerNotFound controllerName=" + controllerName)
	}

	ctl := reflect.New(t)
	typ := reflect.Value(ctl).Elem()
	method := typ.MethodByName(methodName)
	// if method == nil {
	// 	panic("hapgo.methodNotFound methodName=" + methodName)
	// }
	params := make([]reflect.Value, len(args))
	for k, v := range args {
		params[k] = reflect.ValueOf(v)
	}
	log.Print(method, params)
	method.Call(params)
}
