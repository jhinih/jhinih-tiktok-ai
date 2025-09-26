package types

type WebsocketRequest struct {
	Token string `form:"token"`
}
type GetHistoryRequest struct {
	Before int64 `json:"before"`
	Count  int64 `json:"count"`
}
