package types

import (
	"Tiktok/response"
	"errors"
	"github.com/gin-gonic/gin"
)

// fixme 注意这个 bindReq 并不适合符合restful 设计规范的参数绑定,只适合单类型参数绑定

// BindReq Uri绑定 /for/example/:id
func BindUri[T any](c *gin.Context) (req T, err error) {
	if err = c.ShouldBindUri(&req); err != nil {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
	}
	return req, err
}

// BindJson Json绑定 {code:1}
func BindJson[T any](c *gin.Context) (req T, err error) {
	if err = c.ShouldBindJSON(&req); err != nil {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
	}
	return req, err
}

// BindQuery Query绑定 /for/example?code=1
func BindQuery[T any](c *gin.Context) (req T, err error) {
	if err = c.ShouldBindQuery(&req); err != nil {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
	}
	return req, err
}

// BindReq 绑定请求 目前来讲更加通用
func BindRequest[T any](c *gin.Context) (req T, err error) {
	switch c.Request.Method {
	case "GET":
		return BindQuery[T](c)
	case "POST":
		return BindJson[T](c)
	case "PUT":
		return BindJson[T](c)
	case "DELETE":
		return BindUri[T](c)
		// 没有使用 PATCH 方法
	default:
		// 不可能到达这里
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return req, errors.New("method not support")
	}

}
