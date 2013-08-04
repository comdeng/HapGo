package xml

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

var replacer = strings.NewReplacer(
	"&", "&amp;",
	">", "&gt;",
	"<", "&lt;",
	`"`, "&quot;",
	"'", "&apos;",
)

const tabFlag = "\t"

func Map2Xml(data map[string]interface{}, encoding string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="%s"?>\r\n`, encoding))
	datas := map[interface{}]interface{}{
		"HapGo": data,
	}
	convertMap2Xml(datas, buf, uint8(0))
	fh, _ := os.Create("a.log")
	defer fh.Close()
	fh.Write(buf.Bytes())
	return fmt.Sprintf("%s", buf.Bytes())
}

func convertMap2Xml(data map[interface{}]interface{}, buf bytes.Buffer, level uint8) {
	var value string
	for k, v := range data {
		log.Print(reflect.TypeOf(k))
		if _, ok := k.(int); ok {
			k = "item"
		}
		switch v.(type) {
		case bool:
			if v.(bool) {
				value = "true"
			} else {
				value = "false"
			}
		// case int8:
		// 	fallthrough
		// case uint8:
		// 	fallthrough
		// case int16:
		// 	fallthrough
		// case uint16:
		// 	fallthrough
		// case int32:
		// 	fallthrough
		// case uint32:
		// 	fallthrough
		// case int64:
		// 	fallthrough
		// case uint64:
		// 	fallthrough
		// case uint:
		// 	fallthrough
		case int:
			value = fmt.Sprint("%d", v)
		case float32:
		case float64:
			value = fmt.Sprintf("%f", v)
		case string:
			value = replacer.Replace(v.(string))
		case map[interface{}]interface{}:
			buf.WriteString(strings.Repeat(tabFlag, int(level)))
			buf.WriteString("<" + k.(string) + ">\r\n")
			convertMap2Xml(v.(map[interface{}]interface{}), buf, level+1)
			buf.WriteString(strings.Repeat(tabFlag, int(level)))
			buf.WriteString("</" + k.(string) + ">\r\n")
		}
		buf.WriteString(strings.Repeat(tabFlag, int(level)))
		buf.WriteString(fmt.Sprintf(`<%s>%s</%s>\r\n`, k, value, k))
	}

}
