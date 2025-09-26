package logic

import (
	"context"
	"eino_gozero/common/ai_common/ais/aiUtils"
	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/coze-dev/cozeloop-go"
	"github.com/joho/godotenv"
	"log"

	"eino_gozero/apps/AI/rpc/internal/svc"
	"eino_gozero/apps/AI/rpc/types/ai"

	"github.com/zeromicro/go-zero/core/logx"
)

type AICommonChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAICommonChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AICommonChatLogic {
	return &AICommonChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AICommonChatLogic) AICommonChat(in *ai.AICommonChatRequest) (*ai.AICommonChatResponse, error) {
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
		"content": in.Ask,
	})
	if err != nil {
		panic(err)
	}
	//go ais.Test1()
	return &ai.AICommonChatResponse{Answer: answer}, nil
}
