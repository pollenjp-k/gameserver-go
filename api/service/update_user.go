package service

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/repository"
)

//go:generate go run github.com/matryer/moq -out update_user_moq_test.go . UserUpdater
type UserUpdater interface {
	UpdateUser(ctx context.Context, db repository.Execer, newUser *entity.User) error
}

type UpdateUser struct {
	DB   repository.Execer
	Repo UserUpdater
}

func (ru *UpdateUser) UpdateUser(
	ctx context.Context,
	user *entity.User,
) error {
	if err := ru.Repo.UpdateUser(ctx, ru.DB, user); err != nil {
		return err
	}
	return nil
}
