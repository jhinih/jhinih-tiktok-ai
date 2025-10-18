package config

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"log"
	"sync"
)

var inst *Config

func C() *Config { return inst } // 全局可读，无需注入

// Config 全站统一配置，字段用 mapstructure tag
type Config struct {
	service.ServiceConf                     // 继承 go-zero 基础配置（Name/Log/Metrics...）
	App                 ApplicationConfig   `mapstructure:"app"`
	Log                 LoggerConfig        `mapstructure:"log"`
	DB                  DBConfig            `mapstructure:"database"`
	Redis               RedisConfig         `mapstructure:"redis"`
	Email               EmailConfig         `mapstructure:"email"`
	JWT                 JWTConfig           `mapstructure:"jwt"`
	Oss                 OssConfig           `mapstructure:"oss"`
	Elasticsearch       ElasticsearchConfig `mapstructure:"elasticsearch"`
	MQ                  MQConfig            `mapstructure:"mq"`
	Coze                CozeConfig          `mapstructure:"coze"`
}

type ApplicationConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Env         string `mapstructure:"env"`
	LogfilePath string `mapstructure:"logfilePath"`
}
type LoggerConfig struct {
	Level    int8   `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Director string `mapstructure:"director"`
	ShowLine bool   `mapstructure:"show-line"`
}

type DBConfig struct {
	Driver []string `mapstructure:"driver"`
	MySQL  struct {
		AutoMigrate bool   `mapstructure:"migrate"`
		Dsn         string `mapstructure:"dsn"`
	} `mapstructure:"mysql"`
}

type MQConfig struct {
	Enabled  []string `mapstructure:"enabled"`
	RabbitMQ struct {
		Dsn             string `mapstructure:"dsn"`
		ChannelPoolSize string `mapstructure:"channelPoolSize"`
	} `mapstructure:"rabbitmq"`
	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		host    string   `mapstructure:"host"`
		port    int      `mapstructure:"port"`
	} `mapstructure:"kafka"`
}

type RedisConfig struct {
	Enable   bool   `mapstructure:"enable"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type EmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type OssConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"accessKeyID"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	BucketName      string `mapstructure:"bucketName"`
}

type ElasticsearchConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
type CozeConfig struct {
	Token string `mapstructure:"token"`
	//BotID string `mapstructure:"botID"`
}

var (
	once sync.Once
	cfg  Config
)

// MustLoad 单例加载，供所有 main 调用
func MustLoad() *Config {
	once.Do(func() {
		// 1. 先把 .env 灌进去（env 优先级低于代码里显式指定）
		if err := LoadEnv(); err != nil {
			log.Fatalf("load .env error: %v", err)
		}

		// 2. 用 go-zero 的 conf 包填充：先读 yaml（可选），再用环境变量覆盖
		//    如果没有 yaml 文件，直接全 env 也能跑
		conf.MustLoad("config.yaml", &cfg, conf.UseEnv())

		// 3. 兜底默认值
		if cfg.Mode == "" {
			cfg.Mode = "pro"
		}
	})
	return &cfg
}
func init() {
	// 仅第一次 import 时执行
	MustLoad()
}
