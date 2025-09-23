package logic

import (
	"context"
	"eino_gozero/common/ai_common/ais/aiUtils"
	"eino_gozero/common/ai_common/ais/aiUtils/aiTools/login"
	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/cozeloop-go"
	"github.com/joho/godotenv"
	"log"

	"eino_gozero/apps/AI/rpc/internal/svc"
	"eino_gozero/apps/AI/rpc/types/ai"

	"github.com/zeromicro/go-zero/core/logx"
)

type AISendResumeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAISendResumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AISendResumeLogic {
	return &AISendResumeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AISendResumeLogic) AISendResume(in *ai.AIRequest) (*ai.AIResponse, error) {
	err := godotenv.Load() // 加载环境变量
	if err != nil {
		log.Fatal("加载环境变量失败")
	}
	ctx := context.Background()

	SendResumeTool := login.CreateSendResumeTool()
	info, err := SendResumeTool.Info(ctx)

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
		SendResumeTool,
	}
	outsidegraph := aiUtils.AIWithTools("send_resume", "你需要根据用户提供的邮箱来调用工具发送简历", tools, infos)
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
		panic(err)
	}
	return &ai.AIResponse{Answer: answer}, nil
}
