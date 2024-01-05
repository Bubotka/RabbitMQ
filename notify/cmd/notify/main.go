package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"io/ioutil"
	"log"
	"net/http"
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
		sendSMS("Усё ты попал, кончились попытки")
		sendEmail(string(msg.Data))
	}
}

func sendSMS(message string) {
	fmt.Println("Отправляем смс пользовтелю")
	response, err := http.Post("http://proxy:8080/api/sms/send", "application/json", bytes.NewBuffer([]byte(message)))
	fmt.Println("Смс отправили")
	if err != nil {
		fmt.Println(err)
	}

	data, _ := ioutil.ReadAll(response.Body)

	fmt.Println("Содержимое сообщения", string(data))
}

func sendEmail(email string) {
	fmt.Println("Отправляем письмо пользовтелю")
	response, err := http.Post("http://proxy:8080/api/email/send", "application/json", bytes.NewBuffer([]byte(email+"превысил лимит запросов")))
	fmt.Println("Письмо отправили")
	if err != nil {
		fmt.Println(err)
	}

	data, _ := ioutil.ReadAll(response.Body)

	fmt.Println("Содержимое письма", string(data))
}

func ConnectAmqpWithRetry(address string) (*amqp.Connection, error) {
	time.Sleep(15 * time.Second)
	for i := 0; i < 3; i++ {
		conn, err := amqp.Dial(address)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		return conn, nil
	}
	return nil, fmt.Errorf("не удалось подключиться к rabbit")
}
