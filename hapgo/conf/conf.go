package conf

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

var confs = make(map[string]interface{})
var confDir string

func Init(dir string) {
	confDir = dir
}

func All() map[string]interface{} {
	return confs
}

// 载入配置文件
func Load(name string) error {
	if confDir != "" {
		name = filepath.Join(confDir, name)
	}
	fo, err := os.Open(name)
	if err != nil {
		return err
	}
	defer fo.Close()

	decoder := json.NewDecoder(fo)
	data := make(map[string]interface{})
	err = decoder.Decode(&data)
	if err != nil {
		log.Print(err.Error())
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

// 直接将某个键的值赋予v
func Decode(key string, v interface{}) bool {
	value, ok := confs[key]
	if !ok {
		return false
	}
	mapOfValue := value.(map[string]interface{})
	s := reflect.ValueOf(v).Elem()

	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		vOfK, ok := mapOfValue[typeOfT.Field(i).Name]
		if ok && f.CanSet() {
			switch f.Kind() {
			case reflect.String:
				f.SetString(vOfK.(string))
			case reflect.Bool:
				f.SetBool(vOfK.(bool))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				switch reflect.ValueOf(vOfK).Kind() {
				case reflect.Float32:
					f.SetInt((int64)(vOfK.(float32)))
				case reflect.Float64:
					f.SetInt((int64)(vOfK.(float64)))
				default:
					f.SetInt(vOfK.(int64))
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				switch reflect.ValueOf(vOfK).Kind() {
				case reflect.Float32:
					f.SetUint((uint64)(vOfK.(float32)))
				case reflect.Float64:
					f.SetUint((uint64)(vOfK.(float64)))
				default:
					f.SetUint(vOfK.(uint64))
				}
			case reflect.Float32, reflect.Float64:
				f.SetFloat(vOfK.(float64))
			}
		}
	}
	return true
}

func getIntValue(v interface{}) interface{} {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Float32:
		return v.(float32)
	case reflect.Float64:
		return v.(float64)
	default:
		return v
	}
}

// 设置配置项
func Set(key string, value interface{}) {
	confs[key] = value
}

// 删除指定项
func Remove(key string) {
	delete(confs, key)
}
