package aiUtils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetToken() (string, error) {
	// 1. 目标接口
	url := "http://localhost:8080/api/login/login"

	// 3. 要 POST 的 JSON 数据
	payload := map[string]any{"email": "1@gmail.com", "password": "123456"}
	raw, _ := json.Marshal(payload)

	// 4. 创建请求
	req, _ := http.NewRequest("POST", url, bytes.NewReader(raw))
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
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Atoken string `json:"atoken"`
			Rtoken string `json:"rtoken"`
		} `json:"data"`
	}

	var s Resp
	json.NewDecoder(resp.Body).Decode(&s)
	return s.Data.Atoken, nil
}
