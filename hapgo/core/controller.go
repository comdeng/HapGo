package core

import (
	"fmt"
	"github.com/comdeng/HapGo/hapgo/logger"
	"reflect"
	"strings"
	"sync"
)

type Controller struct {
	Name   string
	Method string
	Req    *HttpRequest
	Res    *HttpResponse
}

type controllerInfo struct {
	Typ     reflect.Type
	Name    string
	Methods map[string]bool
}

func (ctrl *Controller) Write(str string) {
	fmt.Fprint(ctrl.Res.Writer, str)
}

var controllers map[string]*controllerInfo = make(map[string]*controllerInfo)
var filterMethods map[string]bool
var once sync.Once

func RegisterController(ctlName string, c interface{}) {
	once.Do(func() {
		initFilterMethods(new(Controller))
	})

	typ := reflect.TypeOf(c)
	num := typ.NumMethod()
	methods := make(map[string]bool)
	for i := 0; i < num; i++ {
		name := typ.Method(i).Name
		if _, ok := filterMethods[name]; !ok {
			methods[name] = true
		}
	}

	controllers[ctlName] = &controllerInfo{
		typ,
		ctlName,
		methods,
	}
}

func initFilterMethods(ctrl *Controller) {
	typ := reflect.TypeOf(ctrl)
	num := typ.NumMethod()
	filterMethods = make(map[string]bool)
	for i := 0; i < num; i++ {
		filterMethods[typ.Method(i).Name] = true
	}
}

func NewController(ctlName string, methodName string, req *HttpRequest, res *HttpResponse) (ctrl reflect.Value) {
	ctlName = strings.ToLower(ctlName)
	var ci *controllerInfo
	var ok bool
	if ci, ok = controllers[ctlName]; !ok {
		logger.Warn("controllerNotFound controllerName=" + ctlName)
		panic("hapgo.u_notfound")
	}

	methodName = strings.Title(methodName)

	if _, ok = ci.Methods[methodName]; !ok {
		logger.Warn("methodNotFound methodName=" + methodName)
		panic("hapgo.u_notfound")
	}

	ctrl = reflect.New(ci.Typ.Elem())

	coreCtrl := &Controller{ctlName, methodName, req, res}

	// 必须将core.Controller 作为第一个参数
	ctrl.Elem().Field(0).Set(reflect.ValueOf(coreCtrl))

	return ctrl
}

func CallMethod(ctrl reflect.Value, methodName string, args ...interface{}) {
	methodName = strings.Title(methodName)

	method := ctrl.MethodByName(methodName)
	if method.IsNil() || !method.IsValid() {
		logger.Warn("hapgo.methodNotFound methodName=" + methodName)
		panic("hapgo.u_notfound")
	}
	params := make([]reflect.Value, len(args))
	for k, v := range args {
		params[k] = reflect.ValueOf(v)
	}
	method.Call(params)
}
