package initialize

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"runtime"
)

func Eve() {
	zlog.Warnf("开始释放资源！")
	errRedis := global.Rdb.Close()
	if errRedis != nil {
		zlog.Errorf("Redis关闭失败 ：%v", errRedis.Error())
	}

	sqlDB, _ := global.DB.DB()
	errDB := sqlDB.Close()
	if errDB != nil {
		zlog.Errorf("数据库关闭失败 ：%v", errDB.Error())
	}

	if global.RabbitMQ != nil {
		if err := global.RabbitMQ.Close(); err != nil {
			zlog.Errorf("RabbitMQ 关闭失败：%v", err)
		}
	}
	runtime.GC()
	if errDB == nil && errRedis == nil {
		zlog.Warnf("资源释放成功！")
	}
}
