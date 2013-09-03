package app

import (
	"github.com/comdeng/HapGo/hapgo/conf"
	"github.com/comdeng/HapGo/hapgo/core"
	"github.com/comdeng/HapGo/hapgo/logger"
	"github.com/comdeng/HapGo/lib/cache"
	"github.com/comdeng/HapGo/lib/util"
	"strings"
	"sync"
	"time"
)

type csrfInfo struct {
	Enable   bool
	Url      string
	Expire   int    // 过期时间，单位为s
	TokenKey string // 字段名称
}

const csrfConfKey = "hapgo.csrf"

var _csrfInfo csrfInfo

const (
	pathSplitter = '/'
	postFlag     = '_'
)

func AppUrlFilter(_app *WebApp) error {
	start := time.Now()

	var once sync.Once
	once.Do(initUrlConf)

	if !initCsrf(_app) {
		initUrl(_app)
	}

	end := time.Now()
	logger.Trace("hapgo.filter.url end cost %dμs", end.Sub(start)/1000)
	return nil
}

func initUrlConf() {
	_csrfInfo = csrfInfo{
		false,
		"",
		0,
		"",
	}
	conf.Decode(csrfConfKey, &_csrfInfo)
}

func initUrl(_app *WebApp) {
	req := _app.Request.Req
	url := strings.ToLower(req.URL.Path)
	url = url[1:]
	urlLen := len(url)
	if url[urlLen-1] == pathSplitter {
		url = url[0 : urlLen-2]
	}
	arr := strings.Split(url, string(pathSplitter))
	// 限制POST请求必须以_开头
	if arr[len(arr)-1][0] == postFlag {
		if req.Method != "POST" {
			panic("hapgo.u_notfound")
		}
	}

	if len(arr) == 2 {
		ctrl := core.NewController(arr[0], arr[1], _app.Request, _app.Response)

		core.CallMethod(ctrl, arr[1])
	}
}

// csrf 校验
func initCsrf(_app *WebApp) (isOver bool) {
	req := _app.Request.Req
	if req.Method == "POST" {
		if _csrfInfo.Enable {
			if _csrfInfo.Url == req.URL.Path {
				csrfId := util.UniqId()
				cache.Set("csrf."+csrfId, _app.Request.UserData["tid"], 10)
				_app.Response.Set(_csrfInfo.TokenKey, csrfId)
				_app.Response.Send()
				return true
			} else {
				tid := _app.Request.UserData["tid"]
				csrfId, ok := _app.Request.Get(_csrfInfo.TokenKey)
				if !ok {
					panic("hapn.u_input csrf param is not defined")
				}
				if len(csrfId) != 32 {
					panic("hapgo.u_csrf tokenIllegal")
				}
				if ctid, ok := cache.Get("csrf." + csrfId); ok {
					if ctid != tid {
						panic("hapgo.u_csrf tokenAndTidNotMatched")
					}
				} else {
					panic("hapgo.u_csrf tokenIsMissed")
				}
			}
		}
	}
	return false
}
