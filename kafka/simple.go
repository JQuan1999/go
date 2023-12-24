package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

func producer(wg *sync.WaitGroup) {
	defer wg.Done()
	topic := "test"
	partition := 0
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader: ", err)
	}
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages: ", err)
	}
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer: ", err)
	}
}

func consumer(wg *sync.WaitGroup) {
	defer wg.Done()
	topic := "test"
	partition := 0
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader: ", err)
	}
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6)
	b := make([]byte, 10e3) // 10kB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		log.Println(string(b[:n]))
	}
	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch: ", err)
	}
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection: ", err)
	}
}
