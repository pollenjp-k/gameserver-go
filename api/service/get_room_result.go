package service

import (
	"context"
	"log"

	"github.com/pollenjp/gameserver-go/api/entity"
)

// handler へのレスポンス
type RoomUserResult struct {
	UserId       entity.UserId
	Score        int
	JudgePerfect int
	JudgeGreat   int
	JudgeGood    int
	JudgeBad     int
	JudgeMiss    int
}

func NewRoomUserResult(
	userId entity.UserId,
	score int,
	judgePerfect int,
	judgeGreat int,
	judgeGood int,
	judgeBad int,
	judgeMiss int,
) *RoomUserResult {
	return &RoomUserResult{
		UserId:       userId,
		Score:        score,
		JudgePerfect: judgePerfect,
		JudgeGreat:   judgeGreat,
		JudgeGood:    judgeGood,
		JudgeBad:     judgeBad,
		JudgeMiss:    judgeMiss,
	}
}

// Repository からの受け取り
type RoomUserAndScore struct {
	UserId       entity.UserId         `db:"user_id"`
	UserStatus   entity.RoomUserStatus `db:"user_status"`
	Score        int                   `db:"score"`
	JudgePerfect int                   `db:"judge_perfect"`
	JudgeGreat   int                   `db:"judge_great"`
	JudgeGood    int                   `db:"judge_good"`
	JudgeBad     int                   `db:"judge_bad"`
	JudgeMiss    int                   `db:"judge_miss"`
}

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out get_room_result_moq_test.go . GetRoomResultRepository
type GetRoomResultRepository interface {
	GetRoomUserAndScoreInRoom(ctx context.Context, db Queryer, roomId entity.RoomId) ([]*RoomUserAndScore, error)
	GetRoomUsers(
		ctx context.Context,
		db Queryer,
		roomId entity.RoomId,
	) ([]*entity.RoomUser, error)
}

type GetRoomResult struct {
	DB   Queryer
	Repo GetRoomResultRepository
}

type RoomUserResultList []*RoomUserResult

func (grr *GetRoomResult) GetRoomResult(
	ctx context.Context,
	roomId entity.RoomId,
) (RoomUserResultList, error) {
	// TODO: roomId auth check
	// ルームに参加していないユーザーは結果を見れない
	// userId, ok := GetUserId(ctx)
	// if !ok {
	// 	return nil, &entity.ErrUnauthorized{}
	// }
	// m := make(map[entity.UserId]*entity.RoomUser, len(roomUsers))
	// for _, ru := range roomUsers {
	// 	m[ru.UserId] = ru
	// }

	roomUsers, err := grr.Repo.GetRoomUsers(ctx, grr.DB, roomId)
	if err != nil {
		return nil, err
	}

	Status2RoomUser := make(map[entity.RoomUserStatus]*entity.RoomUser)
	UserId2RoomUser := make(map[entity.UserId]*entity.RoomUser)
	for _, ru := range roomUsers {
		Status2RoomUser[ru.Status] = ru
		UserId2RoomUser[ru.UserId] = ru
	}

	// もし WaitingUser の人がいる場合は結果を見れない
	if _, ok := Status2RoomUser[entity.RoomUserStatusWaiting]; ok {
		log.Printf("GetRoomResult: waiting user exists")
		return RoomUserResultList{}, nil
	}

	userAndScores, err := grr.Repo.GetRoomUserAndScoreInRoom(ctx, grr.DB, roomId)
	if err != nil {
		return nil, err
	}

	roomUserResults := make(RoomUserResultList, len(userAndScores))
	for i, us := range userAndScores {
		roomUserResults[i] = NewRoomUserResult(
			us.UserId,
			us.Score,
			us.JudgePerfect,
			us.JudgeGreat,
			us.JudgeGood,
			us.JudgeBad,
			us.JudgeMiss,
		)
	}

	return roomUserResults, nil
}
