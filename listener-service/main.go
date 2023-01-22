package main

import (
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yushyn-andriy/listener-service/event"
)

func main() {
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Println("Listening for and consuming RabbitMQ messages...")
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.WARNING"})
	if err != nil {
		log.Println("listening:", err)
		os.Exit(1)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
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
