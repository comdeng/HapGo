package cache

import (
	"time"
)

type cacheInfo struct {
	Value  interface{}
	Expire uint32 // 过期时间 单位：s
	Create time.Time
}

// 最多存储7天时间
const (
	_MAX_CACHE_TIME     = 24 * 3600 * 7
	_DEFAULT_CACHE_TIME = 3600
)

var caches = make(map[string]cacheInfo)

// 设置指定的缓存
func Set(key string, value interface{}, args ...uint32) {
	expire := uint32(0)
	if len(args) < 1 {
		expire = _DEFAULT_CACHE_TIME
	} else {
		expire = args[0]
	}
	if expire > _MAX_CACHE_TIME {
		expire = _MAX_CACHE_TIME
	}

	caches[key] = cacheInfo{
		value,
		expire,
		time.Now(),
	}
}

// 获取某个缓存
func Get(key string) (value interface{}, ok bool) {
	ci, ok := caches[key]
	if !ok {
		return nil, false
	}
	if time.Now().Sub(ci.Create) > (time.Second * time.Duration(ci.Expire)) {
		Del(key)
		return nil, false
	}
	return ci.Value, true
}

// 删除指定的缓存
func Del(key string) {
	delete(caches, key)
}
