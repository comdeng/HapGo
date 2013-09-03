package app

import (
	"github.com/comdeng/HapGo/hapgo/conf"
	"github.com/comdeng/HapGo/hapgo/logger"
	"github.com/comdeng/HapGo/lib/util"
	"net/http"
	"time"
)

type trackIdInfo struct {
	Create bool
	Expire uint16
	Domain string
	Path   string
	Name   string
}

func AppInitFilter(_app *WebApp) error {
	start := time.Now()

	generateTrackId(_app)

	end := time.Now()
	logger.Trace("hapgo.filter.init end cost %dμs", end.Sub(start)/1000)
	return nil
}

// 生成跟踪的ID
func generateTrackId(_app *WebApp) {
	var trackId trackIdInfo
	if ok := conf.Decode("hapgo.tid", &trackId); !ok {
		return
	}

	var value string
	cookie, err := _app.Request.Req.Cookie(trackId.Name)
	if err != nil {
		cookie := new(http.Cookie)
		cookie.Name = trackId.Name
		cookie.Path = trackId.Path
		cookie.Domain = trackId.Domain
		cookie.Expires = time.Now().Add(time.Hour * time.Duration(24*trackId.Expire))
		cookie.Value = util.UniqId()
		http.SetCookie(_app.Response.Writer, cookie)
		value = cookie.Value
	} else {
		value = cookie.Value
	}
	_app.Request.UserData["tid"] = value

	logger.AddBasic(map[string]interface{}{
		"tid": value,
	})
}
