package main

import (
	"fmt"

	"github.com/worldflyingct/mnsql"
)

func main() {
	mnsql.Set("hello", []byte("hello world!!!"))
	w := mnsql.Get("hello")
	str, ok := w.([]byte)
	if ok {
		fmt.Println(string(str))
	} else {
		fmt.Println("fail")
	}
}
