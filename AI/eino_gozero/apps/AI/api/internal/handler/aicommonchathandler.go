package handler

import (
	"eino_gozero/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"eino_gozero/apps/AI/api/internal/logic"
	"eino_gozero/apps/AI/api/internal/svc"
	"eino_gozero/apps/AI/api/internal/types"
)

func aiCommonChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AICommonChatRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAiCommonChatLogic(r.Context(), svcCtx)
		resp, err := l.AiCommonChat(&req)
		response.Response(r, w, resp, err)
	}
}
