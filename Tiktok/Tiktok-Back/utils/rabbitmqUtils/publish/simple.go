package publish

import (
	"Tiktok/utils/rabbitmqUtils"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"log"
)

func SimplePublish(queueName string) {
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
	// 发送消息
	body := fmt.Sprintf("msg-%d", 1)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			Body: []byte(body),
		})
	if err != nil {
		logrus.Errorf("发送失败: %v", err)
		return
	}
	fmt.Printf("发送消息: %s\n", body)
}

//type OrderMsg struct {
//    OrderID uint64 `json:"order_id"`
//    UserID  uint64 `json:"user_id"`
//    Price   uint   `json:"price"`
//}
//
//func publishOrder(ch *amqp.Channel, o OrderMsg) error {
//    body, err := json.Marshal(o)
//    if err != nil {
//        return err
//    }
//    return ch.Publish(
//        "order.exchange", // exchange
//        "order.new",      // routing key
//        false, false,
//        amqp.Publishing{
//            ContentType: "application/json", // 方便消费者识别
//            Body:        body,
//            Timestamp:   time.Now(),
//        })
//}

//var o OrderMsg
//if err := json.Marshal(d.Body, &o); err != nil { ... }
