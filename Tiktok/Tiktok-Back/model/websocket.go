package model

type ConnMeta struct {
	UserID   string
	UserName string // 需要就加
	Role     int64
	ExpireAt int64  // UTC 秒，方便后台定时刷过期
	Token    string // 如果想后台续签可留着
}
type InMessage struct {
	Type    string `json:"type"`    //chat/history
	Content string `json:"content"` //{"all/user/ai/common_ai","to","content"}
}
type ChatMessage struct {
	Content string `json:"content"` // 发送内容
	ToType  string `json:"to_type"` // 群发、私聊
	To      int64  `json:"to"`      // 发送对象 id
}
type OutMessage struct {
	Type      string `json:"type"`
	ID        int64  `json:"id,string"`
	Content   string `json:"content"`
	UserID    string `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
}
