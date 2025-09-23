package ais

import (
	"context"
	"eino_gozero/common/ai_common/jhinih_model/chatmodel"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"io"
	"time"
)

func Test1() {
	startTime := time.Now()
	fmt.Printf("程序开始执行时间: %s\n", startTime.Format("2006-01-02 15:04:05.000"))

	ctx := context.Background()
	arkModel := chatmodel.NewArkChatModel(ctx)
	addtool := GetAddTool()
	subtool := GetSubTool()
	analyzetool := GetAnalyzeTool()
	persona := `#Character:
	你是一个幼儿园老师，会同时判断题目难易程度，给出问题的答案
	`
	// toolCallChecker 用于检查从流中读取的消息是否包含工具调用。
	// 它会持续从流中接收消息，直到找到一个包含工具调用的消息或流结束。
	toolCallChecker := func(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
		// 确保在函数退出时关闭流读取器，以防资源泄漏。
		defer sr.Close()
		// 无限循环，用于持续从流中读取消息。
		for {
			// 从流中接收一条消息。如果流中没有新消息，此调用会阻塞。
			msg, err := sr.Recv()
			// 检查接收过程中是否发生错误。
			if err != nil {
				// 如果错误是 io.EOF，表示流已正常结束，没有更多消息了。
				if errors.Is(err, io.EOF) {
					// 正常结束，跳出循环。
					break
				}

				// 如果是其他类型的错误，则直接返回错误。
				return false, err
			}

			// 检查收到的消息是否包含任何工具调用。
			if len(msg.ToolCalls) > 0 {
				// 如果找到工具调用，立即返回 true，表示检查成功。
				return true, nil
			}
		}
		// 如果完整遍历了流而没有找到任何工具调用，则返回 false。
		return false, nil
	}
	raAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: arkModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools:               []tool.BaseTool{addtool, subtool, analyzetool},
			ExecuteSequentially: false,
		},
		StreamToolCallChecker: toolCallChecker,
	})
	if err != nil {
		fmt.Printf("Agent创建失败: %v", err)
		return
	}
	//构建输入模板schema
	chatmsg := []*schema.Message{
		{
			Role:    schema.System,
			Content: persona,
		},
		{
			Role:    schema.User,
			Content: "请同时告诉我183+192-90这道题的难易程度和答案",
		},
	}
	//流式调用
	//添加了loggerCallback的回调函数
	sr, err := raAgent.Stream(ctx, chatmsg, agent.WithComposeOptions(compose.WithCallbacks(&loggerCallback{})))
	if err != nil {
		fmt.Printf("流式调用失败: %v", err)
		return
	}
	//构建最终回答
	finalContent := ""
	for {
		msg, err := sr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// finish
				break
			}
			fmt.Printf("接受流式输出失败: %v", err)
			return
		}
		finalContent += msg.Content
		fmt.Printf("%v", msg.Content)
	}

	fmt.Printf("\n\n===== 失败回答 =====\n\n")
	fmt.Printf("%s", finalContent)
	fmt.Printf("\n\n===== 完成 =====\n")

	// 计算并打印运行时间
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("\n程序结束执行时间: %s\n", endTime.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("总运行时间: %v\n", duration)
	fmt.Printf("总运行时间(毫秒): %.2f ms\n", float64(duration.Nanoseconds())/1000000)
}

type loggerCallback struct {
	callbacks.HandlerBuilder
}

func (cb *loggerCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	fmt.Println("=========开始=========")
	inputStr, _ := json.MarshalIndent(input, "", "  ")
	fmt.Println(string(inputStr))
	return ctx
}

func (cb *loggerCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	fmt.Println("=========结束=========")
	outputStr, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(outputStr))
	return ctx
}

func (cb *loggerCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	fmt.Println("=========错误=========")
	fmt.Println(err)
	return ctx
}

func (cb *loggerCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo,
	output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {

	var graphInfoName = react.GraphName

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("流式输出结束错误:", err)
			}
		}()

		defer output.Close()

		for {
			frame, err := output.Recv()
			if errors.Is(err, io.EOF) {
				// finish
				break
			}
			if err != nil {
				fmt.Printf("错误: %s\n", err)
				return
			}

			s, err := json.Marshal(frame)
			if err != nil {
				fmt.Printf("错误: %s\n", err)
				return
			}

			if info.Name == graphInfoName { // 仅打印 graph 的输出, 否则每个 stream 节点的输出都会打印一遍
				fmt.Printf("%s: %s\n", info.Name, string(s))
			}
		}

	}()
	return ctx
}

func (cb *loggerCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo,
	input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	defer input.Close()
	return ctx
}
