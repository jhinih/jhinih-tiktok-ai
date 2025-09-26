package ais

import (
	"context"
	"eino_gozero/common/ai_common/ais/aiUtils"
	"fmt"
	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/coze-dev/cozeloop-go"
	"github.com/joho/godotenv"
	"log"
)

func CommonAI(content string) string {
	err := godotenv.Load() // 加载环境变量
	if err != nil {
		log.Fatal("加载环境变量失败")
	}
	ctx := context.Background()

	//扣子罗盘调试
	client, err := cozeloop.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close(ctx)
	// 在服务 init 时 once 调用
	handler := ccb.NewLoopHandler(client)
	callbacks.AppendGlobalHandlers(handler)

	outsidegraph := aiUtils.AI("common", "")
	// 编译
	r, err := outsidegraph.Compile(ctx)
	if err != nil {
		panic(err)
	}
	// 执行
	answer, err := r.Invoke(ctx, map[string]string{
		//"role":    "姐姐",
		"content": content,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(answer + "####################################")
	return answer
}
