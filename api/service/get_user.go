package service

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/repository"
)

//go:generate go run github.com/matryer/moq -out get_user_moq_test.go . UserGetter
type UserGetter interface {
	GetUserFromId(ctx context.Context, db repository.Queryer, userId entity.UserID) (*entity.User, error)
}

type GetUser struct {
	DB   repository.Queryer
	Repo UserGetter
}

func (ru *GetUser) GetUser(
	ctx context.Context,
	userId entity.UserID,
) (*entity.User, error) {
	u, err := ru.Repo.GetUserFromId(ctx, ru.DB, userId)
	if err != nil {
		return nil, err
	}
	return u, nil
}
