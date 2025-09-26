package logic

import (
	"context"

	"eino_gozero/apps/login/rpc/internal/svc"
	"eino_gozero/apps/login/rpc/types/login"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTokenLogic) GetToken(in *login.GetTokenRequest) (*login.GetTokenResponse, error) {
	// todo: add your logic here and delete this line

	return &login.GetTokenResponse{}, nil
}
