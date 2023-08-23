package repository

import (
	"context"
	"fmt"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

func (r *Repository) GetRoomUsers(
	ctx context.Context,
	db service.Queryer,
	roomId entity.RoomId,
) ([]*entity.RoomUser, error) {
	roomUsers := []*entity.RoomUser{}

	sql := `
	SELECT
		room_id,
		user_id,
		live_difficulty,
		status
	FROM
		room_user
	WHERE
		room_id = ?
	;`

	err := db.SelectContext(
		ctx,
		&roomUsers,
		sql,
		roomId,
	)
	if err != nil {
		return nil, fmt.Errorf("GetRoomUsers: %w", err)
	}
	return roomUsers, nil
}

func (r *Repository) GetRoomUsersByStatus(
	ctx context.Context,
	db service.Queryer,
	roomId entity.RoomId,
	status entity.RoomUserStatus,
) ([]*entity.RoomUser, error) {
	roomUsers := []*entity.RoomUser{}

	sql := `
	SELECT
		room_id,
		user_id,
		live_difficulty,
		status
	FROM
		room_user
	WHERE
		room_id = ?
		status = ?
	;`

	err := db.SelectContext(
		ctx,
		&roomUsers,
		sql,
		roomId,
		status,
	)
	if err != nil {
		return nil, fmt.Errorf("GetRoomUsers: %w", err)
	}
	return roomUsers, nil
}
