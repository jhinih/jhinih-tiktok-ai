package logic

import (
	"context"

	"eino_gozero/apps/login/rpc/internal/svc"
	"eino_gozero/apps/login/rpc/types/login"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshTokenLogic) RefreshToken(in *login.RefreshTokenRequest) (*login.RefreshTokenResponse, error) {
	// todo: add your logic here and delete this line

	return &login.RefreshTokenResponse{}, nil
}
