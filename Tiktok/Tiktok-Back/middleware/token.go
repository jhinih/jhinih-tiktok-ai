package middleware

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/response"
	"Tiktok/utils/jwtUtils"
	"github.com/gin-gonic/gin"
	"strings"
)

//	func Authentication(role int) gin.HandlerFunc {
//		return func(c *gin.Context) {
//			ctx := zlog.GetCtxFromGin(c)
//			authorization := c.GetHeader("Authorization")
//			if authorization == "" {
//				zlog.CtxErrorf(ctx, "authorization为空")
//				response.NewResponse(c).Error(response.TOKEN_IS_BLANK)
//				c.Abort()
//				return
//			}
//			// 以空格分割，取出token
//			list := strings.Split(authorization, " ")
//			if len(list) != 2 {
//				zlog.CtxErrorf(ctx, "token格式错误")
//				response.NewResponse(c).Error(response.TOKEN_FORMAT_ERROR)
//				c.Abort()
//				return
//			}
//			token := list[1]
//			//解析token是否有效，并取出上一次的值
//			data, err := jwtUtils.IdentifyToken(token)
//			if err != nil {
//				zlog.CtxErrorf(ctx, "token验证失败:%v", err)
//				response.NewResponse(c).Error(response.TOKEN_IS_EXPIRED)
//				//对应token无效，直接让他返回
//				c.Abort()
//				return
//			}
//			//判断其是否为atoken
//			if data.Email != global.AUTH_ENUMS_ATOKEN {
//				zlog.CtxErrorf(ctx, "token类型错误")
//				response.NewResponse(c).Error(response.TOKEN_TYPE_ERROR)
//				c.Abort()
//				return
//			}
//			// 判断权限是否足够
//			if data.Role < role {
//				zlog.CtxErrorf(ctx, "权限不足")
//				response.NewResponse(c).Error(response.PERMISSION_DENIED)
//				c.Abort()
//				return
//			}
//			//将token内部数据传下去,在logic.token内有对应方法获取userid
//			c.Set(global.TOKEN_USER_ID, data.Userid)
//			c.Set(global.TOKEN_ROLE, data.Role)
//			c.Next()
//		}
//	}
//
// middleware/token.go
func Authentication(role int) gin.HandlerFunc {
	return func(c *gin.Context) {
		zlog.GetCtxFromGin(c)

		// 检查是否是 WebSocket 请求
		if c.Request.Header.Get("Upgrade") == "websocket" {
			token := c.Query("ticket")
			data, _ := jwtUtils.IdentifyToken(token)
			c.Set(global.TOKEN_USER_ID, data.Userid)
			c.Set(global.TOKEN_ROLE, data.Role)
			c.Set(global.TOKEN_USER_NAME, data.Username)
			c.Next()
			return
		}
		token := ""
		auth := c.GetHeader("Authorization")
		if len(auth) > 7 && strings.HasPrefix(auth, "Bearer ") {
			token = auth[7:]
		}
		if token == "" {
			response.NewResponse(c).Error(response.TOKEN_IS_BLANK)
			c.Abort()
			return
		}

		data, err := jwtUtils.IdentifyToken(token)
		if err != nil {
			response.NewResponse(c).Error(response.TOKEN_IS_EXPIRED)
			c.Abort()
			return
		}

		if data.Role < role {
			response.NewResponse(c).Error(response.PERMISSION_DENIED)
			c.Abort()
			return
		}

		c.Set(global.TOKEN_USER_ID, data.Userid)
		c.Set(global.TOKEN_ROLE, data.Role)
		c.Set(global.TOKEN_USER_NAME, data.Username)
		c.Next()
	}
}
