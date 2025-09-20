package svc

import (
	"eino_gozero/apps/AI/rpc/internal/config"
	"eino_gozero/apps/login/rpc/loginclient"
)

type ServiceContext struct {
	Config   config.Config
	LoginRpc loginclient.Login
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		//LoginRpc: loginclient.NewLogin(zrpc.MustNewClient(c.LoginRpc)),
	}
}
