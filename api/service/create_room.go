package service

import (
	"context"
	"fmt"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/repository"
)

//go:generate go run github.com/matryer/moq -out create_room_moq_test.go . CreateRoomRepository
type CreateRoomRepository interface {
	CreateRoom(
		ctx context.Context,
		db repository.Execer,
		liveId entity.LiveId,
		hostUserId entity.UserId,
	) (*entity.Room, error)
	CreateRoomUser(
		ctx context.Context,
		db repository.Execer,
		roomId entity.RoomId,
		userId entity.UserId,
		liveDifficulty entity.LiveDifficulty,
	) (*entity.RoomUser, error)
}

type CreateRoom struct {
	DB   repository.Beginner
	Repo CreateRoomRepository
}

func (cr *CreateRoom) CreateRoom(
	ctx context.Context,
	liveId entity.LiveId,
	hostUserId entity.UserId,
) (*entity.Room, *entity.RoomUser, error) {
	// helper functions
	fail := func(err error) (*entity.Room, *entity.RoomUser, error) {
		return nil, nil, fmt.Errorf("failed to register: %w", err)
	}

	tx, err := cr.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fail(fmt.Errorf("BeginTxx: %w", err))
	}
	defer tx.Rollback()

	room, err := cr.Repo.CreateRoom(ctx, tx, liveId, hostUserId)
	if err != nil {
		return fail(fmt.Errorf("CreateRoom: %w", err))
	}

	roomUser, err := cr.Repo.CreateRoomUser(ctx, tx, room.Id, hostUserId, entity.LiveDifficultyNormal)
	if err != nil {
		return fail(fmt.Errorf("CreateRoomUser: %w", err))
	}

	if err := tx.Commit(); err != nil {
		return fail(fmt.Errorf("commit: %w", err))
	}

	return room, roomUser, nil
}
