package auth

import (
	"context"
	"net/http"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/repository"
)

//go:generate go run github.com/matryer/moq -out auth_moq_test.go . AuthRepository
type AuthRepository interface {
	GetUserFromToken(ctx context.Context, db repository.Queryer, userToken entity.UserTokenType) (*entity.User, error)
}

func NewAuthorizer(db repository.Queryer, repo AuthRepository) *Authorizer {
	return &Authorizer{
		DB:   db,
		Repo: repo,
	}
}

type Authorizer struct {
	DB   repository.Queryer
	Repo AuthRepository
}

type userIDKey struct{}

// *http.Request型から認証情報を context に書き込む
func (au *Authorizer) FillContext(r *http.Request) (*http.Request, error) {
	token, err := ExtractBearerToken(r)
	if err != nil {
		return nil, err
	}

	u, err := au.Repo.GetUserFromToken(r.Context(), au.DB, token)
	if err != nil {
		return nil, err
	}

	ctx := SetUserId(r.Context(), u.Id)

	clone := r.Clone(ctx)
	return clone, nil
}

func SetUserId(ctx context.Context, uid entity.UserId) context.Context {
	return context.WithValue(ctx, userIDKey{}, uid)
}

func GetUserId(ctx context.Context) (entity.UserId, bool) {
	id, ok := ctx.Value(userIDKey{}).(entity.UserId)
	return id, ok
}
