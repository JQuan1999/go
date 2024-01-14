package main

import (
	"fmt"
	"net/url"
)

func TestUrlParse() {
	urlString := "https://example.com/search?q=golang&page=1"
	u, err := url.Parse(urlString) // 用url.Parse函数将URL字符串解析为一个url.URL结构体
	if err != nil {
		panic(err)
	}
	query := u.Query()               // 使用u.Query方法将查询参数解析为一个url.Values类型的值
	fmt.Println(query.Get("q"))      // 查询q的参数值"golang"
	fmt.Println(query.Get("page"))   // 查询page的参数值"1"
	fmt.Println(query.Get("offset")) // ""
}
