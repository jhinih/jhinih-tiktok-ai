package repository

import (
	"Tiktok/types"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type FileRequest struct {
	Rdb *redis.Client
}

func NewFileRequest(Rdb *redis.Client) *FileRequest {
	return &FileRequest{
		Rdb: Rdb,
	}
}

func (r *FileRequest) SetUploadResult(ctx context.Context, id string, res interface{}, exp time.Duration) error {
	b, _ := json.Marshal(res)
	return r.Rdb.Set(ctx, "upload_res:"+id, b, exp).Err()
}

func (r *FileRequest) GetUploadResult(ctx context.Context, req types.MQUploadResultRequest) (string, error) {
	return r.Rdb.Get(ctx, "upload_res:"+req.ID).Result()
}
