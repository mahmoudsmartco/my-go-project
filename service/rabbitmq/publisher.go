package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

// Message struct - النموذج العام للرسائل
type StudentCreatedEvent struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	When  int64  `json:"when"` // unix timestamp
}

// NewPublisher يتصل بالـ RabbitMQ وينشئ exchange
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
	// أعلن Exchange من النوع direct
	if err := ch.ExchangeDeclare(
		exchange, "direct", true, false, false, false, nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}
	return &Publisher{conn: conn, channel: ch, exchange: exchange}, nil
}

func (p *Publisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *Publisher) PublishStudentCreated(ctx context.Context, evt StudentCreatedEvent, routingKey string) error {
	body, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	// مرّة واحدة مع تأكيد (publisher confirm) — بسيط هنا
	err = p.channel.PublishWithContext(ctx,
		p.exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Unix(evt.When, 0),
			DeliveryMode: amqp.Persistent, // لتخزين الرسالة على القرص إن لزم
		})
	if err != nil {
		return err
	}
	return nil
}
