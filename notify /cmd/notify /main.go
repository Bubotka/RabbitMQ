package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Не удалось установить соединение с RabbitMQ: %v", err)
	}
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
