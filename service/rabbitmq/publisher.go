package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher يفتح اتصال ويهيئ exchange
type Publisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	closed   bool
}

type StudentCreatedEvent struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	When  int64  `json:"when"` // unix timestamp
}

func NewPublisher(amqpURL, exchange string) (*Publisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare exchange
	if err := ch.ExchangeDeclare(
		exchange,
		"direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	p := &Publisher{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
		closed:   false,
	}
	return p, nil
}

func (p *Publisher) Close() {
	if p == nil {
		return
	}
	if p.closed {
		return
	}
	p.closed = true
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}

// PublishStudentCreated ينشر حدث StudentCreated إلى exchange مع routing key
func (p *Publisher) PublishStudentCreated(ctx context.Context, evt StudentCreatedEvent, routingKey string) error {
	if p == nil || p.channel == nil {
		return errors.New("publisher not initialized")
	}

	body, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	// استخدام PublishWithContext مع Timeout
	ctxPub, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(ctxPub,
		p.exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Unix(evt.When, 0),
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		// لا نرمي الخطأ الحاد للـ HTTP client في حال فشل النشر (نكتفي بالـ log)
		log.Printf("rabbitmq publish error: %v", err)
		return err
	}
	return nil
}
