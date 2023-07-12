package repository

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
)

// DB内のユーザ情報を取得
func (r *Repository) GetUser(
	ctx context.Context, db Queryer, userId entity.UserID,
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
