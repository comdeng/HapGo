package main

import (
	_ctrl "github.com/comdeng/HapGo/test/ctrl"
	"reflect"

	"fmt"
)

func main() {
	var ctrl interface{} = new(_ctrl.TestController)
	//ctrlValue := reflect.ValueOf(ctrl)

	//ctrlValue.Call(nil)

	fmt.Println(ctrl)

	ctrlType := reflect.TypeOf(ctrl)

	ctrl1 := reflect.New(ctrlType.Elem())
	fmt.Println(ctrl1)

	method := ctrl1.MethodByName("Index")
	method.Call(nil)
}
