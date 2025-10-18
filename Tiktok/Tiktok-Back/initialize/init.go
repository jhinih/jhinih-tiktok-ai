package initialize

import (
	"Tiktok/cmd/flags"
	"Tiktok/global"
	"Tiktok/utils"
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"syscall"
)

var (
	Ch *amqp.Channel
)

func Init() {
	// 解析命令行参数
	flags.Parse()
	// 初始化根目录
	InitPath()
	// 加载配置文件
	InitConfig()
	fmt.Println(global.Config.DB.MySQL.Dsn)
	// 正式初始化日志
	InitLog(global.Config)
	// 初始化数据库
	InitDataBase(*global.Config)
	InitRedis(*global.Config)

	// 初始化消息队列
	//InitMQ(*global.Config)
	// 初始化全局雪花ID生成器
	InitSnowflake()
	// 开启定时任务
	Cron()
	// 初始化OSS服务
	InitOSS()

	//// 初始化ElasticSearch
	//InitElasticsearch()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		global.Stop()
	}()
	// 对命令行参数进行处理
	flags.Run()
}
func InitPath() {
	global.Path = utils.GetRootPath("")
}
