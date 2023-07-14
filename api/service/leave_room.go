package service

import (
	"context"
	"fmt"

	"github.com/pollenjp/gameserver-go/api/entity"
)

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out leave_room_moq_test.go . LeaveRoomRepository
type LeaveRoomRepository interface {
	// RoomUser.Status を Leaved にする
	LeaveRoom(
		ctx context.Context,
		db Execer,
		roomId entity.RoomId,
		userId entity.UserId,
	) error
	DissolveRoom(
		ctx context.Context,
		db Execer,
		roomId entity.RoomId,
	) error
	GetRoomUsers(
		ctx context.Context,
		db Queryer,
		roomId entity.RoomId,
	) ([]*entity.RoomUser, error)
}

type LeaveRoom struct {
	DB   QueryerAndExecer
	Repo LeaveRoomRepository
}

// トランザクションは貼らなくて良いはず...?
func (cr *LeaveRoom) LeaveRoom(
	ctx context.Context,
	roomId entity.RoomId,
	userId entity.UserId,
) error {
	// helper functions
	fail := func(err error) error {
		return fmt.Errorf("LeaveRoom: %w", err)
	}

	db := cr.DB

	if err := cr.Repo.LeaveRoom(ctx, db, roomId, userId); err != nil {
		return fail(err)
	}

	roomUsers, err := cr.Repo.GetRoomUsers(ctx, db, roomId)
	if err != nil {
		return fail(err)
	}
	for _, roomUser := range roomUsers {
		if roomUser.Status != entity.RoomUserStatusLeaved {
			// まだ抜けていない人がいれば、ルームを解散しない
			return nil
		}
	}

	if err := cr.Repo.DissolveRoom(ctx, db, roomId); err != nil {
		return fail(err)
	}

	return nil
}
