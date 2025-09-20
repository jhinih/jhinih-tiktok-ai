package logic

import (
	"context"
	"eino_gozero/apps/AI/rpc/types/ai"

	"eino_gozero/apps/AI/api/internal/svc"
	"eino_gozero/apps/AI/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AIGetUserInfoChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAIGetUserInfoChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AIGetUserInfoChatLogic {
	return &AIGetUserInfoChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AIGetUserInfoChatLogic) AIGetUserInfoChat(req *types.AIRequest) (resp *types.AIResponse, err error) {
	response, err := l.svcCtx.AIRpc.AIGetUserInfo(l.ctx, &ai.AIRequest{
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
