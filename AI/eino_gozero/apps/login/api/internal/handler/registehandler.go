package handler

import (
	"eino_gozero/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"eino_gozero/apps/login/api/internal/logic"
	"eino_gozero/apps/login/api/internal/svc"
	"eino_gozero/apps/login/api/internal/types"
)

func RegisteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRegisteLogic(r.Context(), svcCtx)
		resp, err := l.Registe(&req)
		response.Response(r, w, resp, err)
	}
}
