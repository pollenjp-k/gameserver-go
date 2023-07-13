package service

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
)

//go:generate go run github.com/matryer/moq -out update_user_moq_test.go . UserUpdater
type UserUpdater interface {
	UpdateUser(ctx context.Context, db Execer, newUser *entity.User) error
}

type UpdateUser struct {
	DB   Execer
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
