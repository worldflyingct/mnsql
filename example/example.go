package main

import (
	"fmt"

	"github.com/worldflyingct/mnsql"
)

func main() {
	mnsql.Set("hello", []byte("hello world!!!"))
	w := mnsql.Get("hello")
	fmt.Println(string(w.([]byte)))

	mnsql.Set("hello", "你好，作者沃航科技")
	w = mnsql.Get("hello")
	fmt.Println(w.(string))

	mnsql.Set("hello", 666)
	w = mnsql.Get("world") // 获取一个不存在的对象
	v, ok := w.(int)
	if ok {
		fmt.Println(v)
	} else {
		fmt.Println("fail")
	}

	mnsql.Set("world", 888)
	w = mnsql.Get("world")
	fmt.Println(w.(int))

	mnsql.Incr("incr")
	w = mnsql.Get("incr")
	fmt.Println(w.(int))
}
