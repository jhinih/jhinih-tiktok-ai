package mq

import (
	"Tiktok/configs"
	"Tiktok/log/zlog"
)

type MQ interface {
	InitMQ(config configs.Config) error
	Push(exchange, key string, task interface{}) error
	Consume(queue string, handler func(msg []byte) error) error // ← 新增
	Close() error
}

func InitMQ(mq MQ, config configs.Config) {
	var err error
	err = mq.InitMQ(config)
	if err != nil {
		zlog.Fatalf("无法初始化消息队列 %v", err)
		return
	}
	zlog.Infof("初始化消息队列成功！")
	return
}
