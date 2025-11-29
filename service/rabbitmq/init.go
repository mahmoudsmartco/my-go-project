package rabbitmq

import (
	"log"
)

var DefaultPublisher *Publisher

// InitDefaultPublisher ينشئ DefaultPublisher وإمكانية الوصول له عالمياً
func InitDefaultPublisher(amqpURL, exchange string) {
	p, err := NewPublisher(amqpURL, exchange)
	if err != nil {
		log.Printf("rabbitmq init failed: %v", err)
		// لا نفشل التطبيق هنا، لكن DefaultPublisher سيبقى nil — تحقق قبل الاستخدام
		return
	}
	DefaultPublisher = p
	log.Printf("rabbitmq publisher initialized (exchange=%s)", exchange)
}

// CloseDefaultPublisher يغلق الاتصال عند غلق التطبيق
func CloseDefaultPublisher() {
	if DefaultPublisher != nil {
		DefaultPublisher.Close()
	}
}
