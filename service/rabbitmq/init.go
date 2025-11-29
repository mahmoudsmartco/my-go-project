package rabbitmq

import "log"

var DefaultPublisher *Publisher

func InitDefaultPublisher() {
	p, err := NewPublisher("amqp://guest:guest@localhost:5672/", "students.exchange")
	if err != nil {
		log.Fatalf("rabbitmq publisher init failed: %v", err)
	}
	DefaultPublisher = p
}
