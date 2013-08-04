package html

import (
	"testing"
)

func TestTidy(t *testing.T) {
	dest := "<div id='hello'><a onclick=\"window.location='aa'\" id=\"ss\"></a>"
	str := Tidy(dest)
	t.Error(str)
}
