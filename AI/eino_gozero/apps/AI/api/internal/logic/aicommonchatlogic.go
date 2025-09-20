package logic

import (
	"context"
	"eino_gozero/apps/AI/rpc/types/ai"

	"eino_gozero/apps/AI/api/internal/svc"
	"eino_gozero/apps/AI/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiCommonChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiCommonChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiCommonChatLogic {
	return &AiCommonChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiCommonChatLogic) AiCommonChat(req *types.AICommonChatRequest) (resp *types.AICommonChatResponse, err error) {
	response, err := l.svcCtx.AIRpc.AICommonChat(l.ctx, &ai.AICommonChatRequest{
		Ask: req.Ask,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.AICommonChatResponse{
		Answer: response.Answer,
	}
	return resp, nil
}
