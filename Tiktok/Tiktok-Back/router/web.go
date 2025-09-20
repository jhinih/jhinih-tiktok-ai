package router

import (
	"Tiktok/configs"
	"Tiktok/global"
	"Tiktok/internal/api"
	"Tiktok/log/zlog"
	"Tiktok/manager"
	"Tiktok/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

// RunServer 启动服务器 路由层
func RunServer() {
	r, err := listen()
	if err != nil {
		zlog.Errorf("Listen error: %v", err)
		panic(err.Error())
	}
	r.Run(fmt.Sprintf("%s:%d", configs.Conf.App.Host, configs.Conf.App.Port)) // 启动 Gin 服务器
}

// 自定义ResponseWriter类型
type responseWriter struct {
	gin.ResponseWriter
	headers map[string]string
}

func (w *responseWriter) WriteHeader(code int) {
	for k, v := range w.headers {
		w.ResponseWriter.Header().Set(k, v)
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	for k, v := range w.headers {
		w.ResponseWriter.Header().Set(k, v)
	}
	return w.ResponseWriter.Write(data)
}

// listen 配置 Gin 服务器
func listen() (*gin.Engine, error) {
	r := gin.New()

	// 添加Recovery中间件
	r.Use(gin.Recovery())

	// 注册其他全局中间件
	manager.RequestGlobalMiddleware(r)

	// 静态文件路由（支持视频播放）
	r.Static("/uploads", "uploads")

	// 创建 RouteManager 实例
	routeManager := manager.NewRouteManager(r)
	// 注册各业务路由组的具体路由
	registerRoutes(routeManager)
	// 启动 WebSocket 服务
	startWebsocket(r)
	return r, nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域（生产环境需严格限制）
	},
}

func startWebsocket(r *gin.Engine) {
	manager.WebsocketManager = manager.NewClientManager()
	go manager.WebsocketManager.Start()

	// 注册 WebSocket 路由
	r.GET("/ws", api.WebsocketAPI)
}

// registerRoutes 注册各业务路由的具体处理函数
func registerRoutes(routeManager *manager.RouteManager) {

	// 注册登录相关路由组
	routeManager.RegisterLoginRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/send-code", middleware.Limiter(rate.Every(time.Minute)*4, 4), api.SendCode)
		rg.POST("/register", middleware.Limiter(rate.Every(time.Minute)*4, 4), api.Register)
		rg.POST("/login", middleware.Limiter(rate.Every(time.Minute)*4, 4), api.Login)
		rg.POST("/refresh-token", middleware.Limiter(rate.Every(time.Second)*4, 8), api.RefreshToken)
		rg.POST("/get-token", middleware.Limiter(rate.Every(time.Second)*40, 80), api.GetToken)
	})

	//注册用户相关路由组
	routeManager.RegisterUserRoutes(func(rg *gin.RouterGroup) {
		rg.GET("/get-user-info", middleware.Limiter(rate.Every(time.Second)*20, 40), middleware.Authentication(global.ROLE_GUEST), api.GetUserInfo)
		rg.GET("/get-my-info", middleware.Limiter(rate.Every(time.Second)*5, 10), middleware.Authentication(global.ROLE_GUEST), api.GetMyUserInfo)
		//获取和修改用户资料
		rg.GET("/get-profile", middleware.Limiter(rate.Every(time.Second)*10, 20), middleware.Authentication(global.ROLE_GUEST), api.GetProfile)
		rg.POST("/set-profile", middleware.Limiter(rate.Every(time.Second)*4, 8), middleware.Authentication(global.ROLE_GUEST), api.SetProfile)
		rg.POST("/set-role", middleware.Limiter(rate.Every(time.Second)*4, 8), middleware.Authentication(global.ROLE_GUEST), api.SetRole)

	})
	// 注册文件上传相关路由组
	routeManager.RegisterFileRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/upload", middleware.Limiter(rate.Every(time.Minute)*300, 500), middleware.Authentication(global.ROLE_USER), api.UploadFile)
		rg.POST("/mq_upload", middleware.Limiter(rate.Every(time.Minute)*300, 500), middleware.Authentication(global.ROLE_USER), api.MQUploadFile)
		rg.POST("/mq_upload_result", middleware.Limiter(rate.Every(time.Minute)*3000, 5000), middleware.Authentication(global.ROLE_USER), api.MQUploadResult)
		rg.POST("/upload_sse", middleware.Limiter(rate.Every(time.Minute)*300, 500), middleware.Authentication(global.ROLE_USER), api.UploadSSE)

	})

	// 注册视频相关路由组
	routeManager.RegisterVideosRoutes(func(rg *gin.RouterGroup) {
		// 添加响应头验证
		rg.GET("", middleware.Authentication(global.ROLE_GUEST), api.GetVideos)
		rg.GET("/get-video-likes", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_GUEST), api.GetVideoLikes)
		rg.GET("/get-comment-likes", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_GUEST), api.GetCommentLikes)
		rg.GET("/get-comments", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_GUEST), api.GetComments)
		rg.GET("/get-comment-all", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_GUEST), api.GetCommentAll)
		rg.GET("/get-comment-member", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_GUEST), api.GetCommentsMember)

		rg.POST("/create-video", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.CreateVideo)

		rg.POST("/like-video", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.LikeVideo)
		rg.POST("/like-comments", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.LikeComment)

		rg.POST("/comment-video", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.CommentVideo)
		rg.POST("/comment-comments", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.CommentComment)

	})
	//注册关系相关路由组
	routeManager.RegisterContactRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/add-friend", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.AddFriend)
		rg.POST("/get-friend-list", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.GetFriendList)
		rg.GET("/get-user-list-online", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.GetUserListOnline)
		rg.GET("/get-group-users", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.GetGroupUsers)
		rg.GET("/get-group-list", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.GetGroupList)

		rg.POST("/create-community", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.CreateCommunity)
		rg.POST("/join-community", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.JoinCommunity)
		rg.POST("/load-community", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.LoadCommunity)
		//rg.GET("/get-community-list", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.GetCommunityList)

	})
	//注册AI相关路由组
	routeManager.RegisterAIRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/common_ai", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.CommonAI)
		rg.POST("/video_ai", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.VideoAI)
		rg.POST("/send_code_ai", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.SendCodeAI)
		rg.POST("/get_user_info_ai", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.GetUserInfoAI)
		rg.POST("/ai", middleware.Limiter(rate.Every(time.Minute)*20, 40), middleware.Authentication(global.ROLE_USER), api.AllAI)
	})
}
