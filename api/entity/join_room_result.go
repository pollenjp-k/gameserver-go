package entity

type JoinRoomResult int

const (
	// https://github.com/KLabServerCamp/gameserver/blob/85b37d1c81bb7f4e7b3cba7875c3e0f84bfbcd54/docs/api.md
	// name	value	memo
	// Ok	1	入場OK
	// RoomFull	2	満員
	// Disbanded	3	解散済み
	// OtherError	4	その他エラー
	JoinRoomResultOk        JoinRoomResult = 1
	JoinRoomResultRoomFull  JoinRoomResult = 2
	JoinRoomResultDisbanded JoinRoomResult = 3
	JoinRoomResultOtherErr  JoinRoomResult = 4
)
