package logic

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/model"
	"Tiktok/repository"
	"Tiktok/response"
	"Tiktok/types"
	"Tiktok/utils"
	"context"
	//"encoding/json"
	"errors"
	//"fmt"
	"gorm.io/gorm"
	//"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	CODEFORCES_API_URL             = "https://codeforces.com/api/user.info?handles=%s&checkHistoricHandles=false"
	REDIS_CODEFORCES_IS_UPDATE_KEY = "codeforces_is_update_%s"
)

type UserLogic struct {
}

func NewUserLogic() *UserLogic {
	return &UserLogic{}
}

// GetUserInfo 获取用户信息
func (l *UserLogic) GetUserInfo(ctx context.Context, req types.GetUserInfoRequest) (resp types.GetUserInfoResponse, err error) {
	defer utils.RecordTime(time.Now())()
	zlog.CtxInfof(ctx, "获取用户信息 %s", req.ID)
	// ID 转化为 int64
	userID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%s 转换 int64 错误: %v", req.ID, err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	// 获取用户信息
	var user model.User
	user, err = repository.NewUserRequest(global.DB).GetUserProfileByID(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zlog.CtxErrorf(ctx, "用户并不存在!: %v", err)
		return resp, response.ErrResponse(err, response.USER_NOT_EXIST)
	}
	// 填入参数
	User := types.User{
		ID:            strconv.FormatInt(user.ID, 10),
		CreatedTime:   strconv.FormatInt(user.CreatedTime, 10),
		UpdatedTime:   strconv.FormatInt(user.UpdatedTime, 10),
		Username:      user.Username,
		Password:      user.Password,
		Email:         user.Email,
		Avatar:        user.Avatar,
		Role:          strconv.Itoa(user.Role),
		Phone:         user.Phone,
		ClientIp:      user.ClientIp,
		ClientPort:    user.ClientPort,
		LoginTime:     user.LoginTime,
		HeartbeatTime: user.HeartbeatTime,
		LoginOutTime:  user.LoginOutTime,
		IsLogout:      user.IsLogout,
		DeviceInfo:    user.DeviceInfo,
		Bio:           user.Bio,
	}
	resp.User = User

	zlog.CtxInfof(ctx, "成功获取用户信息: ID=%d, Username=%s", userID, user.Username)
	return resp, nil
}

// GetUserProfile 获取用户信息
func (l *UserLogic) GetUserProfile(ctx context.Context, req types.GetUserProfileRequest) (resp types.GetUserProfileResponse, err error) {
	defer utils.RecordTime(time.Now())()
	zlog.CtxInfof(ctx, "开始获取用户信息，请求ID: %s", req.ID)

	// ID 转化为 int64
	userID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "ID转换错误: %s -> %v", req.ID, err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}

	// 获取用户信息
	var user model.User
	user, err = repository.NewUserRequest(global.DB).GetUserProfileByID(userID)
	if err != nil {
		zlog.CtxErrorf(ctx, "数据库查询错误: %v", err)
		return resp, response.ErrResponse(err, response.DATABASE_ERROR)
	}
	// 填入参数
	User := types.User{
		ID:            strconv.FormatInt(user.ID, 10),
		CreatedTime:   strconv.FormatInt(user.CreatedTime, 10),
		UpdatedTime:   strconv.FormatInt(user.UpdatedTime, 10),
		Username:      user.Username,
		Password:      user.Password,
		Email:         user.Email,
		Avatar:        user.Avatar,
		Role:          strconv.Itoa(user.Role),
		Phone:         user.Phone,
		ClientIp:      user.ClientIp,
		ClientPort:    user.ClientPort,
		LoginTime:     user.LoginTime,
		HeartbeatTime: user.HeartbeatTime,
		LoginOutTime:  user.LoginOutTime,
		IsLogout:      user.IsLogout,
		DeviceInfo:    user.DeviceInfo,
		Bio:           user.Bio,
	}
	resp.User = User

	zlog.CtxInfof(ctx, "成功获取用户信息: ID=%d, Username=%s", userID, user.Username)
	return resp, nil
}

func (l *UserLogic) SetUserProfile(ctx context.Context, req types.SetUserProfileRequest) (resp types.SetUserProfileResponse, err error) {
	defer utils.RecordTime(time.Now())()
	// ID 转化为 int64
	userID, err := strconv.ParseInt(req.User.ID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", req.User.ID, err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	operatorID, err := strconv.ParseInt(req.OperatorID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", req.OperatorID, err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	// 检查权限
	if operatorID != userID {
		zlog.CtxErrorf(ctx, "非法操作: %v", req.OperatorID)
		return resp, response.ErrResponse(err, response.PERMISSION_DENIED)

	}
	// 检验数据
	// 1.用户名去除所有空格，且不能为空，且长度不超过 30
	req.User.Username = strings.ReplaceAll(req.User.Username, " ", "")
	if req.User.Username == "" {
		zlog.CtxErrorf(ctx, "用户名不能为空")
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	} else if len(req.User.Username) > 30 {
		zlog.CtxErrorf(ctx, "用户名长度不能超过 30")
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}

	if len(req.User.Avatar) > 255 {
		zlog.CtxErrorf(ctx, "头像 URL 长度不能超过 255")
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	// 拿出先前的用户信息
	user, err := repository.NewUserRequest(global.DB).GetUserProfileByID(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zlog.CtxErrorf(ctx, "用户并不存在!: %v", err)
		return resp, response.ErrResponse(err, response.USER_NOT_EXIST)
	} else if err != nil {
		zlog.CtxErrorf(ctx, "获取用户信息失败: %v", err)
		return resp, response.ErrResponse(err, response.DATABASE_ERROR)
	}
	// 更新用户信息
	user = model.User{
		ID: userID,
		TimeModel: model.TimeModel{
			CreatedTime: user.CreatedTime,
			UpdatedTime: user.UpdatedTime,
		},
		Username:      user.Username,
		Password:      user.Password,
		Email:         user.Email,
		Avatar:        user.Avatar,
		Role:          user.Role,
		Phone:         user.Phone,
		ClientIp:      user.ClientIp,
		ClientPort:    user.ClientPort,
		LoginTime:     user.LoginTime,
		HeartbeatTime: user.HeartbeatTime,
		LoginOutTime:  user.LoginOutTime,
		IsLogout:      user.IsLogout,
		DeviceInfo:    user.DeviceInfo,
		Bio:           user.Bio,
	}
	// 如果用户的身份是游客，那么这次提交将升级为普通用户
	if user.Role == 0 {
		user.Role = 1
	}
	err = repository.NewUserRequest(global.DB).UpdateUserProfile(user)
	if err != nil {
		zlog.CtxErrorf(ctx, "更新用户信息失败: %v", err)
		return resp, response.ErrResponse(err, response.DATABASE_ERROR)
	}
	return resp, nil
}

// SetUserRole 设置用户权限
func (l *UserLogic) SetUserRole(ctx context.Context, req types.SetUserRoleRequest) (resp types.SetUserRoleResp, err error) {
	defer utils.RecordTime(time.Now())()
	// ID 转化为 int64
	userID, err := strconv.ParseInt(req.ID, 10, 64)
	role, _ := strconv.Atoi(req.Role)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", req.ID, err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	//// 检查权限(必须是超级管理员)
	//if req.OperatorRole < 4 {
	//	zlog.CtxErrorf(ctx, "非法操作")
	//	return resp, response.ErrResponse(err, response.PERMISSION_DENIED)
	//}
	// 修改权限
	err = repository.NewUserRequest(global.DB).SetUserRole(userID, role)
	if err != nil {
		zlog.CtxErrorf(ctx, "修改权限失败: %v", err)
		return resp, response.ErrResponse(err, response.DATABASE_ERROR)
	}
	zlog.CtxInfof(ctx, "修改权限成功")
	user, err := repository.NewUserRequest(global.DB).GetUserProfileByID(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zlog.CtxErrorf(ctx, "用户并不存在!: %v", err)
		return resp, response.ErrResponse(err, response.USER_NOT_EXIST)
	} else if err != nil {
		zlog.CtxErrorf(ctx, "获取用户信息失败: %v", err)
		return resp, response.ErrResponse(err, response.DATABASE_ERROR)
	}
	// 更新用户信息
	err = repository.NewUserRequest(global.DB).UpdateUserProfile(user)
	if err != nil {
		zlog.CtxErrorf(ctx, "更新用户信息失败: %v", err)
		return resp, response.ErrResponse(err, response.DATABASE_ERROR)
	}

	return resp, nil
}
