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

func StartStudentConsumer(amqpURL, exchange, queueName, routingKey string) error {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	// Declare exchange & queue & bind
	if err := ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
		return err
	}

	// Dead-letter queue setup could be added via args
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	if err := ch.QueueBind(queueName, routingKey, exchange, false, nil); err != nil {
		return err
	}

	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var evt StudentCreatedEvent
			if err := json.Unmarshal(d.Body, &evt); err != nil {
				log.Println("invalid message:", err)
				d.Nack(false, false) // drop or send to DLQ
				continue
			}

			// معالجة الرسالة: مثال تسجيل لوج أو إرسال إيميل (هنا سنسجل)
			log.Printf("Processing StudentCreatedEvent: %+v\n", evt)

			// مثال: نكتب سجل في DB (repository مثال)
			// err = repository.LogStudentCreated(evt.ID, evt.Name, evt.When)
			// if err != nil { ... retry ... }

			// successful
			d.Ack(false)
		}
	}()

	fmt.Println("Student consumer started, waiting messages...")
	// لا نغلق conn/ch لأننا نريد worker يعمل طوال الوقت
	return nil
}
