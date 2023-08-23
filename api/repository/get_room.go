package repository

import (
	"context"
	"fmt"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

func (r *Repository) GetRoom(
	ctx context.Context,
	db service.Queryer,
	roomId entity.RoomId,
) (*entity.Room, error) {
	room := &entity.Room{}

	sql := `
	SELECT
		id,
		live_id,
		host_user_id,
		status,
		created_at,
		updated_at
	FROM
		room
	WHERE
		id = ?
	;`

	err := db.GetContext(
		ctx,
		room,
		sql,
		roomId,
	)
	if err != nil {
		return nil, fmt.Errorf("GetRoom: %w", err)
	}
	return room, nil
}
