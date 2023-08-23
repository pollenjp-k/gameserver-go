package entity

// Judge
//
// - Perfect
// - Great
// - Good
// - Bad
// - Miss
type Score struct {
	RoomId       RoomId `json:"room_id" db:"room_id"`
	UserId       UserId `json:"user_id" db:"user_id"`
	Score        int    `json:"score" db:"score"`
	JudgePerfect int    `json:"judge_perfect" db:"judge_perfect"`
	JudgeGreat   int    `json:"judge_great" db:"judge_great"`
	JudgeGood    int    `json:"judge_good" db:"judge_good"`
	JudgeBad     int    `json:"judge_bad" db:"judge_bad"`
	JudgeMiss    int    `json:"judge_miss" db:"judge_miss"`
}

func NewScore(
	roomId RoomId,
	userId UserId,
	score int,
	judgePerfect int,
	judgeGreat int,
	judgeGood int,
	judgeBad int,
	judgeMiss int,
) *Score {
	return &Score{
		RoomId:       roomId,
		UserId:       userId,
		Score:        score,
		JudgePerfect: judgePerfect,
		JudgeGreat:   judgeGreat,
		JudgeGood:    judgeGood,
		JudgeBad:     judgeBad,
		JudgeMiss:    judgeMiss,
	}
}
