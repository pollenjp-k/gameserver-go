package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pollenjp/gameserver-go/api/entity"
)

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out start_room_moq_test.go . StartRoomRepository
type StartRoomRepository interface {
	GetRoom(
		ctx context.Context,
		db Queryer,
		roomId entity.RoomId,
	) (*entity.Room, error)
	UpdateRoomStatus(
		ctx context.Context,
		db Execer,
		roomId entity.RoomId,
		status entity.RoomStatus,
	) error
}

type StartRoom struct {
	DB   Beginner
	Repo StartRoomRepository
}

func (cr *StartRoom) StartRoom(
	ctx context.Context,
	roomId entity.RoomId,
	hostUserId entity.UserId,
) error {
	// helper functions
	fail := func(err error) error {
		return err
	}
	failWithRollBack := func(tx *sqlx.Tx, err error) error {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = fmt.Errorf("rollbacking: %w: %v", rollbackErr, err)
		}
		return fail(err)
	}

	tx, err := cr.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fail(fmt.Errorf("BeginTxx: %w", err))
	}

	room, err := cr.Repo.GetRoom(ctx, tx, roomId)
	if err != nil {
		return failWithRollBack(tx, err)
	}

	if room.HostUserId != hostUserId {
		return failWithRollBack(tx, fmt.Errorf("hostUser mismatch: %v != %v", room.HostUserId, hostUserId))
	}

	if room.Status != entity.RoomStatusWaiting {
		return failWithRollBack(tx, fmt.Errorf("room status: %v", room.Status))
	}

	if err := cr.Repo.UpdateRoomStatus(ctx, tx, roomId, entity.RoomStatusLiveStart); err != nil {
		return failWithRollBack(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return failWithRollBack(tx, fmt.Errorf("committing: %w", err))
	}

	return nil
}
