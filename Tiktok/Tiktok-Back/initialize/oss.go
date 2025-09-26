package initialize

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func InitOSS() {
	// 初始化OSS客户端
	client, err := oss.New(global.Config.Oss.Endpoint, global.Config.Oss.AccessKeyID, global.Config.Oss.AccessKeySecret)
	if err != nil {
		zlog.Errorf("oss初始化失败 %v", err)
		panic(err)
	}
	// 获取Bucket
	bucket, err := client.Bucket(global.Config.Oss.BucketName)
	if err != nil {
		zlog.Errorf("oss初始化失败 %v", err)
		panic(err)
	}
	global.OssClient = client
	global.OssBucket = bucket
	return
}
