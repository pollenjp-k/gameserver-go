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
	GetRoomUsers(ctx context.Context, db Queryer, roomId entity.RoomId) ([]*entity.RoomUser, error)
}

type GetRoomResult struct {
	DB   Queryer
	Repo GetRoomResultRepository
}

func (grr *GetRoomResult) GetRoomResult(
	ctx context.Context,
	roomId entity.RoomId,
) ([]*RoomUserResult, error) {
	roomUsers, err := grr.Repo.GetRoomUsers(ctx, grr.DB, roomId)
	if err != nil {
		return nil, err
	}

	m := make(map[entity.UserId]*entity.RoomUser, len(roomUsers))
	for _, ru := range roomUsers {
		m[ru.UserId] = ru
	}

	userId, ok := GetUserId(ctx)
	if !ok {
		return nil, &entity.ErrUnauthorized{}
	}

	// ルームに参加していないユーザーは結果を見れない
	if _, ok := m[userId]; !ok {
		return nil, &entity.ErrPermissionDenied{}
	}

	// 全員が終わっていない場合は結果を見れない
	allFinished := func(roomUsers []*entity.RoomUser) bool {
		for _, ru := range roomUsers {
			if !entity.HasScore(ru.Status) {
				return false
			}
		}
		return true
	}
	if !allFinished(roomUsers) {
		log.Printf("GetRoomResult: not all users finished yet")
		return []*RoomUserResult{}, nil
	}

	userAndScores, err := grr.Repo.GetRoomUserAndScoreInRoom(ctx, grr.DB, roomId)
	if err != nil {
		return nil, err
	}

	roomUserResults := make([]*RoomUserResult, len(userAndScores))
	for i, us := range userAndScores {
		roomUserResults[i] = &RoomUserResult{
			UserId:       us.UserId,
			Score:        us.Score,
			JudgePerfect: us.JudgePerfect,
			JudgeGreat:   us.JudgeGreat,
			JudgeGood:    us.JudgeGood,
			JudgeBad:     us.JudgeBad,
			JudgeMiss:    us.JudgeMiss,
		}
	}

	return roomUserResults, nil
}
