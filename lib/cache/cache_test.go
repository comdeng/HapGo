package cache

import (
	"testing"
)

func TestSet(t *testing.T) {
	Set("foo", "bar", 1)
	value, ok := Get("foo")
	if ok {
		ret := value.(string)
		if ret != "bar" {
			t.Error("get data wrong:%s", ret)
		}
	} else {
		t.Error("get nil data")
	}
}
