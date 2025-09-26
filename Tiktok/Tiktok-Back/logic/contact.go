package logic

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/model"
	"Tiktok/repository"
	"Tiktok/response"
	"Tiktok/types"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type ContactLogic struct {
}

func NewContactLogic() *ContactLogic {
	return &ContactLogic{}
}

func (l *ContactLogic) AddFriend(ctx context.Context, req types.AddFriendRequest) (resp types.AddFriendResponse, err error) {
	if req.TargetName == "" {
		zlog.CtxErrorf(ctx, "好友用户名不能为空")
		return resp, response.ErrResponse(err, response.PARAM_IS_BLANK)
	}

	// 查询发起用户
	UserID, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "用户ID转换失败: %v, UserID: %s", err, req.UserID)
		return resp, response.ErrResponse(err, response.PARAM_IS_INVALID)
	}
	OwnerUser, err := repository.NewUserRequest(global.DB).FindUserByID(int64(UserID))
	if err != nil {
		zlog.CtxErrorf(ctx, "查询发起用户失败: %v, UserID: %s", err, req.UserID)
		return resp, response.ErrResponse(err, response.USER_NOT_EXIST)
	}

	// 查询目标用户
	TargetUser, err := repository.NewUserRequest(global.DB).FindUserByName(req.TargetName)
	if err != nil {
		zlog.CtxErrorf(ctx, "查询目标用户失败: %v, username: %s", err, req.TargetName)
		return resp, response.ErrResponse(err, response.USER_NOT_EXIST)
	}

	// 检查是否添加自己
	if TargetUser.ID == int64(UserID) {
		zlog.CtxErrorf(ctx, "不能添加自己为好友, userId: %d", UserID)
		return resp, response.ErrResponse(err, response.ME_AND_ME)
	}

	// 检查是否已是好友
	contact, err := repository.NewContactRequest(global.DB).IsFriend(int64(UserID), int64(TargetUser.ID))
	if err != nil {
		zlog.CtxErrorf(ctx, "检查好友关系失败: %v", err)
		return resp, response.ErrResponse(err, response.COMMON_FAIL)
	}
	if contact.ID != 0 {
		zlog.CtxErrorf(ctx, "已是好友, 不能重复添加")
		return resp, response.ErrResponse(err, response.FRIEND_YES_FRIEN)
	}

	// 添加好友
	// 添加详细日志记录
	zlog.CtxDebugf(ctx, "准备添加好友: OwnerID=%d, TargetID=%d", OwnerUser.ID, TargetUser.ID)

	// 添加事务重试机制
	var retryCount = 0
	maxRetries := 3
	var lastErr error

	for retryCount < maxRetries {
		if err := repository.NewContactRequest(global.DB).AddFriend(OwnerUser, TargetUser); err != nil {
			lastErr = err
			zlog.CtxErrorf(ctx, "添加好友失败(尝试 %d/%d): %v", retryCount+1, maxRetries, err)
			retryCount++
			continue
		}
		break
	}

	if retryCount == maxRetries {
		zlog.CtxErrorf(ctx, "添加好友最终失败: %v", lastErr)
		return resp, response.ErrResponse(lastErr, response.COMMON_FAIL)
	}

	zlog.CtxInfof(ctx, "添加好友成功: OwnerID=%d, TargetID=%d", OwnerUser.ID, TargetUser.ID)
	return resp, nil
}

func (l *ContactLogic) GetFriendList(ctx context.Context, req types.GetFriendListRequest) (resp types.GetFriendListResponse, err error) {
	UserID, _ := strconv.ParseInt(req.UserID, 10, 64)
	Page, _ := strconv.ParseInt(req.Page, 10, 64)
	PageSize, _ := strconv.ParseInt(req.Page, 10, 64)
	users, err := repository.NewContactRequest(global.DB).GetFriendList(int64(UserID), Page, PageSize, req.OrderBy)
	if len(users) == 0 {
		zlog.CtxErrorf(ctx, "好友列表获取失败: %v", err)
	}
	Users := make([]types.User, 0)
	for _, user := range users {
		Users = append(Users, types.User{
			ID:            strconv.FormatInt(user.ID, 10),
			CreatedTime:   strconv.FormatInt(user.CreatedTime, 10),
			UpdatedTime:   strconv.FormatInt(user.CreatedTime, 10),
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
		})
	}
	resp.Users = Users
	return resp, nil
}

func (l *ContactLogic) GetUserListOnline(ctx context.Context) (resp types.GetUserListOnlineResponse, err error) {
	var users []model.User
	if err := repository.NewContactRequest(global.DB).SearchUsersOnline(&users); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxErrorf(ctx, "在线列表获取为空: %v", err)
			return resp, nil
		}
		return resp, fmt.Errorf("在线列表获取失败: %w", err)
	}
	Users := make([]types.User, 0)
	for _, user := range users {
		Users = append(Users, types.User{
			ID:            strconv.FormatInt(user.ID, 10),
			CreatedTime:   strconv.FormatInt(user.CreatedTime, 10),
			UpdatedTime:   strconv.FormatInt(user.CreatedTime, 10),
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
		})
	}
	resp.Users = Users
	return resp, nil
}

func (l *ContactLogic) GetGroupUsers(ctx context.Context, req types.GetGroupUsersRequest) (resp types.GetGroupUsersResponse, err error) {
	CommunityID, _ := strconv.ParseInt(req.CommunityID, 10, 64)
	users, err := repository.NewContactRequest(global.DB).GetGroupUsers(int64(CommunityID))
	if len(users) == 0 {
		zlog.CtxErrorf(ctx, "搜索群友失败: %v", err)
	}
	Users := make([]types.User, 0)
	for _, user := range users {
		Users = append(Users, types.User{
			ID:            strconv.FormatInt(user.ID, 10),
			CreatedTime:   strconv.FormatInt(user.CreatedTime, 10),
			UpdatedTime:   strconv.FormatInt(user.CreatedTime, 10),
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
		})
	}
	resp.Users = Users
	return resp, nil
}
func (l *ContactLogic) GetGroupList(ctx context.Context, req types.GetGroupListRequest) (resp types.GetGroupListResponse, err error) {
	UserID, _ := strconv.ParseInt(req.UserID, 10, 64)
	Page, _ := strconv.ParseInt(req.Page, 10, 64)
	PageSize, _ := strconv.ParseInt(req.Page, 10, 64)
	groups, err := repository.NewContactRequest(global.DB).GetGroupList(int64(UserID), Page, PageSize, req.OrderBy)
	if len(groups) == 0 {
		zlog.CtxErrorf(ctx, "获取群聊列表失败: %v", err)
	}
	Groups := make([]types.Community, 0)
	for _, group := range groups {
		Groups = append(Groups, types.Community{
			ID:          strconv.FormatInt(group.ID, 10),
			CreatedTime: strconv.FormatInt(group.CreatedTime, 10),
			UpdatedTime: strconv.FormatInt(group.CreatedTime, 10),
			Name:        group.Name,
			OwnerID:     strconv.FormatInt(group.OwnerID, 10),
			OwnerName:   group.OwnerName,
			Img:         group.Img,
			Desc:        group.Desc,
		})
	}
	resp.Groups = Groups
	return resp, nil
}

//群是否被拥有
//func (l *ContactLogic) SearchUserByGroupId(ctx context.Context, req types.SearchUserByGroupIdRequest) (resp types.SearchUserByGroupIdResponse, err error) {
//	contacts := make([]model.Contact, 0)
//	objIds := make([]int64, 0)
//	clean := strings.Trim(req.CommunityId, `"`)
//	Id, _ := strconv.ParseInt(clean, 10, 64)
//	contacts = repository.NewContactRequest(global.DB).SearchUserByGroupId(int64(Id))
//	for _, v := range contacts {
//		objIds = append(objIds, v.OwnerId)
//	}
//	resp.UserIds = objIds
//	zlog.CtxDebugf(ctx, "查找群友成功: %v", req)
//	return resp, err
//}
