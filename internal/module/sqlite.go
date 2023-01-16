package module

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	DSN string
	DB  *sql.DB
}

func (m *SQLite) Init(ctx context.Context) (err error) {
	if m.DB, err = sql.Open("sqlite3", m.DSN); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return m.DB.PingContext(ctx)
}

func (m *SQLite) Close() error {
	return m.DB.Close()
}
