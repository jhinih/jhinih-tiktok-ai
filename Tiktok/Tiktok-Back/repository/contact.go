package repository

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/model"
	"Tiktok/repository/list"
	"gorm.io/gorm"
	"time"
)

type ContactRequest struct {
	DB *gorm.DB
}

func NewContactRequest(db *gorm.DB) *ContactRequest {
	return &ContactRequest{
		DB: db,
	}
}

// 添加好友
func (r *ContactRequest) AddFriend(ownerUser, targetUser model.User) (err error) {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
		if err != nil {
			tx.Rollback()
		}
	}()
	id := global.SnowflakeNode.Generate().Int64()

	// 创建正向联系人记录
	contact := model.Contact{
		ID: id,
		TimeModel: model.TimeModel{
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
		OwnerID:    ownerUser.ID,
		OwnerName:  ownerUser.Username,
		TargetID:   targetUser.ID,
		TargetName: targetUser.Username,
		Type:       1,
		Desc:       "",
	}

	// 添加详细日志
	zlog.Debugf("Creating forward contact: %+v", contact)
	if err = tx.Create(&contact).Error; err != nil {
		zlog.Errorf("Failed to create forward contact: %v", err)
		return err
	}
	zlog.Debugf("Forward contact created successfully, ID: %d", contact.ID)
	id = global.SnowflakeNode.Generate().Int64()
	// 创建反向联系人记录
	reverseContact := model.Contact{
		ID: id,
		TimeModel: model.TimeModel{
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
		OwnerID:    targetUser.ID,
		OwnerName:  targetUser.Username,
		TargetID:   ownerUser.ID,
		TargetName: ownerUser.Username,
		Type:       1,
		Desc:       "",
	}

	zlog.Debugf("Creating reverse contact: %+v", reverseContact)
	if err = tx.Create(&reverseContact).Error; err != nil {
		zlog.Errorf("Failed to create reverse contact: %v", err)
		return err
	}
	zlog.Debugf("Reverse contact created successfully, ID: %d", reverseContact.ID)

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// 好友列表
func (r *ContactRequest) GetFriendList(userId, page, pageSize int64, orderBy string) ([]model.User, error) {
	//contacts := make([]model.Contact, 0)
	//ObjIDs := make([]int64, 0)
	//err := r.DB.Where("owner_id = ? and type=1", userId).Find(&contacts).Error
	//if err != nil {
	//	return nil, err
	//}
	//for _, v := range contacts {
	//	ObjIDs = append(ObjIDs, v.TargetID)
	//}
	//users := make([]model.User, 0)
	//err = r.DB.Where("id in ?", ObjIDs).Find(&users).Error
	//var Users []model.User
	//Users, _, err = list.Query[model.User](list.Options{
	//	Where: func(db *gorm.DB) *gorm.DB {
	//		return db.Where("")
	//	},
	//	Order: orderBy,
	//	//Preloads: []string{"Author"},
	//	PageInfo: list.PageInfo{Page: page, Limit: pageSize},
	//})
	//return Users, err
	// 只保留这一段即可
	users, _, err := list.Query[model.User](list.Options{
		Where: func(db *gorm.DB) *gorm.DB {
			// 子查询：找出当前用户关注的所有人
			return db.Where("id IN (?)",
				r.DB.Model(&model.Contact{}).
					Select("target_id").
					Where("owner_id = ? AND type = 1", userId),
			)
		},
		Order: orderBy, // latest / popular / random 等
		//Preloads: []string{"Author"}, // 需要时再放开
		PageInfo: list.PageInfo{Page: page, Limit: pageSize},
	})
	return users, err
}

// 是否是好友
func (r *ContactRequest) IsFriend(UserId, TargetID int64) (model.Contact, error) {
	var contact model.Contact
	err := r.DB.Where("owner_id =?  and target_id =? and type=1", UserId, TargetID).Find(&contact).Error
	return contact, err
}

// 创建群聊
func (r *ContactRequest) CreatCommunity(contact model.Contact) error {
	err := r.DB.Create(&contact).Error
	return err
}

// 在线用户列表
func (r *ContactRequest) SearchUsersOnline(users *[]model.User) error {
	if err := r.DB.Where("is_logout = ?", false).Find(users).Error; err != nil {
		return err
	}
	if len(*users) == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 获取群友
func (r *ContactRequest) GetGroupUsers(communityId int64) ([]model.User, error) {
	contacts := make([]model.Contact, 0)
	objIds := make([]int64, 0)
	err := r.DB.Where("target_id = ? and type=2", communityId).Find(&contacts).Error
	if err != nil {
		return nil, err
	}
	for _, v := range contacts {
		objIds = append(objIds, v.TargetID)
	}
	users := make([]model.User, 0)
	err = r.DB.Where("id in ?", objIds).Find(&users).Error
	return users, err
}

// 获取群列表
func (r *ContactRequest) GetGroupList(userId, page, pageSize int64, orderBy string) ([]model.Community, error) {
	//contacts := make([]model.Contact, 0)
	//objIds := make([]int64, 0)
	//err := r.DB.Where("owner_id = ? and type=1", userId).Find(&contacts).Error
	//if err != nil {
	//	return nil, err
	//}
	//for _, v := range contacts {
	//	objIds = append(objIds, v.TargetID)
	//}
	//groups := make([]model.Community, 0)
	//err = r.DB.Where("id in ?", objIds).Find(&groups).Error
	//var Groups []model.Community
	//Groups, _, err = list.Query[model.Community](list.Options{
	//	Where: func(db *gorm.DB) *gorm.DB {
	//		return db.Where("")
	//	},
	//	Order: orderBy,
	//	//Preloads: []string{"Author"},
	//	PageInfo: list.PageInfo{Page: page, Limit: pageSize},
	//})
	//return Groups, err
	// 删掉 contacts/objIds 的代码，直接这样：
	groups, _, err := list.Query[model.Community](list.Options{
		Where: func(db *gorm.DB) *gorm.DB {
			return db.Where("id IN (?)",
				r.DB.Model(&model.Contact{}).
					Select("target_id").
					Where("owner_id = ? AND type = 1", userId),
			)
		},
		Order: orderBy,
		//Preloads: []string{"Author"}, // 需要时再打开
		PageInfo: list.PageInfo{Page: page, Limit: pageSize},
	})
	return groups, err

}

//// 群是否被拥有
//func (r *ContactRequest) SearchUserByGroupId(communityId int64) []model.Contact {
//	contacts := make([]model.Contact, 0)
//	r.DB.Where("target_id = ? and type=2", communityId).Find(&contacts)
//	return contacts
//}
