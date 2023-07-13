package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

// user table にユーザを追加
//
// - DBに登録
// - 以下の値を設定する
//   - `entity.User.ID`
//   - `entity.User.Created`
//   - `entity.User.Modified`
func (r *Repository) CreateUser(
	ctx context.Context, db service.Execer, u *entity.User,
) error {
	u.Token = entity.UserTokenType(uuid.NewString())
	u.CreatedAt = r.Clocker.Now()
	u.UpdatedAt = r.Clocker.Now()

	// TODO: token が unique でなかった場合に5回程リトライする

	sql := `INSERT INTO
		user (
			name,
			token,
			leader_card_id,
			created_at,
			updated_at
		)
	VALUES
		(?, ?, ?, ?, ?)
	;`

	result, err := db.ExecContext(
		ctx,
		sql,
		u.Name,
		u.Token,
		u.LeaderCardId,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError

		// primary key 重複エラー
		if errors.As(err, &mysqlErr) && mysqlErr.Number == service.ErrCodeMySQLDuplicateEntry {
			return fmt.Errorf("cannot create same name user: %w", service.ErrAlreadyEntry)
		}

		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.Id = entity.UserId(id)
	if err := u.ValidateNotEmpty(); err != nil {
		return err
	}
	return nil
}
