package entity

type LiveDifficulty int

const (
	// https://github.com/KLabServerCamp/gameserver/blob/85b37d1c81bb7f4e7b3cba7875c3e0f84bfbcd54/docs/api.md#livedifficulty
	LiveDifficultyNormal LiveDifficulty = 1
	LiveDifficultyHard   LiveDifficulty = 2
)
