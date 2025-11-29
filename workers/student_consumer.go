package workers

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type StudentCreatedEvent struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	When  int64  `json:"when"`
}

// StartStudentConsumer يبدأ consumer ويستمر بالاستماع (blocking)
func StartStudentConsumer(amqpURL, exchange, queueName, routingKey string) error {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return fmt.Errorf("cannot dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("cannot open channel: %w", err)
	}

	// declare exchange
	if err := ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("exchange declare: %w", err)
	}

	// declare queue (durable)
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	// bind
	if err := ch.QueueBind(queueName, routingKey, exchange, false, nil); err != nil {
		return fmt.Errorf("queue bind: %w", err)
	}

	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	log.Println("student consumer started, waiting for messages...")
	for d := range msgs {
		var evt StudentCreatedEvent
		if err := json.Unmarshal(d.Body, &evt); err != nil {
			log.Println("invalid message body:", err)
			// drop message or send to DLQ (here nack without requeue)
			_ = d.Nack(false, false)
			continue
		}

		// ----- هنا ضع منطق المعالجة الحقيقي -----
		log.Printf("Processing StudentCreatedEvent: ID=%d Name=%s Email=%s\n", evt.ID, evt.Name, evt.Email)
		// مثال: استدعاء repository لتسجيل لوق أو ارسال إيميل
		// err := repository.LogStudentCreated(evt)
		// if err != nil { ... retry ... }

		// acknowledge
		_ = d.Ack(false)
	}
	return nil
}
