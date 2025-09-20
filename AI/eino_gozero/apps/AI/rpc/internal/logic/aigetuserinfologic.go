package logic

import (
	"context"
	"eino_gozero/apps/AI/rpc/internal/svc"
	"eino_gozero/apps/AI/rpc/types/ai"
	"eino_gozero/common/ai_common/ais/aiUtils"
	"eino_gozero/common/ai_common/ais/aiUtils/aiTools/user"
	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/cozeloop-go"
	"github.com/joho/godotenv"
	"log"

	"github.com/zeromicro/go-zero/core/logx"
)

type AIGetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAIGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AIGetUserInfoLogic {
	return &AIGetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AIGetUserInfoLogic) AIGetUserInfo(in *ai.AIRequest) (*ai.AIResponse, error) {
	err := godotenv.Load() // 加载环境变量
	if err != nil {
		log.Fatal("加载环境变量失败")
	}
	ctx := context.Background()

	GetUserInfoTool := user.CreateGetUserInfoTool()
	info, err := GetUserInfoTool.Info(ctx)

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
		GetUserInfoTool,
	}
	outsidegraph := aiUtils.AIWithToolsJson("get_user_info", "你需要根据用户提供的用户ID查询用户信息", tools, infos)
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
		log.Fatal(err.Error())
	}
	return &ai.AIResponse{Answer: answer}, nil
}
