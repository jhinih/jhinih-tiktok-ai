package main

import (
	"Tiktok/global"
	"Tiktok/initialize"
	"Tiktok/log/zlog"
	"Tiktok/manager"
	"Tiktok/repository"
	"Tiktok/types"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime"
	"os"
	"time"
)

func main() {
	initialize.Init()

	err := global.RabbitMQ.Consume("upload_queue", func(body []byte) error {
		var task types.UploadTask
		if err := json.Unmarshal(body, &task); err != nil {
			return err
		}
		return handle(task)
	})
	if err != nil {
		zlog.Fatalf("启动消费者失败: %v", err)
	}

	// 阻塞直到程序退出
	<-global.CtxDone()
	initialize.Eve()
}

func handle(task types.UploadTask) error {
	ctx := context.Background()
	fileBytes, err := os.ReadFile(task.TmpPath)
	if err != nil {
		return err
	}
	defer os.Remove(task.TmpPath) // 清理临时文件

	hasher := sha256.New()
	hasher.Write(fileBytes)
	hash := hex.EncodeToString(hasher.Sum(nil))
	newFilename := hash + task.Ext

	exist, _ := global.OssBucket.IsObjectExist(newFilename)
	if !exist {
		contentType := mime.TypeByExtension(task.Ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		err = global.OssBucket.PutObject(newFilename,
			bytes.NewReader(fileBytes),
			oss.ACL(oss.ACLPublicRead),
			oss.ContentType(contentType))
		if err != nil {
			repository.NewFileRequest(global.Rdb).SetUploadResult(ctx, task.ID, map[string]interface{}{
				"id":    task.ID,
				"ok":    false,
				"error": err.Error()},
				5*time.Minute,
			)
			return err
		}
	}
	url := fmt.Sprintf("https://%s.%s/%s", global.Config.Oss.BucketName, global.Config.Oss.Endpoint, newFilename)
	repository.NewFileRequest(global.Rdb).SetUploadResult(ctx, task.ID,
		map[string]interface{}{
			"id":  task.ID,
			"ok":  true,
			"url": url,
		},
		5*time.Minute,
	)
	//SSE
	resp := types.UploadResult{ID: task.ID, OK: true, URL: url}

	manager.Push(task.ID, resp)
	return nil
}
