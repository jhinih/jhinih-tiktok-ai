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
	"strconv"
	"time"
)

type CommunityLogic struct {
}

func NewCommunityLogic() *CommunityLogic {
	return &CommunityLogic{}
}

// 新建群
func (l *CommunityLogic) CreateCommunity(ctx context.Context, req types.CreateCommunityRequest) (resp types.CreateCommunityResponse, err error) {
	OwnerID, _ := strconv.ParseInt(req.OwnerID, 10, 64)
	ID := global.SnowflakeNode.Generate().Int64()
	community := model.Community{
		ID:        ID,
		OwnerID:   int64(OwnerID),
		OwnerName: req.OwnerName,
		Name:      req.Name,
		Img:       req.Icon,
		Desc:      req.Desc,
		TimeModel: model.TimeModel{
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
	}
	if len(community.Name) == 0 {
		//"群名称不能为空"
		response.ErrResponse(err, response.COMMUNITY_IS_BLANK)
	}
	comIS := repository.NewCommunityRequest(global.DB).FindCommunityByName(community.Name)
	if comIS.Name == "" {
		if utils.IsNumeric(community.Name) {
			//"开发者不允许你拿数字建群"
			response.ErrResponse(err, response.FACK_FACK_FACK)
		}
	} else if comIS.Name == community.Name {
		//"群聊已存在"
		response.ErrResponse(err, response.EMAIL_NOT_VALID)
	}
	err = repository.NewCommunityRequest(global.DB).CreateCommunity(community)
	if err != nil {
		zlog.CtxErrorf(ctx, "创建群聊失败: %v", err)
		return resp, response.ErrResponse(err, response.DATABASE_ERROR)
	}
	return resp, nil
}

// 加入群聊
func (l *CommunityLogic) JoinCommunity(ctx context.Context, req types.JoinCommunityRequest) (resp types.JoinCommunityResponse, err error) {
	UserID, _ := strconv.ParseInt(req.UserID, 10, 64)
	CommunityID, _ := strconv.ParseInt(req.CommunityID, 10, 64)

	community := model.Community{}
	community = repository.NewCommunityRequest(global.DB).FindCommunityByNameOrId(int64(CommunityID))
	contact := model.Contact{}
	contact = repository.NewCommunityRequest(global.DB).IsInCommunity(UserID, community)
	if contact.TimeModel.CreatedTime != 0 {
		//"已加过此群"
	} else {
		contact.OwnerID = int64(UserID)
		contact.TargetID = int64(community.ID)
		contact.TargetName = community.Name
		contact.TimeModel.CreatedTime = time.Now().Unix()
		contact.TimeModel.UpdatedTime = time.Now().Unix()
		contact.Type = 2
		contact.Desc = "" //后续添加描述
		err = repository.NewContactRequest(global.DB).CreatCommunity(contact)
		//"加群成功"
	}
	return resp, err
}

// 加载群列表
func (l *CommunityLogic) LoadCommunity(ctx context.Context, req types.LoadCommunityRequest) (resp types.LoadCommunityResponse, err error) {
	OwnerID, _ := strconv.ParseInt(req.OwnerID, 10, 64)
	contacts := make([]model.Contact, 0)
	ObjIDs := make([]int64, 0)
	contacts = repository.NewCommunityRequest(global.DB).LoadUserCommunity(int64(OwnerID))
	for _, v := range contacts {
		ObjIDs = append(ObjIDs, v.TargetID)
	}
	data := repository.NewCommunityRequest(global.DB).LoadCommunityUser([]int64(ObjIDs))
	if len(data) != 0 {
		zlog.CtxDebugf(ctx, "成功获取群数量: len=%d", len(data))
		//response.ErrResponse(err, response.COMMUNITY_IS_NILL)
	} else {
		response.ErrResponse(err, response.COMMUNITY_IS_NILL)
	}
	Groups := make([]types.Community, 0)
	for _, group := range data {
		Groups = append(Groups, types.Community{
			ID:          strconv.FormatInt(group.ID, 10),
			CreatedTime: strconv.FormatInt(group.CreatedTime, 10),
			UpdatedTime: strconv.FormatInt(group.UpdatedTime, 10),
			Name:        group.Name,
			OwnerID:     strconv.FormatInt(group.OwnerID, 10),
			OwnerName:   group.OwnerName,
			Img:         group.Img,
			Desc:        group.Desc,
		})
	}
	resp.Groups = Groups
	return resp, err
}
