package core

import (
	"net/http"
)

const (
	FORMAT_JSON = "json"
	FORMAT_HTML = "html"
	FORMAT_TEXT = "text"
	FORMAT_XML  = "xml"
)

type Responsor struct {
	req       http.ResponseWriter
	outFormat string
}
