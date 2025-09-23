package ais

import (
	"context"
	"eino_gozero/common/ai_common/ais/aiUtils"
	"eino_gozero/common/ai_common/ais/aiUtils/aiTools/videos"
	"fmt"
	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/cozeloop-go"
	"github.com/joho/godotenv"
	"log"
)

func GetVideoAI(content string) string {
	err := godotenv.Load() // 加载环境变量
	if err != nil {
		log.Fatal("加载环境变量失败")
	}
	ctx := context.Background()

	GetVideoTool := videos.CreateGetVideoTool()
	info, err := GetVideoTool.Info(ctx)

	//扣子罗盘调试
	client, err := cozeloop.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close(ctx)
	// 在服务 init 时 once 调用
	handler := ccb.NewLoopHandler(client)
	callbacks.AppendGlobalHandlers(handler)

	if err != nil {
		log.Fatal(err)
	}

	infos := []*schema.ToolInfo{
		info,
	}
	tools := []tool.BaseTool{
		GetVideoTool,
	}
	outsidegraph := aiUtils.AIWithToolsJson("get_Video_info", "当用户需要视频时，你需要获取视频并交给用户", tools, infos)
	// 编译
	r, err := outsidegraph.Compile(ctx)
	if err != nil {
		panic(err)
	}
	// 执行
	answer, err := r.Invoke(ctx, map[string]string{
		"content": content,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(answer + "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@" + "get_video")
	return answer
}
