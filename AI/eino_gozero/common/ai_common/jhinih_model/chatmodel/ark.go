package chatmodel

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func NewArkChatModel(ctx context.Context) *ark.ChatModel {
	err := godotenv.Load() // 加载环境变量
	if err != nil {
		log.Fatal("加载环境变量失败")
	}
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})
	if err != nil {
		log.Fatal(err)
	}
	return model
}
