package logic

import (
	"Tiktok/types"
	"Tiktok/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AILogic struct {
}

func NewAILogic() *AILogic {
	return &AILogic{}
}

func crawler(url string, payload map[string]any) (string, error) {
	// 要 POST 的 JSON 数据
	raw, _ := json.Marshal(payload)

	// 创建请求
	request, _ := http.NewRequest("POST", url, bytes.NewReader(raw))
	request.Header.Set("Content-Type", "application/json")
	//repository.Header.Set("Authorization", "Bearer "+jwt)

	// 发送请求
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return "爬虫失败", err
	}
	defer response.Body.Close()

	// 6. 打印结果
	type Resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Answer string `json:"answer"`
		} `json:"data"`
	}

	var s Resp
	json.NewDecoder(response.Body).Decode(&s)
	fmt.Printf("status: %d\nbody  :  %s\n", response.StatusCode, s.Data.Answer)
	resp := s.Data.Answer
	return resp, nil
}

func (l *AILogic) CommonAI(ctx context.Context, req types.AIRequest) (resp types.AIResponse, err error) {
	defer utils.RecordTime(time.Now())()
	// 1. 目标接口
	url := "http://localhost:8888/api/ai/aiCommonChat"

	//// 2. 你的 JWT
	//jwt := req.Atoken

	// 3. 要 POST 的 JSON 数据
	payload := map[string]any{"ask": req.Ask}
	raw, _ := json.Marshal(payload)

	// 4. 创建请求
	request, _ := http.NewRequest("POST", url, bytes.NewReader(raw))
	request.Header.Set("Content-Type", "application/json")
	//repository.Header.Set("Authorization", "Bearer "+jwt)

	// 5. 发送
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer response.Body.Close()

	// 6. 打印结果
	type Resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Answer string `json:"answer"`
		} `json:"data"`
	}

	var s Resp
	json.NewDecoder(response.Body).Decode(&s)
	fmt.Printf("status: %d\nbody  :  %s\n", response.StatusCode, s.Data.Answer)
	resp = types.AIResponse{
		Anser: s.Data.Answer,
	}
	return resp, nil
}
func (l *AILogic) VideoAI(ctx context.Context, req types.AIRequest) (resp types.AIResponse, err error) {
	defer utils.RecordTime(time.Now())()
	// 1. 目标接口
	url := "http://localhost:8888/api/ai/aiVideoChat"
	//
	//// 2. 你的 JWT
	//jwt := req.Atoken

	// 3. 要 POST 的 JSON 数据
	payload := map[string]any{"ask": req.Ask}
	raw, _ := json.Marshal(payload)

	// 4. 创建请求
	request, _ := http.NewRequest("POST", url, bytes.NewReader(raw))
	request.Header.Set("Content-Type", "application/json")
	//repository.Header.Set("Authorization", "Bearer "+jwt)

	// 5. 发送
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer response.Body.Close()

	// 6. 打印结果
	type Resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Answer string `json:"answer"`
		} `json:"data"`
	}

	var s Resp
	json.NewDecoder(response.Body).Decode(&s)
	fmt.Printf("status: %d\nbody  :  %s\n", response.StatusCode, s.Data.Answer)
	resp = types.AIResponse{
		Anser: s.Data.Answer,
	}
	return resp, nil
}

func (l *AILogic) SendCodeAI(ctx context.Context, req types.AIRequest) (resp types.AIResponse, err error) {
	defer utils.RecordTime(time.Now())()
	url := "http://localhost:8888/api/ai/AISendCodeChat"
	payload := map[string]any{"ask": req.Ask}
	response, err := crawler(url, payload)
	if err != nil {
		fmt.Println("发送请求失败:", err)
	}
	resp = types.AIResponse{
		Anser: response,
	}
	return resp, nil
}
func (l *AILogic) GetUserInfoAI(ctx context.Context, req types.AIRequest) (resp types.AIResponse, err error) {
	defer utils.RecordTime(time.Now())()
	url := "http://localhost:8888/api/ai/AIGetUserInfoChat"
	payload := map[string]any{"ask": req.Ask}
	response, err := crawler(url, payload)
	if err != nil {
		fmt.Println("发送请求失败:", err)
	}
	resp = types.AIResponse{
		Anser: response,
	}
	return resp, nil
}

func (l *AILogic) GetVideoAI(ctx context.Context, req types.AIRequest) (resp types.AIResponse, err error) {
	defer utils.RecordTime(time.Now())()
	url := "http://localhost:8888/api/ai/AIGetVideoChat"
	payload := map[string]any{"ask": req.Ask}
	response, err := crawler(url, payload)
	if err != nil {
		fmt.Println("发送请求失败:", err)
	}
	resp = types.AIResponse{
		Anser: response,
	}
	return resp, nil
}

func (l *AILogic) AllAI(ctx context.Context, req types.AIRequest) (resp types.AIResponse, err error) {
	defer utils.RecordTime(time.Now())()
	url := "http://localhost:8888/api/ai/aiAllChat"
	payload := map[string]any{"ask": req.Ask}
	response, err := crawler(url, payload)
	if err != nil {
		fmt.Println("发送请求失败:", err)
	}
	resp = types.AIResponse{
		Anser: response,
	}
	return resp, nil
}

func (l *AILogic) AI(ctx context.Context, req types.AIRequest) (resp types.AIResponse, err error) {
	defer utils.RecordTime(time.Now())()
	url := "http://localhost:8888/api/ai/aiAllChat"
	payload := map[string]any{"ask": req.Ask}
	response, err := crawler(url, payload)
	if err != nil {
		fmt.Println("请求失败:", err)
	}
	resp = types.AIResponse{
		Anser: response,
	}
	return resp, nil
}
