package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pollenjp/gameserver-go/api/clock"
	"github.com/pollenjp/gameserver-go/api/config"
)

const (
	// ErrCodeMySQLDuplicateEntry はMySQL系のDUPLICATEエラーコード
	// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
	// Error number: 1062; Symbol: ER_DUP_ENTRY; SQLSTATE: 23000
	ErrCodeMySQLDuplicateEntry = 1062
)

var (
	ErrAlreadyEntry = errors.New("duplicate entry")
)

// defer のように複数のCleanUp処理を渡せるようにする
// @pollenjp のお遊び
type CleanUpContainer struct {
	funcs []func()
}

func (c *CleanUpContainer) Add(f func()) {
	// 戦闘に挿入
	c.funcs = append([]func(){f}, c.funcs...)
}

func (c *CleanUpContainer) GetCleanUp() func() {
	return func() {
		for _, f := range c.funcs {
			f()
		}
	}
}

// databaseとのコネクションを確立する
// return. (db, cleanup func, error)
func New(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {
	cleanUpContainer := &CleanUpContainer{}

	cleanUpContainer.Add(func() {
		log.Printf("1st cleanup")
	})

	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.DBUser, cfg.DBPassword,
			cfg.DBHost, cfg.DBPort,
			cfg.DBName,
		),
	)
	if err != nil {
		return nil, cleanUpContainer.GetCleanUp(), err
	}
	cleanUpContainer.Add(func() {
		log.Printf("db close cleanup")
		_ = db.Close()
	})

	pingDB := func(trial int) error {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			log.Printf("DB Connection (trial:%d): %s", trial, err.Error())
			return err
		}
		return nil
	}

	trial := 0
	for {
		trial++
		if err := pingDB(trial); err != nil {
			if trial > 30 {
				log.Println("Couldn't connect to database.")
				return nil, cleanUpContainer.GetCleanUp(), err
			}
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("Database is up. Starting server...")
		break
	}

	xdb := sqlx.NewDb(db, "mysql")
	return xdb, cleanUpContainer.GetCleanUp(), nil
}

// Repository はデータベースへのアクセスを提供する
type Repository struct {
	Clocker clock.Clocker
}

type Beginner interface {
	// https://pkg.go.dev/github.com/jmoiron/sqlx#DB.BeginTxx
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

type Preparer interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

// Execer はデータベースへのExec系クエリを提供する
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

// Queryer はデータベースへのクエリを提供する
type Queryer interface {
	Preparer
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
}

var (
	_ Beginner = (*sqlx.DB)(nil)
	_ Preparer = (*sqlx.DB)(nil)
	_ Queryer  = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.Tx)(nil)
)
