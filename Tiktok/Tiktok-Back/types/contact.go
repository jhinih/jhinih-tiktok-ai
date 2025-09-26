package types

type AddFriendRequest struct {
	UserID     string `json:"user_id"`
	TargetName string `json:"user_name"`
}

type AddFriendResponse struct {
}

//	type SearchFriendRequest struct {
//		UserId string `json:"userId"`
//	}
//
//	type SearchFriendResponse struct {
//		Users []User `json:"users"`
//	}

type GetFriendListRequest struct {
	UserID   string `json:"user_id"`
	Page     string `json:"page"`
	PageSize string `json:"page_size"`
	OrderBy  string `json:"order_by"` // random:随机, latest:最新, popular:最热
}

type GetFriendListResponse struct {
	Users []User `json:"users"`
}

//	type GetFriendsListOnlineRequest struct {
//		UserId string `json:"userId"`
//	}
//
//	type GetFriendsListOnlineResponse struct {
//		Users []User `json:"users"`
//	}
type GetUserListOnlineRequest struct {
}

type GetUserListOnlineResponse struct {
	Users []User `json:"users"`
}
type SearchUserByGroupIdRequest struct {
	CommunityID string `json:"community_id"`
}

type SearchUserByGroupIdResponse struct {
	UserIDs []string `json:"user_ids"`
}

type GetGroupUsersRequest struct {
	CommunityID string `json:"community_id"`
}

type GetGroupUsersResponse struct {
	Users []User `json:"users"`
}

type GetGroupListRequest struct {
	UserID   string `json:"user_id"`
	Page     string `json:"page"`
	PageSize string `json:"page_size"`
	OrderBy  string `json:"order_by"` // random:随机, latest:最新, popular:最热
}
type Community struct {
	ID          string `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	CreatedTime string `gorm:"column:created_time;type:bigint"`
	UpdatedTime string `gorm:"column:updated_time;type:bigint"`
	Name        string `json:"name" gorm:"column:name;type:varchar(255);size:255;unique;not null;"`
	OwnerID     string `json:"owner_id" gorm:"column:owner_id;type:bigint;type:bigint"`
	OwnerName   string `json:"owner_name" gorm:"column:owner_name;type:varchar(255);size:255;"`
	Img         string `json:"img" gorm:"column:img;type:varchar(255);size:255;"`
	Desc        string `json:"desc" gorm:"column:desc;type:varchar(255);size:255;"`
}
type GetGroupListResponse struct {
	Groups []Community `json:"groups"`
}
