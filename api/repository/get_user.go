package repository

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

// DB内のユーザ情報を取得
func (r *Repository) GetUserFromId(
	ctx context.Context, db service.Queryer, userId entity.UserId,
) (*entity.User, error) {
	u := &entity.User{}
	sql := `
		SELECT
			id,
			name,
			token,
			leader_card_id,
			created_at,
			updated_at
		FROM user
		WHERE id = ?
	`
	if err := db.GetContext(ctx, u, sql, userId); err != nil {
		return nil, err
	}
	if err := u.ValidateNotEmpty(); err != nil {
		return nil, err
	}
	return u, nil
}

// DB内のユーザ情報を取得
func (r *Repository) GetUserFromToken(
	ctx context.Context, db service.Queryer, userToken entity.UserTokenType,
) (*entity.User, error) {
	u := &entity.User{}
	sql := `SELECT
			id,
			name,
			token,
			leader_card_id,
			created_at,
			updated_at
		FROM user
		WHERE token = ?
	`
	if err := db.GetContext(ctx, u, sql, userToken); err != nil {
		return nil, err
	}
	if err := u.ValidateNotEmpty(); err != nil {
		return nil, err
	}
	return u, nil
}
