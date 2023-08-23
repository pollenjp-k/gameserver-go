package entity

type RoomStatus int

const (
	// https://github.com/KLabServerCamp/gameserver/blob/85b37d1c81bb7f4e7b3cba7875c3e0f84bfbcd54/docs/api.md
	// Waiting	1	ホストがライブ開始ボタン押すのを待っている
	// LiveStart	2	ライブ画面遷移OK
	// Dissolution	3	解散された
	RoomStatusWaiting     RoomStatus = 1
	RoomStatusLiveStart   RoomStatus = 2
	RoomStatusDissolution RoomStatus = 3
)
