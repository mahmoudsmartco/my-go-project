package main

import (
	"log"
	"app2_http_api_database/workers"
)

func main() {
	amqpURL := "amqp://guest:guest@localhost:5672/"
	exchange := "students.exchange"
	queueName := "students.created.queue"
	routingKey := "students.created"

	log.Println("starting student consumer...")
	if err := workers.StartStudentConsumer(amqpURL, exchange, queueName, routingKey); err != nil {
		log.Fatalf("consumer failed: %v", err)
	}
}
