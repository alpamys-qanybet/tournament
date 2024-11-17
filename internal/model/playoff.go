package model

import (
	"context"
	"errors"
	"time"
	"tournament/dto"
	"tournament/pkg/db"

	"github.com/jackc/pgx/v4"
)

func PreparePlayoffMatches(ctx context.Context, matches []*dto.MatchDTO) (err error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second) // here give more time
	defer cancel()

	tx, err := conn.BeginTx(timeoutCtx, pgx.TxOptions{}) // Start the transaction
	if err != nil {
		return
	}

	defer func() { // rollback on error
		if err != nil {
			tx.Rollback(timeoutCtx)
		}
	}()

	for _, m := range matches {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		default:
			queryCtx, queryCancel := context.WithTimeout(timeoutCtx, 3*time.Second) // separate timeout for query
			defer queryCancel()

			_, err = tx.Exec(queryCtx, `
				INSERT INTO matches(first_team_id, second_team_id, match_type)
				VALUES ($1, $2, $3)
			`, m.FirstTeamId, m.SecondTeamId, m.MatchType)
			if err != nil {
				return
			}
		}
	}

	// Commit the transaction
	if err = tx.Commit(timeoutCtx); err != nil {
		return
	}

	return
}

func PlayoffIsPrepared(ctx context.Context, matchType string, checkPlayed bool) (prepared bool, err error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count uint16
	err = conn.QueryRow(timeoutCtx, `
		SELECT count(*)
		FROM matches
		WHERE match_type = $1
		AND winner_id IS NULL;
	`, matchType).Scan(&count)
	if err != nil {
		return
	}

	prepared = count > 0 // matches prepared

	if checkPlayed {
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var countPlayoff uint16
		err = conn.QueryRow(timeoutCtx, `
			SELECT count(*)
			FROM matches
			WHERE match_type = $1;
		`, matchType).Scan(&countPlayoff)
		if err != nil {
			return
		}

		if !prepared && countPlayoff > 0 {
			if matchType == dto.MatchTypePlayoffQuarter {
				err = errors.New("playoff_quarter_is_already_started")
			} else if matchType == dto.MatchTypePlayoffSemi {
				err = errors.New("playoff_semi_is_already_started")
			} else if matchType == dto.MatchTypePlayoffFinal {
				err = errors.New("playoff_final_is_already_started")
			}
		}
	}

	return
}

func PlayoffIsStarted(ctx context.Context, matchType string) (prepared bool, err error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count uint16
	err = conn.QueryRow(timeoutCtx, `
		SELECT count(*)
		FROM matches
		WHERE match_type = $1
		AND winner_id IS NOT NULL;
	`, matchType).Scan(&count)
	if err != nil {
		return
	}

	prepared = count > 0 // matches started
	return
}

func GeneratePlayoffMatchesScores(ctx context.Context, matches []*dto.MatchDTO) (err error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second) // here give more time
	defer cancel()

	tx, err := conn.BeginTx(timeoutCtx, pgx.TxOptions{}) // Start the transaction
	if err != nil {
		return
	}

	defer func() { // rollback on error
		if err != nil {
			tx.Rollback(timeoutCtx)
		}
	}()

	for _, m := range matches {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		default:
			queryCtx, queryCancel := context.WithTimeout(timeoutCtx, 3*time.Second) // separate timeout for query
			defer queryCancel()

			_, err = tx.Exec(queryCtx, `
				UPDATE matches
				SET first_team_score = $1,
					second_team_score = $2,
					winner_id = $3,
					played = true
				WHERE id = $4`,
				m.FirstTeamScore,
				m.SecondTeamScore,
				m.WinnerId,
				m.Id,
			)
			if err != nil {
				return
			}
		}
	}

	// Commit the transaction
	if err = tx.Commit(timeoutCtx); err != nil {
		return
	}

	return
}
