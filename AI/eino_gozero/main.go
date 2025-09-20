package main

import (
	"context"
	"eino_gozero/common/ai_common/ais/aiUtils"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load() // 加载环境变量
	if err != nil {
		log.Fatal("加载环境变量失败")
	}
	ctx := context.Background()

	outsidegraph := aiUtils.AI("test", "姐姐，我想玩原神，能给我官方工具地址吗？", false)
	// 编译
	r, err := outsidegraph.Compile(ctx)
	if err != nil {
		panic(err)
	}
	// 执行
	answer, err := r.Invoke(ctx, map[string]string{
		"role":    "姐姐",
		"content": "你好呀,可以告诉我原神的url地址嘛？",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}
