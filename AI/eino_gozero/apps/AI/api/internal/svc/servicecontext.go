package svc

import (
	"eino_gozero/apps/AI/api/internal/config"
	"eino_gozero/apps/AI/rpc/aiclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	//DB     *gorm.DB
	//Redis  *redis.Client
	AIRpc aiclient.Ai
}

func NewServiceContext(c config.Config) *ServiceContext {
	clientConf := zrpc.RpcClientConf{
		Endpoints: []string{"127.0.0.1:8081"},
		Timeout:   300000, // ，单位毫秒
	}
	return &ServiceContext{
		Config: c,
		AIRpc:  aiclient.NewAi(zrpc.MustNewClient(clientConf)),
	}
}
