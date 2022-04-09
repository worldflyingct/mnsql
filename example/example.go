package main

import (
	"fmt"

	"github.com/worldflyingct/mnsql"
)

func main() {
	mnsql.Set("hello", []byte("hello world!!!"))
	w, _ := mnsql.Get("hello")
	fmt.Println(string(w.([]byte)))

	mnsql.Set("hello", "你好，作者沃航科技")
	w, _ = mnsql.Get("hello")
	fmt.Println(w.(string))

	mnsql.Set("hello", 666)
	_, r := mnsql.Get("world") // 获取一个不存在的对象
	if r == -2 {
		fmt.Println("key不存在")
	} else {
		fmt.Println("key存在")
	}

	mnsql.Set("world", 888)
	w, _ = mnsql.Get("world")
	fmt.Println(w.(int))

	mnsql.Incr("incr")
	w, _ = mnsql.Get("incr")
	fmt.Println(w.(int))
}
