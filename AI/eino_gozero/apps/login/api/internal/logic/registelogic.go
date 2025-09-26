package logic

import (
	"context"

	"eino_gozero/apps/login/api/internal/svc"
	"eino_gozero/apps/login/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisteLogic {
	return &RegisteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisteLogic) Registe(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
