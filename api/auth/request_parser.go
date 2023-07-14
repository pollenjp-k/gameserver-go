package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pollenjp/gameserver-go/api/entity"
)

func ExtractBearerToken(r *http.Request) (entity.UserTokenType, error) {
	// Authorizationヘッダーの値を取得します
	authHeader := r.Header.Get("Authorization")

	// zero value
	var token entity.UserTokenType

	// Authorizationヘッダーが存在しない場合はエラーを返します
	if authHeader == "" {
		return token, fmt.Errorf("authorization header does not exist")
	}

	// "Bearer <トークン>"の形式であることを確認します
	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
		log.Printf("authParts: %v", authParts)
		return token, fmt.Errorf("authorization header format must be 'Bearer <token>'")
	}

	token = entity.UserTokenType(authParts[1])
	return token, nil
}
