package login

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"net/http"
)

type InputParams struct {
	Email string `json:"email" jsonschema:"description=the Email address"`
}

func SendCode(_ context.Context, params *InputParams) (string, error) {
	// 1. 目标接口
	url := "http://localhost:8080/api/login/send-code"

	// 3. 要 POST 的 JSON 数据
	payload := map[string]any{"email": params.Email}
	raw, _ := json.Marshal(payload)

	// 4. 创建请求
	req, _ := http.NewRequest("POST", url, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	// 5. 发送
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("repository err:", err)
		return "发送验证码请求失败", err
	}
	defer resp.Body.Close()

	// 6. 打印结果
	type Resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Answer string `json:"answer"`
		} `json:"data"`
	}

	var s Resp
	json.NewDecoder(resp.Body).Decode(&s)
	fmt.Printf("status: %d\nbody  :  %s\n", resp.StatusCode, s.Data.Answer)

	return "发送验证码请求成功", nil
}

func CreateSendCodeTool() tool.InvokableTool {
	sendCodeTool := utils.NewTool(&schema.ToolInfo{
		Name: "send_code",
		Desc: "send code to the Email address",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"email": &schema.ParameterInfo{
					Type:     schema.String,
					Desc:     "the Email address",
					Required: true,
				},
			},
		),
	}, SendCode)
	return sendCodeTool
}
