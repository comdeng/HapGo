package conf

import (
	"encoding/json"
	//"log"
	"os"
	"path/filepath"
)

var confs = make(map[string]interface{})
var confDir string

func Init(dir string) {
	confDir = dir
}

// 载入配置文件
func Load(name string) error {
	fo, err := os.Open(filepath.Join(confDir, name))
	if err != nil {
		return err
	}
	defer fo.Close()

	decoder := json.NewDecoder(fo)
	data := make(map[string]interface{})
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}
	// 覆盖掉原有的配置
	for k, v := range data {
		confs[k] = v
	}
	return nil
}

// 获取指定的键
func Get(key string) (value interface{}, ok bool) {
	value, ok = confs[key]
	return
}

// 设置配置项
func Set(key string, value interface{}) {
	confs[key] = value
}

// 删除指定项
func Remove(key string) {
	delete(confs, key)
}
