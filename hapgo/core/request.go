package core

import (
	"net/http"
)

type HttpRequest struct {
	Req      *http.Request
	UserData map[string]interface{}
}

func (request *HttpRequest) Init(req *http.Request) {
	request.Req = req
	request.UserData = map[string]interface{}{}
}

// 获取请求的值
func (request *HttpRequest) Get(key string) (value string, ok bool) {
	values := request.Req.Form[key]
	if len(values) == 0 {
		return "", false
	} else {
		return values[0], true
	}
}

func (request *HttpRequest) Gets(key string) (values []string, ok bool) {
	vals := request.Req.Form[key]
	if len(vals) == 0 {
		return nil, false
	} else {
		values = make([]string, len(vals))
		for k, v := range vals {
			values[k] = v
		}
		return values, true
	}
}

// 获取指定的post字段
func (request *HttpRequest) Post(key string) (value interface{}, ok bool) {
	return getResponse(request.Req.PostForm[key])
}

func (request *HttpRequest) Posts(key string) (values []string, ok bool) {
	vals := request.Req.PostForm[key]
	if len(vals) == 0 {
		return nil, false
	} else {
		values = make([]string, len(vals))
		for k, v := range vals {
			values[k] = v
		}
		return values, true
	}
}

// 获取指定的cookie值
func (request *HttpRequest) Cookie(key string) (value interface{}, ok bool) {
	cookie, err := request.Req.Cookie(key)
	if err != nil {
		return nil, false
	}
	return cookie.Value, true
}

func getResponse(values []string) (value interface{}, ok bool) {
	switch len(values) {
	case 0:
		return nil, false
	case 1:
		return values[0], true
	default:
		return values, true
	}
}
