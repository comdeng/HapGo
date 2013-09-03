package ctrl

import (
	"fmt"
)

type TestController struct {
	Name string
}

func (ctrl *TestController) Index() {
	fmt.Print("TestController.index")
}
