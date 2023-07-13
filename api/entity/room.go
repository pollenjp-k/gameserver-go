package entity

import "time"

type RoomId int64
type LiveId int64

type Room struct {
	Id         RoomId     `db:"id"`
	LiveId     LiveId     `db:"live_id"`
	HostUserId UserId     `db:"host_user_id"`
	Status     RoomStatus `db:"status"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

func NewRoom(
	liveId LiveId,
	hostUseId UserId,
	status RoomStatus,
	createdAt time.Time,
	updatedAt time.Time,
) *Room {
	return &Room{
		LiveId:     liveId,
		HostUserId: hostUseId,
		Status:     status,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}
