package app

import (
	// "fmt"
	"github.com/comdeng/HapGo/hapgo/conf"
	"github.com/comdeng/HapGo/hapgo/core"
	"github.com/comdeng/HapGo/hapgo/logger"
	//"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var appId uint64

type WebApp struct {
	Request  *core.HttpRequest
	Response *core.HttpResponse
}

var redirectConfs = make(map[string]string)

const (
	OUTPUT_FORMAT_KEY   = "_of"
	OUTPUT_ENCODING_KEY = "_oe"
	DEBUG_KEY           = "_d"
	redirectConfKey     = "hapgo.redirect"
)

// type WebFilter interface {
// 	Execute(WebApp *app)
// }

func (_app *WebApp) Init() {
	paths, ok := conf.Get("hapgo.dirs")
	var confDir string
	var logDir string

	if !ok {
		panic("webapp.dirsNotDefined")
	}

	ps := paths.(map[string]string)
	confDir, ok = ps["conf"]
	if !ok {
		panic("webapp.confDirNotDefined")
	}

	logDir, ok = ps["log"]
	if !ok {
		panic("webapp.logDirNotDefined")
	}

	// 初始化配置文件
	conf.Init(confDir)
	err := conf.Load("hapgo.conf")
	if err != nil {
		panic("webapp.confLoadError " + err.Error())
	}

	// 初始化日志
	logger.Init(logDir)

	confs, ok := conf.Get(redirectConfKey)
	for k, v := range confs.(map[string]interface{}) {
		redirectConfs[k] = v.(string)
	}
}

func (_app *WebApp) Execute(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal("hapgo.webappFatal:%v", r)
			if err, ok := r.(string); ok {
				arr := strings.SplitN(err, " ", 2)
				if url, ok := redirectConfs[arr[0]]; ok {
					w.Header().Set("Location", url)
					switch arr[0] {
					case "hapgo.u_notfound":
						w.WriteHeader(http.StatusNotFound)
					case "hapgo.u_login":
						w.WriteHeader(http.StatusTemporaryRedirect)
					case "hapgo.u_error":
						w.WriteHeader(http.StatusInternalServerError)
					case "hapgo.u_power":
						w.WriteHeader(http.StatusForbidden)
					}
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					if url, ok := redirectConfs["hapgo.u_error"]; ok {
						w.Header().Set("Location", url)
					}
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				if url, ok := redirectConfs["hapgo.u_error"]; ok {
					w.Header().Set("Location", url)
				}
			}
		}
	}()

	appId := createAppId()
	logger.SetAppId(appId)

	req := new(core.HttpRequest)
	req.Init(r)
	_app.Request = req

	res := new(core.HttpResponse)
	res.Init(w)
	_app.Response = res

	_app.executeFilter("init", w)
	_app.executeFilter("input", w)

	_app.handlReqAndRes()

	_app.executeFilter("url", w)

	// 开始执行filter
	// logger.Debug("hapgo.filter.init")
	// log.Print(r.RequestURI)
	// log.Print(r.URL)

	//fmt.Fprintf(w, "you URI:%s, URL:%s", r.RequestURI, _app.Request.UserData["tid"].(string))
}

func (_app *WebApp) handlReqAndRes() {
	if _app.Request.Req.Method == "POST" {
		_app.Response.SetFormat(core.FORMAT_JSON)
	}
	if outputFormat, ok := _app.Request.Get(OUTPUT_FORMAT_KEY); ok {
		_app.Response.SetFormat(outputFormat)
	}
	if outputEncoding, ok := _app.Request.Get(OUTPUT_ENCODING_KEY); ok {
		_app.Response.SetEncoding(outputEncoding)
	}
}

func (_app *WebApp) executeFilter(filterName string, w http.ResponseWriter) {
	err := InitFilter(filterName, _app)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (app *WebApp) AppId() uint64 {
	return appId
}

func createAppId() uint64 {
	now := time.Now()
	timeStamp := uint64(now.Unix()*100) + uint64(now.Nanosecond()/100000)
	//log.Print(timeStamp)
	rand := uint64(rand.Float64() * 2 * float64(timeStamp))
	//log.Print(rand)
	id := (int64(timeStamp) ^ int64(rand)) & 0xFFFFFFFF
	//log.Print(id)
	return uint64(math.Floor(float64(id)/100) * 100)
}
