package logic

import (
	"context"
	"eino_gozero/apps/AI/rpc/aiTools/login"
	"eino_gozero/common/ai_common/ais/aiUtils"
	"fmt"
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

type AISendCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAISendCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AISendCodeLogic {
	return &AISendCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AISendCodeLogic) AISendCode(in *ai.AISendCodeRequest) (*ai.AISendCodeResponse, error) {
	err := godotenv.Load() // 加载环境变量
	if err != nil {
		log.Fatal("加载环境变量失败")
	}
	ctx := context.Background()

	SendCodeTool := login.CreateSendCodeTool()
	info, err := SendCodeTool.Info(ctx)

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
		SendCodeTool,
	}
	outsidegraph := aiUtils.AIWithTools("send_code", "你需要根据用户提供的邮箱来调用工具发送验证码", tools, infos)
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
	fmt.Println(answer + "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	return &ai.AISendCodeResponse{Answer: answer}, nil
}
