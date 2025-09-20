package logic

import (
	"context"

	"eino_gozero/apps/login/rpc/internal/svc"
	"eino_gozero/apps/login/rpc/types/login"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCodeLogic {
	return &SendCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendCodeLogic) SendCode(in *login.SendCodeRequest) (*login.SendCodeResponse, error) {
	// todo: add your logic here and delete this line

	return &login.SendCodeResponse{}, nil
}
