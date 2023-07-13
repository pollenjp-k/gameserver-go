package service

import (
	"context"
	"fmt"

	"github.com/pollenjp/gameserver-go/api/entity"
)

// Repository からの受け取りに利用.
// Repository からの受け取りでは IsHost のみ Dummy Value である.
type WaitingRoomUser struct {
	UserId           entity.UserId             `json:"user_id" db:"user_id"`
	Name             string                    `json:"name" db:"name"`
	LeaderCardId     entity.LeaderCardIdIDType `json:"leader_card_id" db:"leader_card_id"`
	SelectDifficulty entity.LiveDifficulty     `json:"select_difficulty" db:"select_difficulty"`
	IsHost           bool                      `json:"is_host" db:"is_host"`
}

// handler への返り値に利用.
type WaitRoomResult struct {
	Room            *entity.Room
	WaitingRoomUser []*WaitingRoomUser
}

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out create_room_moq_test.go . CreateRoomRepository
type WaitRoomRepository interface {
	GetRoom(
		ctx context.Context,
		db Queryer,
		roomId entity.RoomId,
	) (*entity.Room, error)
	GetWaitingUsers(
		ctx context.Context,
		db Queryer,
		roomId entity.RoomId,
	) ([]*WaitingRoomUser, error)
}

type WaitRoom struct {
	DB   Queryer
	Repo WaitRoomRepository
}

func (cr *WaitRoom) WaitRoom(
	ctx context.Context,
	roomId entity.RoomId,
	userId entity.UserId,
) (*WaitRoomResult, error) {
	// helper functions
	fail := func(err error) (*WaitRoomResult, error) {
		return nil, fmt.Errorf("WaitRoom: %w", err)
	}

	db := cr.DB

	room, err := cr.Repo.GetRoom(ctx, db, roomId)
	if err != nil {
		return fail(err)
	}

	waitingRoomUser, err := cr.Repo.GetWaitingUsers(ctx, db, roomId)
	if err != nil {
		return fail(err)
	}

	for _, user := range waitingRoomUser {
		if user.UserId == room.HostUserId {
			user.IsHost = true
		}
	}

	waitRoomResult := &WaitRoomResult{
		Room:            room,
		WaitingRoomUser: waitingRoomUser,
	}

	return waitRoomResult, nil
}
