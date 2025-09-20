package logic

import (
	"context"

	"eino_gozero/apps/login/rpc/internal/svc"
	"eino_gozero/apps/login/rpc/types/login"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *login.LoginRequest) (*login.LoginResponse, error) {

	return &login.LoginResponse{}, nil
}
