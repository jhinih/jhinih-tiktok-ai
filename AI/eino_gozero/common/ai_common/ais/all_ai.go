package ais

import (
	"context"
	common "eino_gozero/common/ai_common"
	"eino_gozero/common/ai_common/ais/aiUtils"
	"eino_gozero/common/ai_common/ais/aiUtils/aiTools/login"
	"eino_gozero/common/ai_common/ais/aiUtils/aiTools/user"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"log"
)

func AllAI(content string) string {
	ctx := context.Background()

	SendCodeTool := login.CreateSendCodeTool()
	info, err := SendCodeTool.Info(ctx)
	GetGameNameTool := common.CreateTool()
	info1, err := GetGameNameTool.Info(ctx)
	GetUserInfoTool := user.CreateGetUserInfoTool()
	info2, err := GetGameNameTool.Info(ctx)

	if err != nil {
		log.Fatal(err)
	}

	infos := []*schema.ToolInfo{
		info,
		info1,
		info2,
	}
	tools := []tool.BaseTool{
		SendCodeTool,
		GetGameNameTool,
		GetUserInfoTool,
	}
	outsidegraph := aiUtils.AIWithTools("all",
		"你必须使用工具完成用户请求，工具如下：\n"+
			"- send_code: 发送验证码，需要 email 参数\n"+
			"- GetUserInfo: 获取用户信息，需要 id 参数\n"+
			"- get_game: 获取游戏链接，需要 name 参数\n\n"+
			"如果用户提到“用户的id”或“用户ID”，"+
			"你必须调用 GetUserInfo 工具，并传入 id。\n",
		tools,
		infos,
	)
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
		log.Fatal(err)
	}
	return answer
}
