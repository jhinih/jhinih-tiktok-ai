package logic

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/manager"
	"Tiktok/model"
	"Tiktok/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

const (
	REDIS_CHAT_MESSAGE_SET = "chat:message:set"
)

type WebsocketLogic struct {
	conn   *websocket.Conn
	userID int64
	meta   *model.ConnMeta
}

func NewWebsocketLogic(conn *websocket.Conn, meta *model.ConnMeta) *WebsocketLogic {
	return &WebsocketLogic{
		conn:   conn,
		userID: manager.WebsocketManager.Clients[conn],
		meta:   meta,
	}
}

func (l *WebsocketLogic) HandleMessage(message string) {
	var Cindata model.InMessage
	err := json.Unmarshal([]byte(message), &Cindata)
	if err != nil {
		zlog.Warnf("websocket 接受消息格式错误: %s", err)
		return
	}
	zlog.Debugf("websocket 收到消息: %v", Cindata)
	// 分类解析消息
	switch Cindata.Type {
	case "chat":
		l.handleTypeChatMessage(Cindata.Content)
	case "ai":
		l.handleTypeAIMessage(Cindata.Content)
	case "common_ai":
		l.handleTypeAICommonMessage(Cindata.Content)
	case "history":
		l.handleTypeGetHistory(Cindata.Content)
	default:
		zlog.Warnf("websocket 未知消息类型: %s", Cindata.Type)
	}
	return
}

func (l *WebsocketLogic) handleTypeGetHistory(content string) {
	// 处理消息内容
	var data types.GetHistoryRequest
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		zlog.Warnf("websocket 接受消息格式错误: %s", err)
		return
	}
	// 从 redis 中获取历史消息
	res, err := global.Rdb.ZRevRangeByScore(context.Background(), REDIS_CHAT_MESSAGE_SET, &redis.ZRangeBy{
		Min: "0",
		Max: strconv.FormatInt(data.Before-1, 10),
	}).Result()

	// 打包前 data.Count 条消息
	for i := 0; i < int(data.Count) && i < len(res); i++ {
		// 解析消息内容
		var resp model.OutMessage
		err := json.Unmarshal([]byte(res[i]), &resp)
		if err != nil {
			zlog.Warnf("websocket 解析消息格式错误: %s", err)
			continue
		}
		resp.Type = "history"
		// 再次转换json
		respJson, err := json.Marshal(resp)
		if err != nil {
			zlog.Warnf("websocket 打包消息格式错误: %s", err)
			continue
		}
		// 单发消息
		msg := model.ChatMessage{
			ToType:  "user",
			To:      l.userID,
			Content: string(respJson),
		}
		manager.WebsocketManager.Broadcast <- msg
	}
}

func (l *WebsocketLogic) handleTypeChatMessage(content string) {
	// 处理消息内容
	var Contentdata model.ChatMessage
	err := json.Unmarshal([]byte(content), &Contentdata)
	if err != nil {
		zlog.Warnf("websocket 接受消息格式错误: %s", err)
		return
	}
	// 打包信息内容
	id := global.SnowflakeNode.Generate().Int64()
	resp := model.OutMessage{
		Type:      "chat",
		ID:        id,                  //消息ID
		Content:   Contentdata.Content, //内容
		UserID:    strconv.FormatInt(l.userID, 10),
		Timestamp: time.Now().UnixMilli(),
	}
	SaveChatMessage(context.Background(), resp)
	respJson, err := json.Marshal(resp)
	if err != nil {
		zlog.Warnf("websocket 打包消息格式错误: %s", err)
		return
	}
	// 处理消息
	msg := model.ChatMessage{
		ToType:  Contentdata.ToType, //all/user
		Content: string(respJson),
		To:      Contentdata.To,
	}
	manager.WebsocketManager.Broadcast <- msg
}
func (l *WebsocketLogic) handleTypeAIMessage(content string) {
	// 处理消息内容
	var Contentdata model.ChatMessage
	err := json.Unmarshal([]byte(content), &Contentdata)
	if err != nil {
		zlog.Warnf("websocket 接受消息格式错误: %s", err)
		return
	}
	// 打包信息内容
	id := global.SnowflakeNode.Generate().Int64()
	resp := model.OutMessage{
		Type:      "ai",
		ID:        id,
		Content:   Contentdata.Content,
		UserID:    strconv.FormatInt(l.userID, 10),
		Timestamp: time.Now().UnixMilli(),
	}
	SaveChatMessage(context.Background(), resp)
	respJson, err := json.Marshal(resp)
	if err != nil {
		zlog.Warnf("websocket 打包消息格式错误: %s", err)
		return
	}
	// 处理消息
	msg := model.ChatMessage{
		ToType:  Contentdata.ToType,
		Content: string(respJson),
		To:      2,
	}
	manager.WebsocketManager.Broadcast <- msg
	go l.AIChat(Contentdata.Content)
}
func (l *WebsocketLogic) handleTypeAICommonMessage(content string) {
	// 处理消息内容
	var Contentdata model.ChatMessage
	err := json.Unmarshal([]byte(content), &Contentdata)
	if err != nil {
		zlog.Warnf("websocket 接受消息格式错误: %s", err)
		return
	}
	// 打包信息内容
	id := global.SnowflakeNode.Generate().Int64()
	resp := model.OutMessage{
		Type:      "common_ai",
		ID:        id,
		Content:   Contentdata.Content,
		UserID:    strconv.FormatInt(l.userID, 10),
		Timestamp: time.Now().UnixMilli(),
	}
	SaveChatMessage(context.Background(), resp)
	respJson, err := json.Marshal(resp)
	if err != nil {
		zlog.Warnf("websocket 打包消息格式错误: %s", err)
		return
	}
	// 处理消息
	msg := model.ChatMessage{
		ToType:  Contentdata.ToType,
		Content: string(respJson),
		To:      1,
	}
	manager.WebsocketManager.Broadcast <- msg
	go l.AICommonChat(Contentdata.Content)
}

func SaveChatMessage(ctx context.Context, message model.OutMessage) {
	// 转化 json 格式
	messageJson, err := json.Marshal(message)
	if err != nil {
		zlog.Errorf("websocket 打包消息格式错误: %s", err)
		return
	}
	//zlog.Debugf("websocket 保存消息: %s", messageJson)
	// 保存消息到 redis
	err = global.Rdb.ZAdd(ctx, REDIS_CHAT_MESSAGE_SET, redis.Z{
		Score:  float64(message.Timestamp),
		Member: messageJson,
	}).Err()
	if err != nil {
		zlog.Errorf("websocket 保存消息到 redis 失败: %s", err)
		return
	}

	// 如果 redis 中的消息数量超过 50 条，清理最早的消息
	for global.Rdb.ZCount(ctx, REDIS_CHAT_MESSAGE_SET, "-inf", "+inf").Val() > 50 {
		err = global.Rdb.ZRemRangeByRank(ctx, REDIS_CHAT_MESSAGE_SET, 0, 0).Err()
		if err != nil {
			zlog.Errorf("websocket 清理消息到 redis 失败: %s", err)
			return
		}
	}
}

func (l *WebsocketLogic) AIChat(content string) {

	c := context.Background()
	req := types.AIRequest{
		Ask: content,
	}
	airesponse, err := NewAILogic().AI(c, req)
	if err != nil {
		zlog.Errorf("AI调用失败: %s", err)
	}
	fmt.Println(airesponse)
	// 打包信息内容
	id := global.SnowflakeNode.Generate().Int64()
	resp := model.OutMessage{
		Type:      "ai",
		ID:        id,
		Content:   airesponse.Anser,
		UserID:    "2",
		Timestamp: time.Now().UnixMilli(),
	}
	SaveChatMessage(context.Background(), resp)
	respJson, err := json.Marshal(resp)
	if err != nil {
		zlog.Warnf("websocket 打包消息格式错误: %s", err)
		return
	}
	// 处理消息
	msg := model.ChatMessage{
		ToType:  "user",
		To:      l.userID,
		Content: string(respJson),
	}
	manager.WebsocketManager.Broadcast <- msg
}

func (l *WebsocketLogic) AICommonChat(content string) {
	c := context.Background()
	req := types.AIRequest{
		Ask: content,
	}
	airesponse, err := NewAILogic().CommonAI(c, req)
	if err != nil {
		zlog.Errorf("AI调用失败: %s", err)
	}

	// 打包信息内容
	id := global.SnowflakeNode.Generate().Int64()
	resp := model.OutMessage{
		Type:      "ai",
		ID:        id,
		Content:   airesponse.Anser,
		UserID:    "1",
		Timestamp: time.Now().UnixMilli(),
	}
	SaveChatMessage(context.Background(), resp)
	respJson, err := json.Marshal(resp)
	if err != nil {
		zlog.Warnf("websocket 打包消息格式错误: %s", err)
		return
	}
	// 处理消息
	msg := model.ChatMessage{
		ToType:  "user",
		To:      l.userID,
		Content: string(respJson),
	}
	manager.WebsocketManager.Broadcast <- msg
}
