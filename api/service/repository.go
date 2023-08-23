package service

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

type QueryerAndExecer interface {
	Queryer
	Execer
}

var (
	_ Beginner         = (*sqlx.DB)(nil)
	_ Preparer         = (*sqlx.DB)(nil)
	_ Queryer          = (*sqlx.DB)(nil)
	_ Execer           = (*sqlx.DB)(nil)
	_ QueryerAndExecer = (*sqlx.DB)(nil)
	_ Execer           = (*sqlx.Tx)(nil)
)
