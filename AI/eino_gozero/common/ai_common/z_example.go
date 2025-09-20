package common

import (
	"context"
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
)

//	type State struct {
//		History map[string]any
//	}
//
//	func genFunc(ctx context.Context) *State {
//		return &State{
//			History: make(map[string]any),
//		}
//	}
func TextAI(content string) string {
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
		compose.WithGenLocalState(genFunc),
	)

	lambda0 := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output map[string]string, err error) {
		_ = compose.ProcessState[*State](ctx, func(ctx context.Context, state *State) error {
			state.History["毫猫"] = "你会喜欢我吗？"
			state.History["耄耋"] = "你好呀，初次见面，"
			//"我喜欢你"
			return nil
		})
		if input["role"] == "妹妹" {
			return map[string]string{"role": "毫猫", "content": input["content"]}, nil
		} else if input["role"] == "姐姐" {
			return map[string]string{"role": "耄耋", "content": input["content"]}, nil
		}
		return map[string]string{"role": "user", "content": input["content"]}, nil
	})

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

	lambda2 := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个高冷傲娇的大小姐，每次都会用傲娇无理的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	lamada1PreHandler := func(ctx context.Context, input map[string]string, state *State) (output map[string]string, err error) {
		input["content"] = input["content"] + state.History["毫猫"].(string)
		return input, nil
	}
	lamada2PreHandler := func(ctx context.Context, input map[string]string, state *State) (output map[string]string, err error) {
		input["content"] = input["content"] + state.History["耄耋"].(string)
		return input, nil
	}

	chatmodel := chatmodel.NewArkChatModel(ctx)

	//GetGametool := CreateTool()
	//info, err := GetGametool.Info(ctx)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//infos := []*schema.ToolInfo{
	//	info,
	//}
	//chatmodel.BindTools(infos)
	//ToolNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
	//	Tools: []tool.BaseTool{
	//		GetGametool,
	//	},
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}

	//加入节点
	err = g.AddLambdaNode("lambda0", lambda0)
	if err != nil {
		panic(err)
	}
	err = g.AddLambdaNode("lambda1", lambda1, compose.WithStatePreHandler(lamada1PreHandler))
	if err != nil {
		panic(err)
	}
	err = g.AddLambdaNode("lambda2", lambda2, compose.WithStatePreHandler(lamada2PreHandler))
	if err != nil {
		panic(err)
	}
	err = g.AddChatModelNode("chatmodel", chatmodel)
	if err != nil {
		panic(err)
	}
	//err = g.AddToolsNode("tools", ToolNode)
	//if err != nil {
	//	panic(err)
	//}
	//分支连接
	err = g.AddEdge(compose.START, "lambda0")
	if err != nil {
		panic(err)
	}
	//加入分支
	err = g.AddBranch("lambda0", compose.NewGraphBranch(func(ctx context.Context, in map[string]string) (endNode string, err error) {
		if in["role"] == "毫猫" {
			return "lambda1", nil
		} else if in["role"] == "耄耋" {
			return "lambda2", nil
		}
		return "lambda2", nil
	}, map[string]bool{"lambda1": true, "lambda2": true}))
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("lambda1", "chatmodel")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("lambda2", "chatmodel")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("chatmodel", compose.END)
	if err != nil {
		panic(err)
	}
	//err = g.AddEdge("chatmodel", "tools")
	//if err != nil {
	//	panic(err)
	//}
	//err = g.AddEdge("tools", compose.END)
	//if err != nil {
	//	panic(err)
	//}

	outsidegraph := compose.NewGraph[map[string]string, string]()

	writeLambda := compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (output string, err error) {
		f, err := os.OpenFile("orc_graph_withgraph.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}
		defer f.Close()
		if _, err := f.WriteString(input.Content + "\n---\n"); err != nil {
			return "", err
		}
		return "已经写入文件，请前往文件内查看内容", nil
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
		"content": content,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
	return answer
}
