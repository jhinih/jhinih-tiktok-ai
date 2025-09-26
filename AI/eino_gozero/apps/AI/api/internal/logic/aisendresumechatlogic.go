package logic

import (
	"context"
	"eino_gozero/apps/AI/rpc/types/ai"

	"eino_gozero/apps/AI/api/internal/svc"
	"eino_gozero/apps/AI/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AISendResumeChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAISendResumeChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AISendResumeChatLogic {
	return &AISendResumeChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AISendResumeChatLogic) AISendResumeChat(req *types.AIRequest) (resp *types.AIResponse, err error) {
	response, err := l.svcCtx.AIRpc.AISendResume(l.ctx, &ai.AIRequest{
		Ask: req.Ask,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.AIResponse{
		Answer: response.Answer,
	}
	return resp, nil
}
