package xml

import (
	"testing"
)

func TestConvert(t *testing.T) {
	data := map[string]interface{}{
		"Err": "hapgo.ok",
		"Data": map[string]interface{}{
			"List": 3,
		},
	}
	str := Map2Xml(data, "utf-8")
	t.Error(str)
}
