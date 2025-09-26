package svc

import (
	{{.configImport}}
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
)

type ServiceContext struct {
	Config {{.config}}
    DB      *gorm.DB
    Redis   *redis.Client
    //AIRpc aiclient.Ai
	{{.middleware}}
}

func NewServiceContext(c {{.config}}) *ServiceContext {
	return &ServiceContext{
		Config: c,
        //DB:      gorms.InitGorm(c.Mysql.DataSource),
        //Redis:   redis.NewClient(&redis.Options{}),
        //AIRpc: aiclient.NewAi(zrpc.MustNewClient(c.AIRpc)),
		{{.middlewareAssignment}}
	}
}
