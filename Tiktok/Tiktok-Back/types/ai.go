package types

type AIRequest struct {
	Ask   string `json:"ask" binding:"required"`
	Token string `json:"token" binding:"-"`
}
type AIResponse struct {
	Anser string `json:"anser" binding:"-"`
}
type CozeRequest struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Token   string `json:"token" binding:"-"`
}

type CozeResponse struct {
	Seq       int    `json:"seq"`
	Type      string `json:"type"`
	ID        int64  `json:"id,string"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
