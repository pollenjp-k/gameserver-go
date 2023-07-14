package entity

type RoomUserStatus int

const (
	// Waiting	ホストがライブ開始ボタン押すのを待っている
	// End	ライブ終了
	RoomUserStatusWaiting  RoomUserStatus = 1
	RoomUserStatusFinished RoomUserStatus = 2
)

type RoomUser struct {
	RoomId         RoomId         `db:"room_id"`
	UserId         UserId         `db:"user_id"`
	LiveDifficulty LiveDifficulty `db:"live_difficulty"`
	Status         RoomUserStatus `db:"status"`
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
		Status:         RoomUserStatusWaiting,
	}
}
