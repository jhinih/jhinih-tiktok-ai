package logic

import (
	"context"
	"eino_gozero/common/ai_common/ais/aiUtils"
	"eino_gozero/common/ai_common/jhinih_model/chatmodel"
	"fmt"
	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/cozeloop-go"
	"github.com/joho/godotenv"
	"log"
	"os"

	"eino_gozero/apps/AI/rpc/internal/svc"
	"eino_gozero/apps/AI/rpc/types/ai"

	"github.com/zeromicro/go-zero/core/logx"
)

type AIVideoChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAIVideoChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AIVideoChatLogic {
	return &AIVideoChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AIVideoChatLogic) AIVideoChat(in *ai.AIVideoChatRequest) (*ai.AIVideoChatResponse, error) {
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

	g := compose.NewGraph[map[string]string, *schema.Message](
		compose.WithGenLocalState(aiUtils.GenFunc),
	)

	lambda1 := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个较小可爱的邻家妹妹，每次都会用可爱羞涩的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})

	lamada1PreHandler := func(ctx context.Context, input map[string]string, state *aiUtils.State) (output map[string]string, err error) {
		input["content"] = input["content"] + state.History["毫猫"].(string)
		return input, nil
	}
	chatmodel := chatmodel.NewArkChatModel(ctx)

	//加入节点
	err = g.AddLambdaNode("lambda1", lambda1, compose.WithStatePreHandler(lamada1PreHandler))
	if err != nil {
		panic(err)
	}
	err = g.AddChatModelNode("chatmodel", chatmodel)
	if err != nil {
		panic(err)
	}
	//分支连接
	err = g.AddEdge(compose.START, "lambda1")
	if err != nil {
		panic(err)
	}
	//加入分支
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("lambda1", "chatmodel")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("chatmodel", compose.END)
	if err != nil {
		panic(err)
	}

	outsidegraph := compose.NewGraph[map[string]string, string]()

	writeLambda := compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (output string, err error) {
		f, err := os.OpenFile("video_ai.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}
		defer f.Close()
		if _, err := f.WriteString(input.Content + "\n---\n"); err != nil {
			return "", err
		}
		return "video_ai已经写入文件，请前往文件内查看内容", nil
	})

	err = outsidegraph.AddGraphNode("inside", g)
	if err != nil {
		panic(err)
	}
	err = outsidegraph.AddLambdaNode("write", writeLambda)
	if err != nil {
		panic(err)
	}

	err = outsidegraph.AddEdge(compose.START, "inside")
	err = outsidegraph.AddEdge("inside", "write")
	if err != nil {
		panic(err)
	}
	err = outsidegraph.AddEdge("write", compose.END)
	if err != nil {
		panic(err)
	}

	// 编译
	r, err := outsidegraph.Compile(ctx)
	if err != nil {
		panic(err)
	}
	// 执行
	answer, err := r.Invoke(ctx, map[string]string{
		"role":    "姐姐",
		"content": in.Ask,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
	return &ai.AIVideoChatResponse{Answer: answer}, nil
}
