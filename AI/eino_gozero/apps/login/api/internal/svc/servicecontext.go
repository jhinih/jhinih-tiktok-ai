package svc

import (
	"eino_gozero/apps/login/api/internal/config"
	"eino_gozero/apps/login/rpc/loginclient"
	"eino_gozero/common/pkg/gorms"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	LoginRpc loginclient.Login
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		DB:       gorms.InitGorm(c.Mysql.DataSource),
		Redis:    redis.NewClient(&redis.Options{}),
		LoginRpc: loginclient.NewLogin(zrpc.MustNewClient(c.LoginRpc)),
	}
}
