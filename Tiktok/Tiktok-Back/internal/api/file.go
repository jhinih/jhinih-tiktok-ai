package api

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/manager"
	"Tiktok/repository"
	"Tiktok/response"
	"Tiktok/types"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

// UploadFile 上传文件,因为比较难以分类，所以只在API中实现
func UploadFile(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	resp := types.UploadFileResponse{}
	// 限制上传文件大小（示例为1000MB）
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1000*1024*1024)

	// 获得上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		zlog.CtxErrorf(ctx, "获取上传文件失败: %v", err)
		err = response.ErrResponse(err, response.INTERNAL_ERROR)
		response.Response(c, resp, err)
		return
	}
	defer file.Close()

	// 读取文件内容到内存
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		zlog.CtxErrorf(ctx, "读取文件内容失败: %v", err)
		err = response.ErrResponse(err, response.INTERNAL_ERROR)
		response.Response(c, resp, err)
		return
	}

	// 计算SHA-256哈希值
	hasher := sha256.New()
	hasher.Write(fileBytes)
	hash := hex.EncodeToString(hasher.Sum(nil))

	// 获取文件扩展名并构造新文件名
	ext := filepath.Ext(header.Filename)
	newFilename := hash + ext

	// 检查OSS中是否已存在该文件
	exist, err := global.OssBucket.IsObjectExist(newFilename)
	if err != nil {
		zlog.CtxErrorf(ctx, "检查文件存在失败: %v", err)
		err = response.ErrResponse(err, response.INTERNAL_ERROR)
		response.Response(c, resp, err)
		return
	}

	if !exist {
		// 设置正确的Content-Type
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		// 上传到OSS
		reader := bytes.NewReader(fileBytes)
		err = global.OssBucket.PutObject(newFilename, reader,
			oss.ACL(oss.ACLPublicRead),
			oss.ContentType(contentType),
		)
		if err != nil {
			zlog.CtxErrorf(ctx, "上传文件到OSS失败: %v", err)
			err = response.ErrResponse(err, response.INTERNAL_ERROR)
			response.Response(c, resp, err)
			return
		}
	}

	// 构造访问URL
	url := fmt.Sprintf("https://%s.%s/%s", global.Config.Oss.BucketName, global.Config.Oss.Endpoint, newFilename)
	resp.Url = url

	zlog.CtxInfof(ctx, "上传成功，访问URL: %s", url)
	response.Response(c, resp, nil)
}

func init() { os.MkdirAll(global.TMPDIR, 0755) }

// UploadFile 入口：只落盘+发 MQ
func MQUploadFile(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	// 限制上传文件大小（示例为1000MB）
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1000*1024*1024)

	// 获得上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		zlog.CtxErrorf(ctx, "获取上传文件失败: %v", err)
		err = response.ErrResponse(err, response.INTERNAL_ERROR)
		response.Response(c, nil, err)
		return
	}
	defer file.Close()

	//resp := gin.H{}

	ID := global.SnowflakeNode.Generate().String()
	ext := filepath.Ext(header.Filename)
	tmpPath := filepath.Join(global.TMPDIR, ID+ext)

	out, err := os.Create(tmpPath)
	if err != nil {
		zlog.CtxErrorf(ctx, "创建临时文件失败: %v", err)
		response.Response(c, nil, response.ErrResponse(err, response.INTERNAL_ERROR))
		return
	}
	defer out.Close()

	// 将上传的文件内容写入临时文件
	_, err = io.Copy(out, file)
	if err != nil {
		zlog.CtxErrorf(ctx, "保存上传文件失败: %v", err)
		response.Response(c, nil, response.ErrResponse(err, response.INTERNAL_ERROR))
		return
	}
	task := types.UploadTask{
		ID:       ID,
		TmpPath:  tmpPath,
		FileName: header.Filename,
		Ext:      ext,
	}
	if err := global.RabbitMQ.Push("upload", "upload_key", task); err != nil {
		response.Response(c, nil, response.ErrResponse(err, response.RABBITMQ_ERROR))
		return
	}
	c.JSON(http.StatusOK, gin.H{"ID": ID})
}

// UploadResult 轮询接口
func MQUploadResult(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.MQUploadResultRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "轮询视频上传判断请求: %v", req)
	resp, err := repository.NewFileRequest(global.Rdb).GetUploadResult(ctx, req)
	if err != nil {
		response.Response(c, nil, err)
		return
	}
	response.Response(c, resp, nil)
}

// UploadSSE 长连接等待结果
func UploadSSE(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.UploadSSERequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "SSE视频上传判断请求: %v", req)

	// 设置 SSE 头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	//注册
	conn, cancel := manager.Register(req.ID)
	defer cancel()
	// 阻塞等消息
	msg := <-conn.Chan()

	// SSE 格式
	c.Writer.Write([]byte("data: "))
	c.Writer.Write(msg)
	c.Writer.Write([]byte("\n\n"))
	c.Writer.(http.Flusher).Flush()
}
