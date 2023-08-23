package repository

import (
	"context"
	"fmt"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

func (r *Repository) CreateScore(
	ctx context.Context,
	db service.Execer,
	score *entity.Score,
) error {
	sql := `
	INSERT INTO
		score
		(
			room_id,
			user_id,
			score,
			judge_perfect,
			judge_great,
			judge_good,
			judge_bad,
			judge_miss
		)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?)
	;`

	_, err := db.ExecContext(
		ctx,
		sql,
		score.RoomId,
		score.UserId,
		score.Score,
		score.JudgePerfect,
		score.JudgeGreat,
		score.JudgeGood,
		score.JudgeBad,
		score.JudgeMiss,
	)
	if err != nil {
		return fmt.Errorf("CreateScore: %w", err)
	}
	return nil
}
