package entity

type RoomUser struct {
	RoomId         RoomId         `db:"room_id"`
	UserId         UserId         `db:"user_id"`
	LiveDifficulty LiveDifficulty `db:"live_difficulty"`
}

func NewRoomUser(
	roomId RoomId,
	userId UserId,
	liveDifficulty LiveDifficulty,
) *RoomUser {
	return &RoomUser{
		RoomId:         roomId,
		UserId:         userId,
		LiveDifficulty: liveDifficulty,
	}
}
