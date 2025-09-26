package types

type AIRequest struct {
	Ask string `json:"ask" binding:"required"`
}
type AIResponse struct {
	Anser string `json:"anser" binding:"-"`
}
