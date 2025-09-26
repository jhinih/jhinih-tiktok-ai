package types

type CreateCommunityRequest struct {
	OwnerID   string `json:"owner_id"`
	OwnerName string `json:"owner_name"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Desc      string `json:"desc"`
}

type CreateCommunityResponse struct {
}

// 加入群聊
type JoinCommunityRequest struct {
	UserID      string `json:"user_id"`
	CommunityID string `json:"community_id"`
}

type JoinCommunityResponse struct {
}
type LoadCommunityRequest struct {
	OwnerID string `json:"owner_id"`
}

type LoadCommunityResponse struct {
	Groups []Community `json:"groups"`
}
