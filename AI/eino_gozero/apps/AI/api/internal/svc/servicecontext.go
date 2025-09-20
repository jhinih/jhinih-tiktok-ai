package svc

import (
	"eino_gozero/apps/AI/api/internal/config"
	"eino_gozero/apps/AI/rpc/aiclient"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

//type ServiceContext struct {
//	Config config.Config
//	DB     *gorm.DB
//	Redis  *redis.Client
//	AIRpc  aiclient.Ai
//}

//func NewServiceContext(c config.Config) *ServiceContext {
//	return &ServiceContext{
//		Config: c,
//		//DB:      gorms.InitGorm(c.Mysql.DataSource),
//		//Redis:   redis.NewClient(&redis.Options{}),
//		AIRpc: aiclient.NewAi(zrpc.MustNewClient(c.AIRpc)),
//	}
//}

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
	AIRpc  aiclient.Ai
}

func NewServiceContext(c config.Config) *ServiceContext {
	clientConf := zrpc.RpcClientConf{
		Endpoints: []string{"127.0.0.1:8081"},
		Timeout:   300000, // ✅ 5分钟，单位毫秒
	}
	return &ServiceContext{
		Config: c,
		AIRpc:  aiclient.NewAi(zrpc.MustNewClient(clientConf)),
	}
}
