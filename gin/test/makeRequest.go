package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ProxyStatus map[string]string
type CommonProxyStatus map[string]string

type QueryResult struct {
	Status   ProxyStatus
	SQLStats []CommonProxyStatus
}

type Result struct {
	Ret   *QueryResult
	Error string
}

type ResultWithTime struct {
	ResultList     []*Result
	MaxProcessTime int64
}

// 自定义http请求
type H map[string]string

// 自定义http request
func makeRequest(ctx context.Context, method string, addr string, headers H, args H, data []byte) (*http.Request, error) {
	// 1. build url
	// 该函数将URL字符串解析为一个url.URL结构体，其中包含了URL的各个组成部分，例如协议、主机、路径、查询参数
	u, _ := url.Parse(strings.Trim(addr, "/"))

	q := u.Query() // 获取query字典
	if args != nil {
		for arg, val := range args {
			q.Add(arg, val) // 添加query参数
		}
	}
	u.RawQuery = q.Encode() // 重新设置url对象的query值

	if headers == nil {
		headers = make(H)
	}

	var body io.Reader
	if data != nil {
		body = bytes.NewBuffer(data)
		headers["content-type"] = "application/json;charset=utf-8"
	} else {
		body = nil
		headers["Content-Type"] = "application/x-www-form-urlencoded;charset=utf-8"
	}
	request, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// 设置request的头部
	if headers != nil {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}
	return request, nil
}

func NewClient() *http.Client {
	return &http.Client{}
}

// 调用GET请求
func DoGet(ctx context.Context, reqURL string, header H, args H) (int, []byte, error) {
	request, err := makeRequest(ctx, "GET", reqURL, header, args, nil)
	if err != nil {
		return 0, nil, err
	}

	client := NewClient()
	resp, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}

	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body) // 从response的body读取数据
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, result, nil
}

// 调用POST请求
func DoPost(ctx context.Context, reqURL string, header H, args H, body []byte) (int, []byte, error) {
	request, err := makeRequest(ctx, "POST", reqURL, header, args, body)
	if err != nil {
		return 0, nil, err
	}
	client := NewClient()
	resp, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}
	// TODO:检查post response的状态码

	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body) // 读取响应体
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, result, nil
}

func TestDoRequest() {
	type ProxyAddress struct {
		Address string `json:"address"`
	}
	proxyList := []ProxyAddress{{Address: "10.177.54.121:3311"}, {Address: "10.177.54.121:3388"}}
	reqURL := "http://127.0.0.1:8887/collect/sqlstats"
	body, err := json.Marshal(proxyList)
	if err != nil {
		panic(err)
	}
	_, responseData, err := DoPost(context.Background(), reqURL, nil, nil, body)
	if err != nil {
		fmt.Println("dopost failed, err: ", err)
		return
	}
	// 反序列化结果
	var result ResultWithTime
	result.ResultList = make([]*Result, 0)
	if err = json.Unmarshal(responseData, &result); err != nil {
		fmt.Println("unmarshal to result failed, err: ", err)
		return
	}
	fmt.Println("max process time= ", result.MaxProcessTime)
	fmt.Println("max process seconds= ", time.Duration(result.MaxProcessTime).Seconds())
	for i := range result.ResultList {
		// 判断result[i]是否为nil, 为nil表示proxy没有采集到数据发生了错误
		if result.ResultList[i].Ret == nil {
			if len(result.ResultList[i].Error) != 0 {
				fmt.Println("collect failed, err: ", result.ResultList[i].Error)
			} else {
				fmt.Println("unexpect error happened")
			}
		} else {
			// 遍历proxy的sql stats
			fmt.Println("result of proxy: ", proxyList[i])
			for _, sqlStats := range result.ResultList[i].Ret.SQLStats {
				fmt.Printf("=======record for sql[%s]=========\n", sqlStats["DigestCode"])
				for k, v := range sqlStats {
					fmt.Println(k, "=", v)
				}
			}

		}
	}
}
