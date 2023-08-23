package repository

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

func (r *Repository) CreateRoom(
	ctx context.Context,
	db service.Execer,
	liveId entity.LiveId,
	hostUserId entity.UserId,
) (*entity.Room, error) {
	room := entity.NewRoom(
		liveId,
		hostUserId,
		entity.RoomStatusWaiting,
		r.Clocker.Now(),
		r.Clocker.Now(),
	)

	sql := `
	INSERT INTO
		room
		(
			live_id,
			host_user_id,
			status,
			created_at,
			updated_at
		)
	VALUES
		(?, ?, ?, ?, ?)
	;`

	result, err := db.ExecContext(
		ctx,
		sql,
		room.LiveId,
		room.HostUserId,
		room.Status,
		room.CreatedAt,
		room.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	room.Id = entity.RoomId(id)
	return room, nil
}
