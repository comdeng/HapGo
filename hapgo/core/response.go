package core

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

const (
	FORMAT_JSON = "json"
	FORMAT_HTML = "html"
	FORMAT_TEXT = "text"
	FORMAT_XML  = "xml"
	STATUS_OK   = "hapgo.ok"
)

type responseData struct {
	Err  string
	Data map[string]interface{}
}

type HttpResponse struct {
	Res         http.ResponseWriter
	outFormat   string
	outEncoding string
	data        responseData
	rawData     string
	template    string
}

const (
	dateLayout     = "Mon, 02 Jan 2006 15:04:05 GMT"
	contentTypeKey = "Content-Type"
)

func (response *HttpResponse) Init(res http.ResponseWriter) {
	response.Res = res
	response.outEncoding = "UTF-8"
	response.data.Data = make(map[string]interface{})
	response.data.Err = STATUS_OK
}

func (response *HttpRequest) SetNoCache() {
	h := response.Req.Header
	h.Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
	h.Set("Last-Modified", time.Now().Format(dateLayout))
	h.Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	h.Set("Pragma", "no-cache")
}

func (response *HttpResponse) SetP3P() {
	//TODO
}

func BuildContentType(w http.ResponseWriter, format string, encoding string) {
	switch format {
	case FORMAT_JSON:
		w.Header().Set(contentTypeKey, "application/json; charset="+encoding)
	case FORMAT_HTML:
		w.Header().Set(contentTypeKey, "text/html; charset="+encoding)
	case FORMAT_XML:
		w.Header().Set(contentTypeKey, "text/xml; charset="+encoding)
	default:
		w.Header().Set(contentTypeKey, "text/plain; charset="+encoding)
	}
}

// 构建视图
func BuildView(template string, data map[string]interface{}) string {
	return "template"
}

func (res *HttpResponse) Set(key string, data interface{}) {
	res.data.Data[key] = data
}

func (res *HttpResponse) SetRaw(rawData string) {
	res.rawData = rawData
}

func (res *HttpResponse) SetFormat(format string) {
	if format == FORMAT_JSON || format == FORMAT_HTML || format == FORMAT_TEXT || format == FORMAT_XML {
		res.outFormat = format
	}
}

func (res *HttpResponse) SetEncoding(encoding string) {
	res.outEncoding = encoding
}

func (res *HttpResponse) getContent() string {
	BuildContentType(res.Res, res.outFormat, res.outEncoding)

	if res.rawData != "" {
		return res.rawData
	} else if res.template != "" {
		return BuildView(res.template, res.data.Data)
	}

	switch res.outFormat {
	case FORMAT_JSON:
		b, err := json.Marshal(res.data)
		if err != nil {
			panic("hapgo.responseError:" + err.Error())
		}
		return fmt.Sprintf("%s", b)
	case FORMAT_HTML:
		return ""
	case FORMAT_XML:
		b, err := xml.Marshal(res.data)
		if err != nil {
			panic("hapgo.responseError:" + err.Error())
		}
		return fmt.Sprintf("%s", b)
	default:
		return ""
	}
}

// 发送数据
func (res *HttpResponse) Send() {
	data := res.getContent()
	if len(data) > 0 {
		res.Res.Write([]byte(data))
	}
}
