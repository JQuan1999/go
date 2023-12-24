package main

import (
	"log"

	"github.com/segmentio/kafka-go"
)

func ListTopic() {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		log.Fatal("create connection err: ", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions() // get all partitions
	if err != nil {
		log.Fatal("read partitions err: ", err)
	}
	m := map[string]struct{}{}
	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for k, _ := range m {
		log.Println(k)
	}
}
