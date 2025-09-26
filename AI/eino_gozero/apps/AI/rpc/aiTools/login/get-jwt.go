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

type GetJWTInputParams struct {
}

func GetJWT(_ context.Context, params *GetJWTInputParams) (string, error) {
	// 1. 目标接口
	url := "http://localhost:8080/api/login/get-token"

	// 3. 要 POST 的 JSON 数据
	payload := map[string]any{"email": "params.Email"}
	raw, _ := json.Marshal(payload)

	// 4. 创建请求
	req, _ := http.NewRequest("GET", url, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	// 5. 发送
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("repository err:", err)
		return "获取Token失败", err
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

	return s.Data.Answer, nil
}

func CreateGetJWTTool() tool.InvokableTool {
	GetJWTTool := utils.NewTool(&schema.ToolInfo{
		Name: "GetJWT",
		Desc: "get jwt tool(always is atoken)",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{},
		),
	}, GetJWT)
	return GetJWTTool
}
