package handler

import (
	"eino_gozero/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"eino_gozero/apps/login/api/internal/logic"
	"eino_gozero/apps/login/api/internal/svc"
	"eino_gozero/apps/login/api/internal/types"
)

func RefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshTokenRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRefreshTokenLogic(r.Context(), svcCtx)
		resp, err := l.RefreshToken(&req)
		response.Response(r, w, resp, err)
	}
}
