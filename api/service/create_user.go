package service

import (
	"context"
	"fmt"

	"github.com/pollenjp/gameserver-go/api/entity"
)

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out create_user_moq_test.go . CreateUserRepository
type CreateUserRepository interface {
	CreateUser(ctx context.Context, db Execer, u *entity.User) error
}

type CreateUser struct {
	DB   Execer
	Repo CreateUserRepository
}

func (ru *CreateUser) CreateUser(
	ctx context.Context,
	name string,
	leaderCard entity.LeaderCardIdIDType,
) (*entity.User, error) {
	u := &entity.User{
		Name: name,
	}

	if err := ru.Repo.CreateUser(ctx, ru.DB, u); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return u, nil
}
