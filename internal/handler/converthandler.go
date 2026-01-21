// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"

	"Go_SHORTENER/internal/logic"
	"Go_SHORTENER/internal/svc"
	"Go_SHORTENER/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ConvertHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//解析请求参数
		var req types.ConvertRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			fmt.Printf("解析失败: %v\n", err)
			return
		}
		fmt.Println("收到请求 Body:", r.Body)
		//参数规则校验
		if err := validator.New().StructCtx(r.Context(), &req); err != nil {
			logx.Errorw("validator check failed", logx.LogField{Key: "err", Value: err.Error()})
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		//执行业务逻辑
		l := logic.NewConvertLogic(r.Context(), svcCtx)
		resp, err := l.Convert(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
