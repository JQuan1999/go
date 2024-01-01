package main

import "fmt"

func main() {
	// 创建gateway
	// TODO:gateway定义指标、暴露接口、采集数据、返回数据
	// gateway工作流程：redis获取proxy节点数据、创建worker关联proxy、下发任务、聚合结果写回
	m := map[string]string{"name": "abc", "age": "19"}
	for _, value := range m {
		value = value + value
	}
	for key, value := range m {
		fmt.Println(key, value)
	}
}
