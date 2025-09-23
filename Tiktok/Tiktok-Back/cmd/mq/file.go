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
	//// 初始化
	//initialize.Init()
	////ch := rabbitmqUtils.RabbitMQChannel()
	//ch := initialize.GetChannel()
	//// 简单只开一个 channel 消费
	//deliveries, _ := ch.Consume(
	//	"upload_queue",
	//	"",
	//	false,
	//	false,
	//	false,
	//	false,
	//	nil,
	//)
	//// 创建动态池
	//p := workerpoolUtils.New(core, capacity)
	//defer p.Stop()
	//
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//var wg sync.WaitGroup
	//for i := 0; i < runtime.NumCPU(); i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		for {
	//			select {
	//			case <-ctx.Done():
	//				return
	//			case d, ok := <-deliveries:
	//				if !ok {
	//					return
	//				}
	//				// 拷贝一份，避免并发读写
	//				msg := d
	//				p.Submit(func() {
	//					process(msg)
	//				})
	//			}
	//		}
	//	}()
	//}
	//wg.Wait()
	////
	////for d := range deliveries {
	////	var task types.UploadTask
	////	json.Unmarshal(d.Body, &task)
	////	if err := handle(task); err == nil {
	////		d.Ack(false)
	////	} else {
	////		d.Nack(false, true) // 失败重入队
	////	}
	////}
	//// 工程进入前夕，释放资源
	//defer initialize.Eve()
	//zlog.Infof("程序运行完成！")

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

//
//func process(d amqp.Delivery) {
//	var task types.UploadTask
//	if err := json.Unmarshal(d.Body, &task); err != nil {
//		log.Printf("垃圾信息: %v", err)
//		d.Nack(false, false) // 无法解析直接丢弃
//		return
//	}
//	if err := handle(task); err == nil {
//		d.Ack(false)
//	} else {
//		// 失败重新入队
//		d.Nack(false, true)
//	}
//}

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
	resp := types.UploadResult{ID: task.ID, OK: true, URL: url}

	manager.Push(task.ID, resp)
	return nil
}
