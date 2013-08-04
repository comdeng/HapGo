package main

import (
	"encoding/xml"
	"fmt"
)

type xmlInfo struct {
	Name string
	Sex  string
	Age  int
}

type XmlInfo map[string]interface{}

func main() {
	xi := &XmlInfo{
		"Name": "ronnie",
		"Sex":  "male",
		"Age":  32,
	}
	b, err := xml.Marshal(xi)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("%s", b)
}
