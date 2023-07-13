package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pollenjp/gameserver-go/api/config"
	"github.com/pollenjp/gameserver-go/api/entity"
)

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out create_room_moq_test.go . CreateRoomRepository
type JoinRoomRepository interface {
	GetRoom(
		ctx context.Context,
		db Queryer,
		roomId entity.RoomId,
	) (*entity.Room, error)
	GetRoomUsers(
		ctx context.Context,
		db Queryer,
		roomId entity.RoomId,
	) ([]*entity.RoomUser, error)
	CreateRoomUser(
		ctx context.Context,
		db Execer,
		roomId entity.RoomId,
		userId entity.UserId,
		liveDifficulty entity.LiveDifficulty,
	) (*entity.RoomUser, error)
}

type JoinRoom struct {
	DB   Beginner
	Repo JoinRoomRepository
}

func (cr *JoinRoom) JoinRoom(
	ctx context.Context,
	roomId entity.RoomId,
	userId entity.UserId,
	liveDifficulty entity.LiveDifficulty,
) (entity.JoinRoomResult, error) {
	// helper functions
	fail := func(err error) (entity.JoinRoomResult, error) {
		return entity.JoinRoomResultOtherErr, err
	}
	failWithRollBack := func(tx *sqlx.Tx, err error) (entity.JoinRoomResult, error) {
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

	switch room.Status {
	case entity.RoomStatusWaiting:
		// do nothing
	case entity.RoomStatusLiveStart:
		if result, err := failWithRollBack(tx, nil); err != nil {
			return result, err
		}
		return entity.JoinRoomResultOtherErr, nil
	case entity.RoomStatusDissolution:
		if result, err := failWithRollBack(tx, nil); err != nil {
			return result, err
		}
		return entity.JoinRoomResultDisbanded, nil
	default:
		return failWithRollBack(tx, fmt.Errorf("unknown room status: %v", room.Status))
	}

	// check the number of users in the room
	roomUsers, err := cr.Repo.GetRoomUsers(ctx, tx, roomId)
	if err != nil {
		return failWithRollBack(tx, err)
	}
	if len(roomUsers) >= config.MaxUserCount {
		if result, err := failWithRollBack(tx, nil); err != nil {
			return result, err
		}
		return entity.JoinRoomResultRoomFull, nil
	}

	// if the user is already in the room, return the result
	for _, roomUser := range roomUsers {
		if roomUser.UserId == userId {
			if result, err := failWithRollBack(tx, nil); err != nil {
				return result, err
			}
			return entity.JoinRoomResultOtherErr, nil
		}
	}

	if _, err := cr.Repo.CreateRoomUser(ctx, tx, room.Id, userId, liveDifficulty); err != nil {
		return failWithRollBack(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return failWithRollBack(tx, fmt.Errorf("committing: %w", err))
	}

	return entity.JoinRoomResultOk, nil
}
