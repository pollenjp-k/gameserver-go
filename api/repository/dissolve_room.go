package repository

import (
	"context"
	"fmt"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

func (r *Repository) DissolveRoom(
	ctx context.Context, db service.Execer, roomId entity.RoomId,
) error {
	sql := `
		UPDATE
			room
		SET
			status = ?
			updated_at = ?
		WHERE
			id = ?
	;`
	if _, err := db.ExecContext(
		ctx,
		sql,
		entity.RoomStatusDissolution,
		r.Clocker.Now(),
		roomId,
	); err != nil {
		return fmt.Errorf("DissolveRoom: %w", err)
	}
	return nil
}
