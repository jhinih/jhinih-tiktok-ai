package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/response"
	"Tiktok/types"
	"Tiktok/utils/jwtUtils"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

//func Chats(c *gin.Context) {
//	ctx := zlog.GetCtxFromGin(c)
//	req, err := types.BindRequest[types.ChatRequest](c)
//	if err != nil {
//		response.Response(c, nil, err)
//		return
//	}
//	zlog.CtxInfof(ctx, "Chat repository: %v", req)
//
//	// Implement chat logic here
//	resp, err := logic.NewChatLogic().Chat(c, req)
//	response.Response(c, resp, err)
//}
//
//func RedisMsg(c *gin.Context) {
//
//	ctx := zlog.GetCtxFromGin(c)
//	req, err := types.BindRequest[types.RedisMsgRequest](c)
//	if err != nil {
//		return
//	}
//	zlog.CtxInfof(ctx, "发送个人消息请求: %v", req)
//	resp, err := logic.NewChatLogic().RedisMsg(ctx, req)
//	response.Response(c, resp, err)
//}
//

//func SendMsg(c *gin.Context) {
//	ctx := zlog.GetCtxFromGin(c)
//	req, err := types.BindRequest[types.SendMsgRequest](c)
//	if err != nil {
//		return
//	}
//	zlog.CtxInfof(ctx, "发送消息请求: %v", req)
//	resp, err := logic.NewChatLogic().SendMsg(req)
//	response.Response(c, resp, err)
//}
//func SendUserMsg(c *gin.Context) {
//	ctx := zlog.GetCtxFromGin(c)
//	req, err := types.BindRequest[types.SendMsgRequest](c)
//	if err != nil {
//		return
//	}
//	zlog.CtxInfof(ctx, "发送个人消息请求: %v", req)
//	resp, err := logic.NewChatLogic().SendMsg(req)
//	response.Response(c, resp, err)
//}
//func SendGroupMsg(c *gin.Context) {
//	ctx := zlog.GetCtxFromGin(c)
//	req, err := types.BindRequest[types.SendGroupMsgRequest](c)
//	if err != nil {
//		return
//	}
//	zlog.CtxInfof(ctx, "发送群消息请求: %v", req)
//	resp, err := logic.NewChatLogic().SendGroupMsg(req)
//	response.Response(c, resp, err)
//}

func Connection(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.ConnectionRequest](c)
	if err != nil {
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "连接websocket请求: %v", req)
	zlog.CtxInfof(ctx, "WebSocket连接: %s", c.Request.RemoteAddr)
	req.UserID = jwtUtils.GetUserId(c)
	req.UserName = req.UserID
	req.Role = req.Role
	logic.NewChatLogic().Connection(c, req)
	zlog.CtxInfof(ctx, "WebSocket连接: %s", c.Request.RemoteAddr)
}

//func HandleSendMessage(c *gin.Context) {
//	ctx := zlog.GetCtxFromGin(c)
//	req, err := types.BindRequest[types.SendMessageRequest](c)
//	if err != nil {
//		response.Response(c, req, err)
//		return
//	}
//	req.UserId = jwtUtils.GetUserId(c)
//	req.UserName = jwtUtils.GetUserName(c)
//	zlog.CtxInfof(ctx, "发送消息请求: %v", req)
//	logic.NewChatLogic().HandleSendMessage(c, req)
//}
//
//func SendUserMessage(c *gin.Context) {
//	ctx := zlog.GetCtxFromGin(c)
//	req, err := types.BindRequest[types.SendMessageRequest](c)
//	if err != nil {
//		response.Response(c, req, err)
//		return
//	}
//	req.UserId = jwtUtils.GetUserId(c)
//	req.UserName = jwtUtils.GetUserName(c)
//	zlog.CtxInfof(ctx, "发送个人消息请求: %v", req)
//	logic.NewChatLogic().SendUserMessage(c, req)
//}
//func SendGroupMessage(c *gin.Context) {
//	ctx := zlog.GetCtxFromGin(c)
//	req, err := types.BindRequest[types.SendMessageRequest](c)
//	if err != nil {
//		response.Response(c, req, err)
//		return
//	}
//	req.UserId = jwtUtils.GetUserId(c)
//	req.UserName = jwtUtils.GetUserName(c)
//	zlog.CtxInfof(ctx, "发送群聊消息请求: %v", req)
//	logic.NewChatLogic().SendGroupMessage(c, req)
//}
//
//// 加载离线信息
//func DeliverOfflineMessages(c *gin.Context) {
//	ctx := zlog.GetCtxFromGin(c)
//	req, err := types.BindRequest[types.DeliverOfflineMessagesRequest](c)
//	if err != nil {
//		return
//	}
//	req.UserID = jwtUtils.GetUserId(c)
//	zlog.CtxInfof(ctx, "加载离线信息: %v", req)
//	logic.NewChatLogic().DeliverOfflineMessages(c, req)
//	response.Response(c, nil, err)
//}

// 生成一次性 30 秒 WebSocket ticket
func GenWSTicket(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GenWSTicketRequest](c)
	if err != nil {
		return
	}
	req.UserID = jwtUtils.GetUserId(c)
	role, _ := strconv.Atoi(req.Role)
	zlog.CtxInfof(ctx, "生成一次性 30 秒 WebSocket ticket: %v", req)
	resp := types.GenWSTicketResponse{}
	resp.Ticket, _ = jwtUtils.GenToken(req.UserID, req.UserName, role, 30*time.Second)
	response.Response(c, resp, err)
}
