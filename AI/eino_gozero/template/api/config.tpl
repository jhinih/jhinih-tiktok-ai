package config

import {{.authImport}}

type Config struct {
	rest.RestConf
    Mysql struct {
        DataSource string
    }
    Redis struct {
        Host     string
        Port     int
        Password string
        DB       int
    }
    //AIRpc zrpc.RpcClientConf
	{{.auth}}
	{{.jwtTrans}}
}
