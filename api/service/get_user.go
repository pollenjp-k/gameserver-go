package service

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
)

//go:generate go run github.com/matryer/moq -out get_user_moq_test.go . UserGetter
type UserGetter interface {
	GetUserFromId(ctx context.Context, db Queryer, userId entity.UserId) (*entity.User, error)
}

type GetUser struct {
	DB   Queryer
	Repo UserGetter
}

func (ru *GetUser) GetUser(
	ctx context.Context,
	userId entity.UserId,
) (*entity.User, error) {
	u, err := ru.Repo.GetUserFromId(ctx, ru.DB, userId)
	if err != nil {
		return nil, err
	}
	return u, nil
}
