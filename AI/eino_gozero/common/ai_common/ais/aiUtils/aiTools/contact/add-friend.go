package contact

import (
	"bytes"
	"context"
	"eino_gozero/common/ai_common/ais/aiUtils"
	_ "eino_gozero/config"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"net/http"
)

type InputParams struct {
	TargetName string `json:"user_name" jsonschema:"description=the TargetName of the user who you want to add friend"`
}

func AddFriend(_ context.Context, params *InputParams) (string, error) {
	// 1. 目标接口
	url := "http://localhost:8080/api/contact/add-friend"
	// 2. 获取 JWT
	jwt, _ := aiUtils.GetToken()

	// 3. 要 POST 的 JSON 数据
	payload := map[string]any{"email": params.TargetName}
	raw, _ := json.Marshal(payload)

	// 4. 创建请求
	req, _ := http.NewRequest("POST", url, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Cache-Control", "no-cache")

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

func CreateAddFriendTool() tool.InvokableTool {
	AddFriendTool := utils.NewTool(&schema.ToolInfo{
		Name: "send_code",
		Desc: "send code to the Email address",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"user_name": &schema.ParameterInfo{
					Type:     schema.String,
					Desc:     "the Email address",
					Required: true,
				},
			},
		),
	}, AddFriend)
	return AddFriendTool
}
