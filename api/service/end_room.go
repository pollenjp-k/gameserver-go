package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pollenjp/gameserver-go/api/entity"
)

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out end_room_list_moq_test.go . EndRoomRepository
type EndRoomRepository interface {
	CreateScore(
		ctx context.Context,
		db Execer,
		score *entity.Score,
	) error

	UpdateRoomUserStatus(
		ctx context.Context,
		db Execer,
		roomId entity.RoomId,
		userId entity.UserId,
		status entity.RoomUserStatus,
	) error
}

type EndRoom struct {
	DB   Beginner
	Repo EndRoomRepository
}

// - Score の格納
// - RoomUser の状態を変更する end など
func (er *EndRoom) EndRoom(
	ctx context.Context,
	score *entity.Score,
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

	tx, err := er.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fail(fmt.Errorf("BeginTxx: %w", err))
	}

	if err := er.Repo.UpdateRoomUserStatus(ctx, tx, score.RoomId, score.UserId, entity.RoomUserStatusFinished); err != nil {
		// TODO: error が起きた場合でも Rollback せずに Status は End にしたほうが良いのか？
		return failWithRollBack(tx, err)
	}

	if err := er.Repo.CreateScore(ctx, tx, score); err != nil {
		return failWithRollBack(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return failWithRollBack(tx, fmt.Errorf("committing: %w", err))
	}

	return nil
}
