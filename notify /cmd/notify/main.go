package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"log"
	"time"
)

func main() {
	conn, err := ConnectAmqpWithRetry("amqp://guest:guest@rabbitmq:5672/")
	fmt.Println("Удалось подключиться к rabbit")
	rabbitMQ, err := rabbitmq.NewRabbitMQ(conn)
	if err != nil {
		log.Fatal(err)
	}

	messages, err := rabbitMQ.Subscribe("limit")
	for msg := range messages {
		fmt.Println(msg)
	}

	fmt.Println("Вышли из приложения")
}

func ConnectAmqpWithRetry(address string) (*amqp.Connection, error) {
	for i := 0; i < 5; i++ {
		conn, err := amqp.Dial(address)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		return conn, nil
	}
	return nil, fmt.Errorf("не удалось подключиться к rabbit")
}
