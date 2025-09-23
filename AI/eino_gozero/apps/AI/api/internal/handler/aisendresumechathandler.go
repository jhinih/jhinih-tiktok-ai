package handler

import (
	"eino_gozero/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"eino_gozero/apps/AI/api/internal/logic"
	"eino_gozero/apps/AI/api/internal/svc"
	"eino_gozero/apps/AI/api/internal/types"
)

func AISendResumeChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AIRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAISendResumeChatLogic(r.Context(), svcCtx)
		resp, err := l.AISendResumeChat(&req)
		response.Response(r, w, resp, err)
	}
}
