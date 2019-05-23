package healthz

import (
	"context"
	"database/sql"
)

// CheckDBPing pings a database. Assumes your database driver supports Ping.
type CheckDBPing struct {
	db *sql.DB
}

// NewCheckDBPing initializes a CheckDBPing.
func NewCheckDBPing(db *sql.DB) CheckDBPing {
	return CheckDBPing{
		db: db,
	}
}

// Check is called by the checker and attempts to ping the database.
func (c CheckDBPing) Check(ctx context.Context) *Response {
	err := c.db.PingContext(ctx)
	return &Response{Error: err}
}
