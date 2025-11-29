package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type StudentCreatedEvent struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func connectWithRetry() (*amqp.Connection, error) {
	var attempt int

	for {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err == nil {
			fmt.Println("🟢 RabbitMQ Connected Successfully")
			return conn, nil
		}

		attempt++
		wait := time.Duration(1<<attempt) * time.Second // exponential backoff

		if wait > 30*time.Second {
			wait = 30 * time.Second // max wait time
		}

		log.Printf("🔴 RabbitMQ connection failed. Retrying in %v ...", wait)
		time.Sleep(wait)
	}
}

func StartConsumer() {
	for {
		conn, _ := connectWithRetry()
		ch, err := conn.Channel()
		if err != nil {
			log.Println("❌ Failed to open channel:", err)
			continue
		}

		q, err := ch.QueueDeclare(
			"student_created",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Println("❌ Queue declare failed:", err)
			ch.Close()
			continue
		}

		msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Println("❌ Failed to consume messages:", err)
			ch.Close()
			continue
		}

		fmt.Println("📡 Consumer Waiting for messages...")

		for msg := range msgs {
			var event StudentCreatedEvent
			json.Unmarshal(msg.Body, &event)

			log.Printf("🟢 Processing StudentCreatedEvent: ID=%d Name=%s Email=%s",
				event.ID, event.Name, event.Email)
		}

		log.Println("🔴 Connection lost. Reconnecting...")
	}
}
