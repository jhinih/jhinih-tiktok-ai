package logic

import (
	"context"
	"eino_gozero/apps/AI/rpc/aiTools/login"
	"eino_gozero/apps/AI/rpc/aiTools/user"
	common "eino_gozero/common/ai_common"
	"eino_gozero/common/ai_common/ais/aiUtils"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"log"

	"eino_gozero/apps/AI/rpc/internal/svc"
	"eino_gozero/apps/AI/rpc/types/ai"

	"github.com/zeromicro/go-zero/core/logx"
)

type AILogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAILogic(ctx context.Context, svcCtx *svc.ServiceContext) *AILogic {
	return &AILogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AILogic) AI(in *ai.AIRequest) (*ai.AIResponse, error) {
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
		"content": in.Ask,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &ai.AIResponse{Answer: answer}, nil
}
