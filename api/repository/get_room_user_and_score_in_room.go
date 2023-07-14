package repository

import (
	"context"
	"fmt"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

func (r *Repository) GetRoomUserAndScoreInRoom(
	ctx context.Context,
	db service.Queryer,
	roomId entity.RoomId,
) ([]*service.RoomUserAndScore, error) {
	roomUserAndScoreList := []*service.RoomUserAndScore{}

	sql := `
	SELECT
		room_user.user_id AS user_id,
		room_user.status AS user_status,
		score.score AS score,
		score.judge_perfect AS judge_perfect,
		score.judge_great AS judge_great,
		score.judge_good AS judge_good,
		score.judge_bad AS judge_bad,
		score.judge_miss AS judge_miss
	FROM
		room_user
		INNER JOIN score
			ON
				room_user.room_id = score.room_id
				AND
				room_user.user_id = score.user_id
	WHERE
		room_user.room_id = ?
	ORDER BY
		room_user.user_id ASC;
	;`

	err := db.SelectContext(
		ctx,
		&roomUserAndScoreList,
		sql,
		roomId,
	)
	if err != nil {
		return nil, fmt.Errorf("GetRoomUserAndScoreInRoom: %w", err)
	}
	return roomUserAndScoreList, nil
}
