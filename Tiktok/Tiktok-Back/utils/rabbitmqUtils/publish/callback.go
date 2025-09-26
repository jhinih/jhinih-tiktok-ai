package publish

import (
	"Tiktok/global"
	"Tiktok/utils/rabbitmqUtils"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func CallBackPublish(queueName string) {
	ch := rabbitmqUtils.RabbitMQChannel()
	defer ch.Close()

	//声明普通队列
	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 声明回调队列
	callbackQueue, err := ch.QueueDeclare(
		"",    // 空名称表示让RabbitMQ自动生成队列名
		false, // 非持久化（因为响应消息不需要长期保存）
		true,  // 自动删除（当消费者断开连接时）
		true,  // 排他性（只允许当前连接访问）
		false, // 非阻塞
		nil,
	)
	msgs, err := ch.Consume(
		callbackQueue.Name, // 回调队列名
		"",                 // 消费者标签
		true,               // 自动确认
		false,              // 非排他性
		false,              // 不阻塞
		false,              // 无参数
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("消息 %d", i)
		msgID := global.SnowflakeNode.Generate().String() // 消息id，用于匹配回调队列中的消息
		err = ch.Publish(
			"",     // 交换机
			q.Name, // 路由键（队列名称）
			false,  // 非强制
			false,  // 非立即
			amqp.Publishing{
				Body:          []byte(msg),
				CorrelationId: msgID,
				ReplyTo:       callbackQueue.Name,
			})
		if err != nil {
			fmt.Println("消息发送失败", err)
			continue
		}
		log.Printf("已发送: %s 消息id %s", msg, msgID)
	}
	for d := range msgs {
		fmt.Printf("收到回调消息 %s %s\n", d.Body, d.CorrelationId)
	}
}
