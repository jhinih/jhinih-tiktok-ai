package logic

import (
	"Tiktok/global"
	"Tiktok/model"
	"Tiktok/repository"
	"Tiktok/types"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 完整WebSocket聊天服务实现
type ChatLogic struct {
	// WebSocket配置
	upgrader websocket.Upgrader
	// 连接管理
	clients   map[string]*Client
	clientsMu sync.RWMutex
	// Redis客户端
	redisClient *redis.Client
	// 消息广播通道
	broadcast chan model.Message
}

// 客户端连接
type Client struct {
	conn     *websocket.Conn
	userID   string
	userName string
	role     int
	lastPing time.Time
}

// 初始化聊天服务
func NewChatLogic() *ChatLogic {
	ChatLogic := &ChatLogic{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // 生产环境应校验来源
			},
		},
		clients:     make(map[string]*Client),
		redisClient: global.Rdb,
		broadcast:   make(chan model.Message, 100),
	}
	// 启动消息广播
	go ChatLogic.RunBroadcaster()
	// 启动连接健康检查
	go ChatLogic.HealthCheck()
	return ChatLogic
}

// 消息广播器
func (l *ChatLogic) RunBroadcaster() {
	for msg := range l.broadcast {
		// 这里可以实现全局广播逻辑
		log.Printf("广播消息: %+v", msg)
	}
}

// 连接健康检查
func (l *ChatLogic) HealthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		var staleClients []string
		l.clientsMu.RLock()
		for userID, client := range l.clients {
			if now.Sub(client.lastPing) > 90*time.Second {
				staleClients = append(staleClients, userID)
			}
		}
		l.clientsMu.RUnlock()
		// 清理不活跃连接
		for _, userID := range staleClients {
			l.clientsMu.Lock()
			if client, ok := l.clients[userID]; ok {
				err := client.conn.Close()
				if err != nil {
					return
				}
				delete(l.clients, userID)
				log.Printf("清理不活跃连接: %s", userID)
			}
			l.clientsMu.Unlock()
		}
	}
}

// WebSocket连接处理
func (l *ChatLogic) Connection(c *gin.Context, req types.ConnectionRequest) {
	// 升级为WebSocket连接
	UserID := req.UserID
	conn, err := l.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket升级失败:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	role, _ := strconv.Atoi(req.Role)
	// 注册客户端
	client := &Client{
		conn:     conn,
		userID:   UserID,
		userName: req.UserName,
		role:     role,
		lastPing: time.Now(),
	}
	l.RegisterClient(client)
	defer l.UnregisterClient(UserID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	l.DeliverOfflineMessages(ctx, UserID, conn)

	// 消息处理循环
	for {
		var msg model.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("用户 %s 断开连接: %v", UserID, err)
			break
		}

		// 更新最后活跃时间
		client.lastPing = time.Now()

		// 处理不同类型消息
		switch msg.Type {
		case 1:
			l.SendUserMessageConn(c, msg)
		case 2:
			l.SendGroupMessageConn(c, msg)
		case 3:
			// 心跳响应
			err := conn.WriteJSON(model.Message{Type: 3})
			if err != nil {
				return
			}
		}
	}
}

// 处理私聊消息
func (l *ChatLogic) SendUserMessageConn(c *gin.Context, msg model.Message) {
	// 存储消息
	l.StoreMessageConn(c, msg)

	// 尝试实时推送
	l.clientsMu.RLock()
	TargetId, ok := l.clients[msg.TargetId]
	l.clientsMu.RUnlock()

	if ok {
		if err := TargetId.conn.WriteJSON(msg); err != nil {
			log.Printf("实时消息推送失败: %v", err)
		}
	}
}

// 处理群组消息
func (l *ChatLogic) SendGroupMessageConn(c *gin.Context, msg model.Message) {
	// 存储消息
	l.StoreMessageConn(c, msg)

	//获取群组成员
	clean := strings.Trim(msg.UserId, `"`)
	Id, _ := strconv.ParseInt(clean, 10, 64)
	members, _ := repository.NewContactRequest(global.DB).GetGroupUsers(Id)

	// 广播给在线成员
	for _, memberID := range members {
		if memberID.ID == Id {
			continue // 不发送给自己
		}

		l.clientsMu.RLock()
		client, ok := l.clients[strconv.FormatInt(memberID.ID, 10)]
		l.clientsMu.RUnlock()

		if ok {
			if err := client.conn.WriteJSON(msg); err != nil {
				log.Printf("群消息推送失败(%d): %v", memberID.ID, err)
			}
		}
	}
}

//// 处理私聊消息
//func (l *ChatLogic) SendUserMessage(c *gin.Context, req types.SendMessageRequest) {
//	msg := model.Message{
//		Model:      gorm.Model{},
//		UserId:     req.UserId,
//		UserName:   req.UserName,
//		TargetId:   req.TargetId,
//		TargetName: req.TargetName,
//		Type:       req.Type,
//		Media:      req.Media,
//		Content:    req.Content,
//		CreateTime: uint64(time.Now().Unix()),
//		ReadTime:   req.ReadTime,
//		Pic:        req.Pic,
//		Url:        req.Url,
//		Desc:       req.Desc,
//		Amount:     req.Amount,
//	}
//	// 存储消息
//	l.StoreMessage(c, req)
//
//	// 尝试实时推送
//	l.clientsMu.RLock()
//	TargetId, ok := l.clients[msg.TargetId]
//	l.clientsMu.RUnlock()
//
//	if ok {
//		if err := TargetId.conn.WriteJSON(msg); err != nil {
//			log.Printf("实时消息推送失败: %v", err)
//		}
//	}
//}
//
//// 处理群组消息
//func (l *ChatLogic) SendGroupMessage(c *gin.Context, req types.SendMessageRequest) {
//	msg := model.Message{
//		Model:      gorm.Model{},
//		UserId:     req.UserId,
//		UserName:   req.UserName,
//		TargetId:   req.TargetId,
//		TargetName: req.TargetName,
//		Type:       req.Type,
//		Media:      req.Media,
//		Content:    req.Content,
//		CreateTime: uint64(time.Now().Unix()),
//		ReadTime:   req.ReadTime,
//		Pic:        req.Pic,
//		Url:        req.Url,
//		Desc:       req.Desc,
//		Amount:     req.Amount,
//	}
//	// 存储消息
//	l.StoreMessage(c, req)
//
//	//获取群组成员
//	clean := strings.Trim(msg.UserId, `"`)
//	Id, _ := strconv.ParseInt(clean, 10, 64)
//	members, _ := repository.NewContactRequest(global.DB).GetGroupUsers(Id)
//
//	// 广播给在线成员
//	for _, memberID := range members {
//		if memberID.ID == Id {
//			continue // 不发送给自己
//		}
//
//		l.clientsMu.RLock()
//		client, ok := l.clients[strconv.FormatInt(memberID.ID, 10)]
//		l.clientsMu.RUnlock()
//
//		if ok {
//			if err := client.conn.WriteJSON(msg); err != nil {
//				log.Printf("群消息推送失败(%d): %v", memberID.ID, err)
//			}
//		}
//	}
//}
//
//// 发送信息
//func (l *ChatLogic) HandleSendMessage(c *gin.Context, req types.SendMessageRequest) {
//	msg := model.Message{
//		Model:      gorm.Model{},
//		UserId:     req.UserId,
//		UserName:   req.UserName,
//		TargetId:   req.TargetId,
//		TargetName: req.TargetName,
//		Type:       req.Type,
//		Media:      req.Media,
//		Content:    req.Content,
//		CreateTime: uint64(time.Now().Unix()),
//		ReadTime:   req.ReadTime,
//		Pic:        req.Pic,
//		Url:        req.Url,
//		Desc:       req.Desc,
//		Amount:     req.Amount,
//	}
//
//	if err := json.NewDecoder(c.Request.Body).Decode(&msg); err != nil {
//		c.Writer.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	// 存储消息
//	l.StoreMessage(c, req)
//
//	// 尝试实时推送
//	l.SendUserMessage(c, req)
//
//	c.Writer.WriteHeader(http.StatusOK)
//}

// 注册客户端
func (l *ChatLogic) RegisterClient(client *Client) {
	l.clientsMu.Lock()
	defer l.clientsMu.Unlock()
	l.clients[client.userID] = client
}

// 注销客户端
func (l *ChatLogic) UnregisterClient(userID string) {
	l.clientsMu.Lock()
	defer l.clientsMu.Unlock()
	delete(l.clients, userID)
}

// redis缓存
func (l *ChatLogic) StoreMessageConn(c *gin.Context, msg model.Message) {
	// 存储到Redis
	msgJSON, _ := json.Marshal(msg)
	l.redisClient.RPush(c.Request.Context(), "messages:"+msg.TargetId, msgJSON)
}

// // redis缓存
//
//	func (l *ChatLogic) StoreMessage(c *gin.Context, req types.SendMessageRequest) {
//		msg := model.Message{
//			Model:      gorm.Model{},
//			UserId:     req.UserId,
//			UserName:   req.UserName,
//			TargetId:   req.TargetId,
//			TargetName: req.TargetName,
//			Type:       req.Type,
//			Media:      req.Media,
//			Content:    req.Content,
//			CreateTime: uint64(time.Now().Unix()),
//			ReadTime:   req.ReadTime,
//			Pic:        req.Pic,
//			Url:        req.Url,
//			Desc:       req.Desc,
//			Amount:     req.Amount,
//		}
//		// 存储到Redis
//		msgJSON, _ := json.Marshal(msg)
//		l.redisClient.RPush(c.Request.Context(), "messages:"+msg.TargetId, msgJSON)
//	}
func (l *ChatLogic) DeliverOfflineMessages(ctx context.Context, UserID string, Conn *websocket.Conn) {
	key := "messages:" + UserID
	list, err := l.redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil || len(list) == 0 {
		return
	}
	for _, msgJSON := range list {
		var msg model.Message
		if err = json.Unmarshal([]byte(msgJSON), &msg); err == nil {
			_ = Conn.WriteJSON(msg)
		}
	}
	_, _ = l.redisClient.Del(ctx, key).Result()
}

//func (l *ChatLogic) DeliverOfflineMessages(ctx context.Context, req types.DeliverOfflineMessagesRequest) {
//	key := "messages:" + req.UserID
//	list, err := l.redisClient.LRange(ctx, key, 0, -1).Result()
//	if err != nil || len(list) == 0 {
//		return
//	}
//	for _, msgJSON := range list {
//		var msg model.Message
//		if err = json.Unmarshal([]byte(msgJSON), &msg); err == nil {
//			_ = req.Conn.WriteJSON(msg)
//		}
//	}
//	_, _ = l.redisClient.Del(ctx, key).Result()
//}

/*
// ============================ 3. 离线消息推送 ============================
// 当用户上线时，把 Redis 中属于他的离线消息一次性推给他。
func (l *ChatLogic) deliverOfflineMessages(userID string, conn *websocket.Conn) {
	ctx := context.Background()
	key := "messages:" + userID
	list, err := l.redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil || len(list) == 0 {
		return
	}
	for _, msgJSON := range list {
		var msg model.Message
		if err := json.Unmarshal([]byte(msgJSON), &msg); err != nil {
			continue
		}
		_ = conn.WriteJSON(msg)
	}
	// 清空离线队列
	_, _ = l.redisClient.Del(ctx, key).Result()
}

// 在 Connection() 里注册完客户端后调用
// l.registerClient(client)
// go l.deliverOfflineMessages(userID, conn)

// ============================ 4. ACK 机制 ============================
// 给客户端一条「消息已送达」回执，客户端收到后可做 UI 标记。
func (l *ChatLogic) sendACK(conn *websocket.Conn, msgID string) {
	ack := model.Message{
		Type: 4,
		Msg:  msgID,
	}
	_ = conn.WriteJSON(ack)
}

// 在 SendUserMessage / SendGroupMessage 成功推送后调用：
// l.sendACK(TargetId.conn, msg.UUID)
// ============================ 6. Prometheus 指标 ============================
/*
import "github.com/prometheus/client_golang/prometheus"
import "github.com/prometheus/client_golang/prometheus/promhttp"

var (
	conns = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "chat_active_connections",
		Help: "当前活跃 WebSocket 连接数",
	})
	msgCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chat_messages_total",
			Help: "已处理的消息总数",
		},
		[]string{"type"},
	)
)

func init() {
	prometheus.MustRegister(conns, msgCounter)
}

// 在 registerClient / unregisterClient 里分别 conns.Inc()/Dec()
// 在 SendUserMessage / SendGroupMessage 里 msgCounter.WithLabelValues("1"/"2").Inc()

// 暴露 /metrics
// http.Handle("/metrics", promhttp.Handler())
*/

// ============================ 7. 优雅重启 ============================
/*
import "context"
import "os/signal"
import "syscall"

func gracefulStart(addr string) {
	srv := &http.Server{Addr: addr}
	go func() {
		log.Println("服务启动，监听", addr)
		log.Fatal(srv.ListenAndServe())
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("服务关闭中…")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("服务已退出")
}
*/

// ============================ 8. 完整 main 示例 ============================
/*
func main() {
	chatService := NewChatLogic()

	http.HandleFunc("/ws", chatService.Connection)
	http.HandleFunc("/api/send", chatService.HandleSendMessage)
	http.Handle("/metrics", promhttp.Handler())

	gracefulStart(":8080")
}
*/
