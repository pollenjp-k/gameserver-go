package repository

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
)

func (r *Repository) CreateRoomUser(
	ctx context.Context,
	db Execer,
	roomId entity.RoomId,
	userId entity.UserId,
	liveDifficulty entity.LiveDifficulty,
) (*entity.RoomUser, error) {
	roomUser := entity.NewRoomUser(
		roomId,
		userId,
		liveDifficulty,
	)

	sql := `
	INSERT INTO
		room_user
		(
			room_id,
			user_id,
			live_difficulty
		)
	VALUES
		(?, ?, ?)
	;`

	_, err := db.ExecContext(
		ctx,
		sql,
		roomUser.RoomId,
		roomUser.UserId,
		roomUser.LiveDifficulty,
	)
	if err != nil {
		return nil, err
	}
	return roomUser, nil
}
