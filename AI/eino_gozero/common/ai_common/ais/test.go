package ais

import (
	"context"
	common "eino_gozero/common/ai_common"
	"eino_gozero/common/ai_common/ais/aiUtils"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load() // 加载环境变量
	if err != nil {
		log.Fatal("加载环境变量失败")
	}
	ctx := context.Background()
	GetGametool := common.CreateTool()
	info, err := GetGametool.Info(ctx)
	if err != nil {
		log.Fatal(err)
	}
	infos := []*schema.ToolInfo{
		info,
	}
	tools := []tool.BaseTool{
		GetGametool,
	}
	outsidegraph := aiUtils.AIWithTools("test", "你是一个较小可爱的邻家妹妹，每次都会用可爱羞涩的语气回答我的问题", tools, infos)
	// 编译
	r, err := outsidegraph.Compile(ctx)
	if err != nil {
		panic(err)
	}
	// 执行
	answer, err := r.Invoke(ctx, map[string]string{
		"role":    "姐姐",
		"content": "你好呀",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}
