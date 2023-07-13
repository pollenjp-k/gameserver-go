package repository

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
)

// User情報として以下の内容を更新
//
// - Name
// - LeaderCardId
// - UpdatedAt (自動設定)
func (r *Repository) UpdateUser(
	ctx context.Context, db Execer, newUser *entity.User,
) error {
	sql := `
		UPDATE
			user
		SET
			name = ?,
			leader_card_id = ?,
			updated_at = ?
		WHERE id = ?
	`
	if _, err := db.ExecContext(
		ctx,
		sql,
		newUser.Name,
		newUser.LeaderCardId,
		r.Clocker.Now(),
		newUser.Id,
	); err != nil {
		return err
	}
	return nil
}
