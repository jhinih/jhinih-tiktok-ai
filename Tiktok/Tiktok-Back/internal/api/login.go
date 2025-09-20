package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/response"
	"Tiktok/types"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	//c.Redirect(http.StatusMovedPermanently, "https://21f3a408.r17.cpolar.top")
}

// SendCode 发送验证码
func SendCode(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.SendCodeRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "发送验证码请求: %v", req)
	resp, err := logic.NewLoginLogic().SendCode(ctx, req)
	response.Response(c, resp, err)
}

// Register 注册
func Register(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.RegisterRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "注册请求: %v", req)
	resp, err := logic.NewLoginLogic().Register(ctx, req)
	response.Response(c, resp, err)
}

// Login 登录
func Login(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.LoginRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "登录请求: %v", req)
	resp, err := logic.NewLoginLogic().Login(ctx, req)
	response.Response(c, resp, err)
}

// RefreshToken 刷新token
func RefreshToken(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.RefreshTokenRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "刷新token请求: %v", req)
	resp, err := logic.NewLoginLogic().RefreshToken(ctx, req)
	response.Response(c, resp, err)
}

// GetToken 刷新token
func GetToken(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetTokenRequest](c)
	if err != nil {
		return
	}
	req.UserName = "ai"
	req.UserID = "1"
	req.Role = "1"
	zlog.CtxInfof(ctx, "刷新token请求: %v", req)
	resp, err := logic.NewLoginLogic().GetToken(ctx, req)
	response.Response(c, resp, err)
}
