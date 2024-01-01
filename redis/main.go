package main

import "redis/history"

type Task struct {
	Address string `json:"address"`
}

func main() {
	history.TestPipelineGet()
}
