package logic

import (
	"context"
	"eino_gozero/apps/AI/rpc/types/ai"

	"eino_gozero/apps/AI/api/internal/svc"
	"eino_gozero/apps/AI/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiVideoChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiVideoChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiVideoChatLogic {
	return &AiVideoChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiVideoChatLogic) AiVideoChat(req *types.AIVideoChatRequest) (resp *types.AIVideoChatResponse, err error) {
	response, err := l.svcCtx.AIRpc.AIVideoChat(l.ctx, &ai.AIVideoChatRequest{
		Ask: req.Ask,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.AIVideoChatResponse{
		Answer: response.Answer,
	}
	return resp, nil
}
