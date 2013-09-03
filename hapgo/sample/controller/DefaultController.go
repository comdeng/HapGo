package controller

import (
	"github.com/comdeng/HapGo/hapgo/core"
)

type DefaultController struct {
	*core.Controller
}

func (ctrl *DefaultController) Index() {
	ctrl.Write("hello,world!")
}

func (ctrl *DefaultController) Test() {
	ctrl.Write("Test")
}
