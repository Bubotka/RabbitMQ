package main

import (
	"fmt"
	"github.com/Bubotka/Microservices/proxy/internal/infrastructure/clients/auth/grpc/client_adapter"
	clientadaptergeo "github.com/Bubotka/Microservices/proxy/internal/infrastructure/clients/geo/grpc/client_adapter"
	"github.com/Bubotka/Microservices/proxy/internal/infrastructure/router"
	"github.com/Bubotka/Microservices/proxy/internal/infrastructure/server"
	clientadapteruser "github.com/Bubotka/Microservices/proxy/pkg/clients/user/grpc/client_adapter"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"log"
	"net/http"
	"time"
)

func main() {
	geoConnect, err := clientadaptergeo.Connect("geo:8081")
	if err != nil {
		log.Fatal(err)
	}

	authConnect, err := client_adapter.Connect("auth:8082")
	if err != nil {
		log.Fatal(err)
	}

	userConnect, err := clientadapteruser.Connect("user:8083")
	if err != nil {
		log.Fatal(err)
	}

	geoClientGRpcAdapter := clientadaptergeo.NewGeoClientGRpcAdapter(geoConnect)
	userClientGRpcAdapter := clientadapteruser.NewUserClientGRpcAdapter(userConnect)
	authClientGRpcAdapter := client_adapter.NewAuthClientGRpcAdapter(authConnect)

	conn, err := ConnectAmqpWithRetry("amqp://guest:guest@rabbitmq:5672/")

	rabbitmq.CreateExchange(conn, "limit", "topic")
	rabbitMQ, err := rabbitmq.NewRabbitMQ(conn)
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr: ":8080",
		Handler: router.NewRouter(
			geoClientGRpcAdapter,
			userClientGRpcAdapter,
			authClientGRpcAdapter,
			rabbitMQ,
		),
	}

	serv := server.NewServer(s)

	serv.Serve()
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
