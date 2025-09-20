package rabbitmqUtils

import (
	"Tiktok/global"
	"github.com/streadway/amqp"
	"log"
)

func RabbitMQChannel() *amqp.Channel {
	ch, err := global.MQ.Channel()
	if err != nil {
		log.Fatalf("无法打开通道: %v", err)
	}
	return ch
}
