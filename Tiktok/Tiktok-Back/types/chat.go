package types

import (
	"github.com/gorilla/websocket"
)

//	type ChatRequest struct {
//		UserId string `json:"userId" binding:"-"`
//	}
//
// type ChatResponse struct {
// }
type ConnectionRequest struct {
	UserID   string `json:"userID" binding:"-"`
	UserName string `json:"userName" binding:"-"`
	Role     string `json:"role" binding:"-"`
}
type ConnectionResponse struct {
}

type SendMessageRequest struct {
	UserId     string `json:"userID" binding:"-"`
	UserName   string `json:"userName" binding:"-"`
	TargetId   string `json:"targetID" binding:"required"`
	TargetName string `json:"targetName"`
	Type       string `json:"type"`
	Media      string `json:"media"`
	Content    string `json:"content" binding:"required"`
	CreateTime string `json:"createTime"`
	ReadTime   string `json:"readTime"`
	Pic        string `json:"pic"`
	Url        string `json:"url"`
	Desc       string `json:"desc"`
	Amount     string `json:"amount"`
}
type SendMessageResponse struct {
}

type DeliverOfflineMessagesRequest struct {
	UserID string          `json:"userID" binding:"-"`
	Conn   *websocket.Conn `json:"-"`
}
type DeliverOfflineMessagesResponse struct {
}
type GenWSTicketRequest struct {
	UserID   string `json:"userID" binding:"-"`
	UserName string `json:"userName"`
	Role     string `json:"role"`
}
type GenWSTicketResponse struct {
	Ticket string `json:"ticket"`
}
