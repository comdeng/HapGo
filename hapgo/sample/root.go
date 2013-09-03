package main

import (
	"github.com/comdeng/HapGo/hapgo/app"
	"github.com/comdeng/HapGo/hapgo/conf"
	"github.com/comdeng/HapGo/hapgo/core"
	"github.com/comdeng/HapGo/hapgo/logger"
	"github.com/comdeng/HapGo/hapgo/sample/controller"
	//"log"
	// "fmt"
	"net/http"
	"os"
	"path/filepath"
)

var (
	// 根目录
	RootDir string
	_app    app.Application
)

func main() {
	RootDir, _ := os.Getwd()

	defer func() {
		logger.Release()
	}()

	conf.Set("hapgo.dirs", map[string]string{
		"root": RootDir,
		"conf": filepath.Join(RootDir, "conf"),
		"log":  filepath.Join(RootDir, "log"),
		"tmp":  filepath.Join(RootDir, "tmp"),
		"app":  filepath.Join(RootDir, "app"),
		"page": filepath.Join(RootDir, "page"),
	})

	_app = app.NewWebApp()
	_app.Init()

	// 注册控制器
	core.RegisterController("default", &controller.DefaultController{nil})

	http.HandleFunc("/", handle)
	http.ListenAndServe("127.0.0.1:11000", nil)

}

func handle(rw http.ResponseWriter, req *http.Request) {
	_app.Execute(rw, req)
}
