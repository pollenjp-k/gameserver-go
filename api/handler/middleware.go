package handler

import (
	"net/http"

	"github.com/pollenjp/gameserver-go/api/auth"
)

func AuthMiddleware(au *auth.Authorizer) func(next http.Handler) http.Handler {
	// request の context に認証情報を埋め込む
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, err := au.FillContext(r)
			if err != nil {
				RespondJson(r.Context(), w, ErrResponse{
					Message: "not find auth info",
					Details: []string{err.Error()},
				}, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}
