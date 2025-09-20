package initialize

import (
	"Tiktok/configs"
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/pkg/database"
	"Tiktok/pkg/mysqlx"
	"Tiktok/pkg/redisx"
)

func InitDataBase(config configs.Config) {
	for _, name := range config.DB.Driver {
		switch name {
		case "mysql":
			database.InitDataBases(mysqlx.NewMySql(), config)
		default:
			zlog.Fatalf("不支持的数据库驱动：%s", name)
		}
	}
	if config.App.Env != "pro" {
		err := global.DB.AutoMigrate()
		if err != nil {
			zlog.Fatalf("数据库迁移失败！")
		}
	}
}
func InitRedis(config configs.Config) {
	if config.Redis.Enable {
		var err error
		global.Rdb, err = redisx.GetRedisClient(config)
		if err != nil {
			zlog.Errorf("无法初始化Redis : %v", err)
		}
		zlog.Infof("初始化Redis成功！")
	} else {
		zlog.Warnf("不使用Redis")
	}
}
