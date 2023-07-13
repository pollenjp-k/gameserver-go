package repository

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

func (r *Repository) GetWaitingUsers(
	ctx context.Context,
	db service.Queryer,
	roomId entity.RoomId,
) ([]*service.WaitingRoomUser, error) {
	waitingRoomUser := []*service.WaitingRoomUser{}

	sql := `
	SELECT
		user.id AS "user_id",
		user.name AS "name",
		user.leader_card_id AS "leader_card_id",
		room_user.live_difficulty AS "select_difficulty"
	FROM
		room_user
		INNER JOIN user
			ON
				room_user.user_id = user.id
	WHERE
		room_user.room_id = ?
	;`

	err := db.SelectContext(
		ctx,
		&waitingRoomUser,
		sql,
		roomId,
	)
	if err != nil {
		return nil, err
	}
	return waitingRoomUser, nil
}
