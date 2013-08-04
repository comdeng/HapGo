package app

import (
	"github.com/comdeng/HapGo/hapgo/conf"
	"github.com/comdeng/HapGo/hapgo/logger"
	"github.com/comdeng/HapGo/lib/cache"
	"github.com/comdeng/HapGo/lib/util"
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

func AppUrlFilter(_app *WebApp) error {
	start := time.Now()

	var once sync.Once
	once.Do(initUrlConf)

	if initCsrf(_app) {

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
