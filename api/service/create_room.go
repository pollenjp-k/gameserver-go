package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pollenjp/gameserver-go/api/entity"
)

//go:generate go run github.com/matryer/moq -out create_room_moq_test.go . CreateRoomRepository
type CreateRoomRepository interface {
	CreateRoom(
		ctx context.Context,
		db Execer,
		liveId entity.LiveId,
		hostUserId entity.UserId,
	) (*entity.Room, error)
	CreateRoomUser(
		ctx context.Context,
		db Execer,
		roomId entity.RoomId,
		userId entity.UserId,
		liveDifficulty entity.LiveDifficulty,
	) (*entity.RoomUser, error)
}

type CreateRoom struct {
	DB   Beginner
	Repo CreateRoomRepository
}

func (cr *CreateRoom) CreateRoom(
	ctx context.Context,
	liveId entity.LiveId,
	hostUserId entity.UserId,
) (*entity.Room, *entity.RoomUser, error) {
	// helper functions
	failWithRollBack := func(tx *sqlx.Tx, err error) (*entity.Room, *entity.RoomUser, error) {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = fmt.Errorf("rollbacking: %w: %v", rollbackErr, err)
		}
		return nil, nil, err
	}

	tx, err := cr.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("BeginTxx: %w", err)
	}

	room, err := cr.Repo.CreateRoom(ctx, tx, liveId, hostUserId)
	if err != nil {
		return failWithRollBack(tx, fmt.Errorf("CreateRoom: %w", err))
	}

	roomUser, err := cr.Repo.CreateRoomUser(ctx, tx, room.Id, hostUserId, entity.LiveDifficultyNormal)
	if err != nil {
		return failWithRollBack(tx, fmt.Errorf("CreateRoomUser: %w", err))
	}

	if err := tx.Commit(); err != nil {
		return failWithRollBack(tx, fmt.Errorf("committing: %w", err))
	}

	return room, roomUser, nil
}
