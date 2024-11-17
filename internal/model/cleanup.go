package model

import (
	"context"
	"time"
	"tournament/pkg/db"
)

func Cleanup(ctx context.Context) (err error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second) // here give more time
	defer cancel()

	_, err = conn.Exec(timeoutCtx, `
		BEGIN;
		TRUNCATE TABLE matches, teams RESTART IDENTITY CASCADE;
		--- if you have more cleanup operations
		COMMIT;
	`)

	return
}
