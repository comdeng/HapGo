package app

/*
	主要用来格式化输入参数，用来防止sql注入，xss攻击等
*/

import (
	"github.com/comdeng/HapGo/hapgo/conf"
	"github.com/comdeng/HapGo/hapgo/logger"
	_html "github.com/comdeng/HapGo/lib/html"
	// "log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const trimReg = " \t\r\n\x0b\v"
const inputConfKey = "hapgo.input"

var replacer = strings.NewReplacer(
	"&", "&amp;",
	"'", "&#039;",
	`"`, "&quot;",
	">", "&gt;",
	"<", "&lt;",
	"`", "&#069;",
	"{", "&#123;",
	"\r\n", "<br/>",
	"\r", "<br/>",
	"\n", "<br/>",
)

type inputInfo struct {
	SafeCheck    bool
	PrefixOfHtml string
	PrefixOfRaw  string
	PrefixOfSafe string
}

var iInfo inputInfo

func AppInputFilter(_app *WebApp) error {
	start := time.Now()

	var once sync.Once
	once.Do(initInputConf)

	_app.Request.Req.ParseForm()
	encodeInput(_app.Request.Req)

	end := time.Now()
	logger.Trace("hapgo.filter.input end cost %dμs", end.Sub(start)/1000)
	return nil
}

func initInputConf() {
	iInfo = inputInfo{false, "ph_", "pr_", "ps_"}
	conf.Decode(inputConfKey, &iInfo)
}

func encodeInput(r *http.Request) {
	for k, v := range r.Form {
		if strings.Index(k, iInfo.PrefixOfHtml) == 0 {
			for _k, _v := range v {
				r.Form[k][_k] = filterHtml(_v)
			}
		} else if strings.Index(k, iInfo.PrefixOfRaw) == 0 {
			//
		} else if iInfo.SafeCheck && strings.Index(k, iInfo.PrefixOfSafe) == 0 {

		} else {
			for _k, _v := range v {
				r.Form[k][_k] = filterText(_v)
			}
		}

	}
}

func safeCheckText(text string) string {
	//TODY
	return text
}

func filterHtml(html string) string {
	return _html.Tidy(html)
}

func filterText(text string) string {
	strings.Trim(text, trimReg)
	return replacer.Replace(text)
}
