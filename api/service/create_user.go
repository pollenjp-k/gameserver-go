package service

import (
	"context"
	"fmt"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/repository"
)

//go:generate go run github.com/matryer/moq -out register_user_test.go . UserRegistrar
type UserRegistrar interface {
	CreateUser(ctx context.Context, db repository.Execer, u *entity.User) error
}

type CreateUser struct {
	DB   repository.Execer
	Repo UserRegistrar
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
