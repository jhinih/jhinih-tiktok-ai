package logic

import (
	"context"

	"eino_gozero/apps/login/rpc/internal/svc"
	"eino_gozero/apps/login/rpc/types/login"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *login.RegisterRequest) (*login.RegisterResponse, error) {
	// todo: add your logic here and delete this line

	return &login.RegisterResponse{}, nil
}
