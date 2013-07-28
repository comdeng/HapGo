package app

import (
	"fmt"
	"github.com/comdeng/HapGo/hapgo/conf"
	"github.com/comdeng/HapGo/hapgo/logger"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var appId uint64

type WebApp struct {
}

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
	conf.Load("hapgo.conf")

	// 初始化日志
	logger.Init(logDir)
}

func (app *WebApp) Execute(w http.ResponseWriter, r *http.Request) {
	appId := createAppId()
	logger.SetAppId(appId)

	// 开始执行filter
	logger.Debug("hapgo.filter.init")
	log.Print(r.RequestURI)
	log.Print(r.URL)

	fmt.Fprintf(w, "uu", "a")
}

func (app *WebApp) AppId() uint64 {
	return appId
}

func createAppId() uint64 {
	now := time.Now()
	timeStamp := uint64(now.Unix()*100) + uint64(now.Nanosecond()/100000)
	log.Print(timeStamp)
	rand := uint64(rand.Float64() * 2 * float64(timeStamp))
	log.Print(rand)
	id := (int64(timeStamp) ^ int64(rand)) & 0xFFFFFFFF
	log.Print(id)
	return uint64(math.Floor(float64(id)/100) * 100)
}
