package repository

import (
	"context"
	"fmt"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

func (r *Repository) UpdateRoomStatus(
	ctx context.Context,
	db service.Execer,
	roomId entity.RoomId,
	status entity.RoomStatus,
) error {
	sql := `
	UPDATE
		room
	SET
		status = ?
	WHERE
		id = ?
	;`

	if _, err := db.ExecContext(
		ctx,
		sql,
		status,
		roomId,
	); err != nil {
		return fmt.Errorf("UpdateRoomStatus: %w", err)
	}

	return nil
}
