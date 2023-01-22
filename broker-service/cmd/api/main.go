package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbitmq, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	app := Config{Rabbit: rabbitmq}

	log.Printf("Starting broker service on port %s\n", webPort)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("RabbitMQ not yet ready...")
			counts += 1

		} else {
			log.Println("Connected to rabbitmq")
			connection = c
			break
		}

		if counts > 5 {
			return nil, err
		}

		time.Sleep(backoff)
	}
	return connection, nil
}
