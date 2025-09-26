package aiUtils

import (
	"context"
	"eino_gozero/common/ai_common/jhinih_model/chatmodel"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func AIWithToolsJson(name, content string, tools []tool.BaseTool, infos []*schema.ToolInfo) *compose.Graph[map[string]string, string] {
	err := godotenv.Load() // 加载环境变量
	if err != nil {
		log.Fatal("加载环境变量失败")
	}
	ctx := context.Background()

	g := compose.NewGraph[map[string]string, *schema.Message](
		compose.WithGenLocalState(GenFunc),
	)
	//具体节点
	lambda1 := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: content,
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	//节点前处理
	lamada1PreHandler := func(ctx context.Context, input map[string]string, state *State) (output map[string]string, err error) {
		//input["content"] = input["content"] + state.History["毫猫"].(string)
		return input, nil
	}

	//AI引入
	chatmodel := chatmodel.NewArkChatModel(ctx)

	//工具引入
	chatmodel.BindTools(infos)
	ToolNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: tools,
	})
	if err != nil {
		log.Fatal(err)
	}

	convertLambda := compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (output *schema.Message, err error) {
		if len(input) == 0 {
			return nil, fmt.Errorf("tools returned empty message list")
		}
		// 取最后一个消息作为最终结果
		return input[len(input)-1], nil
	})

	//加入节点
	err = g.AddLambdaNode("lambda1", lambda1, compose.WithStatePreHandler(lamada1PreHandler))
	if err != nil {
		panic(err)
	}
	err = g.AddChatModelNode("chatmodel", chatmodel)
	if err != nil {
		panic(err)
	}
	err = g.AddToolsNode("tools", ToolNode)
	if err != nil {
		panic(err)
	}
	err = g.AddLambdaNode("convert", convertLambda)
	if err != nil {
		panic(err)
	}

	//连接节点
	err = g.AddEdge(compose.START, "lambda1")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("lambda1", "chatmodel")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("chatmodel", "tools")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("tools", "convert")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("convert", compose.END)
	if err != nil {
		panic(err)
	}

	outsidegraph := compose.NewGraph[map[string]string, string]()

	writeLambda := compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (output string, err error) {
		f, err := os.OpenFile(name+".md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}
		defer f.Close()
		if _, err := f.WriteString(input.Content + "\n---\n"); err != nil {
			return "", err
		}
		return input.Content, nil
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

	return outsidegraph
}
