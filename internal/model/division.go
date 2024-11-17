package model

import (
	"context"
	"time"
	"tournament/dto"
	"tournament/pkg/db"

	"github.com/jackc/pgx/v4"
)

func DivisionIsPrepared(ctx context.Context, division string) (prepared bool, err error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count uint16
	err = conn.QueryRow(timeoutCtx, `
		SELECT count(*)
		FROM teams
		WHERE division = $1;
	`, division).Scan(&count)
	if err != nil {
		return
	}

	prepared = count > 0 // teams prepared, so matches
	return
}

func PrepareDivisions(ctx context.Context, teamsA, teamsB []*dto.TeamDTO, matchesA, matchesB []*dto.MatchDTO) (err error) {
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

	// we could use single sql query with all commands, but multi query statements does not support $1, $2 passing <=== so code will be UGLY with string operations
	// our database is small, so it is trade-off.
	// in other cases should use string operations to run all queries in BEGIN; queries....; COMMIT;

	// again, our database is small, and using one column as STR with indexes is trade-off instead of using FK extra table.
	// in cases of infinite division and infinite tournament need to restructure database and re-write code to use FK values instead of STR values

	// set teams into division A, just put A or B, extra table division(id, name) <---> team(id, ...., division_id) <==== is not needed here for small database and project

	{
		queryCtx, queryCancel := context.WithTimeout(timeoutCtx, 3*time.Second) // separate timeout for query
		defer queryCancel()

		_, err = tx.Exec(queryCtx, `
			UPDATE teams SET division = $1
			WHERE id IN (`+prepareTeamInIdsStr(teamsA)+`)
		`, dto.DivisionA)
		if err != nil {
			return
		}
	}

	{
		queryCtx, queryCancel := context.WithTimeout(timeoutCtx, 3*time.Second) // separate timeout for query
		defer queryCancel()

		// set teams into division B
		_, err = tx.Exec(queryCtx, `
			UPDATE teams SET division = $1
			WHERE id IN (`+prepareTeamInIdsStr(teamsB)+`)
		`, dto.DivisionB)
		if err != nil {
			return
		}
	}

	// just put DA, DB, P4, P2, PF instead of match_type(id, name) <----> team(id, ...., division_id, match_type_id) <----> division(id, name) or maybe more complex relations
	// our database and project are small, so it's trade-off
	for _, m := range matchesA {
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

	for _, m := range matchesB {
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

func DivisionIsStarted(ctx context.Context) (started bool, err error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count uint16
	err = conn.QueryRow(timeoutCtx, `
		SELECT COALESCE(SUM(wins + draws + loses), 0) AS played
		FROM teams;
	`).Scan(&count)
	if err != nil {
		return
	}

	started = count > 0 // matches played, so started
	return
}

func StartDivisions(ctx context.Context, teamsA, teamsB []*dto.TeamDTO, matchesA, matchesB []*dto.MatchDTO) (err error) {
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

	for _, t := range teamsA {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		default:
			queryCtx, queryCancel := context.WithTimeout(timeoutCtx, 3*time.Second) // separate timeout for query
			defer queryCancel()

			_, err = tx.Exec(queryCtx, `
				UPDATE teams
				SET wins = $1,
					draws = $2,
					loses = $3,
					goals_scored = $4,
					goals_conceded = $5,
					goal_diff = $6,
					points = $7
				WHERE id = $8`,
				t.Wins,
				t.Draws,
				t.Loses,
				t.GoalsScored,
				t.GoalsConceded,
				t.GoalDiff,
				t.Points,
				t.Id,
			)
			if err != nil {
				return
			}
		}
	}

	for _, t := range teamsB {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		default:
			queryCtx, queryCancel := context.WithTimeout(timeoutCtx, 3*time.Second) // separate timeout for query
			defer queryCancel()

			_, err = tx.Exec(queryCtx, `
				UPDATE teams
				SET wins = $1,
					draws = $2,
					loses = $3,
					goals_scored = $4,
					goals_conceded = $5,
					goal_diff = $6,
					points = $7
				WHERE id = $8`,
				t.Wins,
				t.Draws,
				t.Loses,
				t.GoalsScored,
				t.GoalsConceded,
				t.GoalDiff,
				t.Points,
				t.Id,
			)
			if err != nil {
				return
			}
		}
	}

	for _, m := range matchesA {
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

	for _, m := range matchesB {
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
