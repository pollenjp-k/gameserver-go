package repository

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

func (r *Repository) GetRoomList(
	ctx context.Context,
	db service.Queryer,
	RoomStatus entity.RoomStatus,
) ([]*service.RoomInfoItem, error) {
	roomList := []*service.RoomInfoItem{}

	sql := `
	WITH
		-- live_id と room.status で絞り込んだ room テーブル
		filtered_room AS (
			SELECT
				id,
				live_id,
				host_user_id,
				status,
				created_at
			FROM
				room
			WHERE
				status = ?
		)
	SELECT
		filtered_room.id AS "room_id",
		filtered_room.live_id AS "live_id",
		COUNT(room_user.user_id) AS "joined_user_count"
	FROM
		room_user
		INNER JOIN filtered_room
			ON
				room_user.room_id = filtered_room.id
	GROUP BY
		room_id
	ORDER BY
		filtered_room.created_at ASC
	;`

	err := db.SelectContext(
		ctx,
		&roomList,
		sql,
		RoomStatus,
	)
	if err != nil {
		return nil, err
	}
	return roomList, nil
}
