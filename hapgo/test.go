package main

import (
	"github.com/comdeng/HapGo/hapgo/core"
	"github.com/comdeng/HapGo/hapgo/sample/controller"
)

func main() {
	core.RegisterController(controller.DefaultController{nil})
	core.CallMethod("DefaultController", "Foo")
}
