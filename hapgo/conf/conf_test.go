package conf

import (
	"testing"
)

type TrackIdInfo struct {
	Create bool
	Name   string
	Expire int
	Domain string
	Path   string
}

func TestDecode(t *testing.T) {
	Load("test.conf")
	var tid TrackIdInfo
	Decode("hapgo.tid", &tid)
	if !tid.Create {
		t.Error("Decode() failed", tid.Create, "Expected true")
	}
}

// "hapgo.tid": 			{
// 		"create" 	: true,
// 		"name" 		: "tid",
// 		"expire"	: "360",
// 		"domain"	: "localhost",
// 		"path"		: "/"
// 	}
