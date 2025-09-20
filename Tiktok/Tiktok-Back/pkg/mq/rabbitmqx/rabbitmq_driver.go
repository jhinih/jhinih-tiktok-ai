package rabbitmqx

import (
	"Tiktok/configs"
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/pkg/mq"
	"Tiktok/utils/workerpoolUtils"
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"strconv"
	"sync"
	"sync/atomic"
)

//type RabbitMQ struct {
//}
//
//// InitMQ 初始化
//func (m *RabbitMQ) InitMQ(config configs.Config) (*amqp.Connection, error) {
//	dsn := m.GetMQDsn(config)
//	// 连接 MQ
//	conn, err := amqp.Dial(dsn)
//	if err != nil {
//		zlog.Panicf("RabbitMQ无法连接！: %v", err)
//	}
//	zlog.Infof("RabbitMQ连接成功！")
//	return conn, nil
//}
//func (m *RabbitMQ) GetMQDsn(config configs.Config) string {
//	return config.MQ.RabbitMQ.Dsn
//}
//func NewRabbitMQ() mq.MQ {
//	return &RabbitMQ{}
//}

type RabbitMQ struct {
	conn  *amqp.Connection
	pool  []*amqp.Channel
	index uint32
	mu    sync.Mutex
}

func NewRabbitMQ() mq.MQ {
	return &RabbitMQ{}
}

const (
	core     = 4  // 最小并发
	capacity = 32 // 最大并发
)

func (r *RabbitMQ) InitMQ(config configs.Config) error {
	conn, err := amqp.Dial(config.MQ.RabbitMQ.Dsn)
	if err != nil {
		return err
	}
	r.conn = conn
	PoolSize, _ := strconv.ParseInt(config.MQ.RabbitMQ.ChannelPoolSize, 10, 64)
	poolSize := int(PoolSize)
	if poolSize <= 0 {
		poolSize = 8
	}

	r.pool = make([]*amqp.Channel, poolSize)
	for i := 0; i < poolSize; i++ {
		ch, err := conn.Channel()
		if err != nil {
			return err
		}
		r.pool[i] = ch
	}
	global.RabbitMQ = r

	zlog.Infof("RabbitMQ 初始化成功，channel 池大小：%d", poolSize)
	return r.declareInfrastructure()
}

// 基础设施 =  exchange + queue + binding 映射表
func (r *RabbitMQ) declareInfrastructure() error {
	infra := []struct {
		Exchange string
		Type     string
		Queue    string
		Key      string
	}{
		{"upload", "direct", "upload_queue", "upload_key"},
		{"email", "direct", "email_queue", "email_key"},
		{"order", "topic", "order_event", "order.created"},
		// 有新业务继续加一行，不改业务代码
	}

	ch := r.getChannel()
	for _, v := range infra {
		if err := ch.ExchangeDeclare(v.Exchange, v.Type, true, false, false, false, nil); err != nil {
			return err
		}
		if _, err := ch.QueueDeclare(v.Queue, true, false, false, false, nil); err != nil {
			return err
		}
		if err := ch.QueueBind(v.Queue, v.Key, v.Exchange, false, nil); err != nil {
			return err
		}
	}
	return nil
}
func (r *RabbitMQ) Push(exchange, key string, task interface{}) error {
	ch := r.getChannel()
	body, _ := json.Marshal(task)
	return ch.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (r *RabbitMQ) Consume(queue string, handler func(msg []byte) error) error {
	ch := r.getChannel()
	if err := ch.Qos(capacity, 0, false); err != nil {
		return err
	}
	deliveries, err := ch.Consume(
		queue,
		"",    // consumer tag
		false, // manual ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,
	)
	if err != nil {
		return err
	}

	wp := workerpoolUtils.New(core, capacity) // 你用过的池
	defer wp.Stop()

	// 优雅退出
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-global.CtxDone() // 你程序退出时调用
		cancel()
		ch.Close()
	}()

	for d := range deliveries {
		msg := d
		wp.Submit(func() {
			if err := handler(msg.Body); err != nil {
				zlog.Errorf("消费失败: %v", err)
				_ = msg.Nack(false, true) // 重入队
				return
			}
			_ = msg.Ack(false)
		})
	}
	<-global.CtxDone() //  阻塞在这里
	return nil
}

// pkg/mq/rabbitmqx/rabbitmqx.go
func (r *RabbitMQ) Close() error {
	// 1. 关闭所有 channel
	for _, ch := range r.pool {
		if ch != nil {
			_ = ch.Close()
		}
	}
	// 2. 关闭连接
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
func (r *RabbitMQ) getChannel() *amqp.Channel {
	idx := atomic.AddUint32(&r.index, 1)
	return r.pool[idx%uint32(len(r.pool))]
}
