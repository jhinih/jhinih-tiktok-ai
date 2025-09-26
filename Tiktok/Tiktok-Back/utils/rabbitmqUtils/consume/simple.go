package consume

import (
	"Tiktok/utils/rabbitmqUtils"
	"fmt"
	"log"
)

func SimpleConsume(queueName string) {
	// 创建通道
	ch := rabbitmqUtils.RabbitMQChannel()
	defer ch.Close()

	// 声明队列
	q, err := ch.QueueDeclare(
		queueName, // 队列名称
		false,     // 持久性
		false,     // 自动删除
		false,     // 排他性
		false,     // 非阻塞
		nil,       // 其他参数
	)
	if err != nil {
		log.Fatalf("无法声明队列: %v", err)
	}

	// 接收消息
	msgs, err := ch.Consume(
		q.Name, // 队列
		"",     // 消费者
		false,  // 自动确认
		false,  // 排他性
		false,  // 非本地
		false,  // 非阻塞
		nil,    // 其他参数
	)
	if err != nil {
		log.Fatalf("无法注册消费者: %v", err)
	}

	fmt.Println("等待接收消息")
	for d := range msgs {
		fmt.Printf("收到消息: %s\n", d.Body)
		fmt.Println("回复消息")
		d.Ack(false)
	}
}
