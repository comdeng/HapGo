package html

import (
	//"code.google.com/p/go-html-transform/h5"
	//"code.google.com/p/go-html-transform/html/transform"
	"bytes"
	_html "code.google.com/p/go.net/html"
	// "io"
	"strings"

//	"fmt"
//	"log"
)

func Tidy(html string) string {
	// tree, _ := h5.NewFromString(html)
	// t := transform.New(tree)
	// return t.String()
	buf := new(bytes.Buffer)
	doc, err := _html.Parse(strings.NewReader(html))
	if err != nil {
		return html
	}
	_html.Render(buf, doc)
	return buf.String()
}
