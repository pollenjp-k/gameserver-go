package entity

type RoomUserStatus int

const (
	// Waiting	ホストがライブ開始ボタン押すのを待っている
	// End	ライブ終了
	RoomUserStatusWaiting  RoomUserStatus = 1
	RoomUserStatusFinished RoomUserStatus = 2
	RoomUserStatusLeaved   RoomUserStatus = 3
)

// TODO: 部屋の出入りとライブの終了は別のフラグで管理したほうが良いかもしれない
func HasScore(roomUserStatus RoomUserStatus) bool {
	switch {
	case roomUserStatus == RoomUserStatusFinished:
		return true
	case roomUserStatus == RoomUserStatusLeaved:
		return true
	default:
		return false
	}
}

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
