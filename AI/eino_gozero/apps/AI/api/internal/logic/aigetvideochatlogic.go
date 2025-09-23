package logic

import (
	"context"
	"eino_gozero/apps/AI/rpc/types/ai"

	"eino_gozero/apps/AI/api/internal/svc"
	"eino_gozero/apps/AI/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AIGetVideoChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAIGetVideoChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AIGetVideoChatLogic {
	return &AIGetVideoChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AIGetVideoChatLogic) AIGetVideoChat(req *types.AIRequest) (resp *types.AIResponse, err error) {
	response, err := l.svcCtx.AIRpc.AIGetVideo(l.ctx, &ai.AIRequest{
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
