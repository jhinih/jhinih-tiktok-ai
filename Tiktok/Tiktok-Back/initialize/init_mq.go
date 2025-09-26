package initialize

import (
	"Tiktok/configs"
	"Tiktok/log/zlog"
	"Tiktok/pkg/mq"
	"Tiktok/pkg/mq/rabbitmqx"
)

func InitMQ(config configs.Config) {
	for _, name := range config.MQ.Enabled {
		switch name {
		case "rabbitmq":
			mq.InitMQ(rabbitmqx.NewRabbitMQ(), config)
		//case "kafka":
		//	mq.InitMQ(kafkax.NewKafKa(), config)
		default:
			zlog.Fatalf("不支持的消息队列驱动：%s", name)
		}
	}
}
