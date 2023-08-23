package repository

import (
	"context"
	"fmt"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

func (r *Repository) UpdateRoomUserStatus(
	ctx context.Context,
	db service.Execer,
	roomId entity.RoomId,
	userId entity.UserId,
	status entity.RoomUserStatus,
) error {
	sql := `
	UPDATE
		room_user
	SET
		status = ?
	WHERE
		room_id = ?
		AND
		user_id = ?
	;`

	if _, err := db.ExecContext(
		ctx,
		sql,
		status,
		roomId,
		userId,
	); err != nil {
		return fmt.Errorf("UpdateRoomUserStatus: %w", err)
	}

	return nil
}
