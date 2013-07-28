package app

import (
	//"github.com/comdeng/HapGo/hapgo/conf"
	//"log"
	//"path/filepath"
	"net/http"
)

type Application interface {
	Init()
	//http.HandlerFunc
	Execute(w http.ResponseWriter, req *http.Request)
	AppId() uint64
}

const (
	MODE_WEB    = "web"
	MODE_TOOL   = "tool"
	MODE_SOCKET = "socket"
)

func NewWebApp() Application {
	var app Application = new(WebApp)
	return app
}
