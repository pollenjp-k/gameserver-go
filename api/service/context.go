package service

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
)

type userIDKey struct{}

func SetUserId(ctx context.Context, uid entity.UserId) context.Context {
	return context.WithValue(ctx, userIDKey{}, uid)
}

func GetUserId(ctx context.Context) (entity.UserId, bool) {
	id, ok := ctx.Value(userIDKey{}).(entity.UserId)
	return id, ok
}
