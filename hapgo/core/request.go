package core

import (
	"net/http"
)

type Requestor struct {
	http.Request
}
