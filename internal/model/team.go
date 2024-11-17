package model

import (
	"context"
	"errors"
	"time"
	"tournament/dto"
	"tournament/pkg/db"
)

// table teams:
// id
// name
// division
// wins
// draws
// loses
// goals_scored
// goals_conceded
// goal_diff
// points

func GetTeamList(ctx context.Context) ([]*dto.TeamDTO, error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return nil, err
	}

	result := []*dto.TeamDTO{}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second) // if it takes longer than 5s it is too long
	defer cancel()                                                // parent ctx will be ok

	rows, err := conn.Query(timeoutCtx, `
		SELECT id, name
		FROM teams
		ORDER BY id ASC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			var item dto.TeamDTO

			err = rows.Scan(&item.Id, &item.Name)
			if err != nil {
				return nil, err
			}
			result = append(result, &item)
		}
	}
	return result, rows.Err()
}

func GetTeamListByDivision(ctx context.Context, division string, top4 bool) ([]*dto.TeamDTO, error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return nil, err
	}

	result := []*dto.TeamDTO{}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sqlQuery := `
		SELECT id, name, wins, draws, loses, goals_scored, goals_conceded, goal_diff, points
		FROM teams
		WHERE division = $1
	`

	if top4 {
		sqlQuery += " AND points > 0"
	}

	sqlQuery += " ORDER BY points DESC, goal_diff DESC, goals_scored DESC, name ASC"

	if top4 {
		sqlQuery += " LIMIT 4"
	}

	sqlQuery += ";"

	rows, err := conn.Query(timeoutCtx, sqlQuery, division)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			var item dto.TeamDTO

			err = rows.Scan(&item.Id, &item.Name, &item.Wins, &item.Draws, &item.Loses, &item.GoalsScored, &item.GoalsConceded, &item.GoalDiff, &item.Points)
			if err != nil {
				return nil, err
			}

			item.Division = division

			result = append(result, &item)
		}
	}
	return result, rows.Err()
}

func GetTeamListByPlayoffWinners(ctx context.Context, matchType string) ([]*dto.TeamDTO, error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return nil, err
	}

	result := []*dto.TeamDTO{}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := conn.Query(timeoutCtx, `
		SELECT t.id, t.name
		FROM matches m, teams t 
		WHERE m.match_type = $1
		AND m.winner_id = t.id
		ORDER BY m.id ASC;
	`, matchType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			var item dto.TeamDTO

			err = rows.Scan(&item.Id, &item.Name)
			if err != nil {
				return nil, err
			}

			result = append(result, &item)
		}
	}
	return result, rows.Err()
}

func CreateTeam(ctx context.Context, name string) (id uint16, err error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return
	}

	{
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var count uint16
		err = conn.QueryRow(timeoutCtx, `
			SELECT count(*)
			FROM teams;
		`).Scan(&count)
		if err != nil {
			return
		}

		if count >= 16 {
			err = errors.New("max_16_teams_allowed")
			return
		}
	}

	{
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = conn.QueryRow(timeoutCtx, `
			INSERT INTO teams(name)
			VALUES ($1) RETURNING id;`,
			name,
		).Scan(&id)
	}

	return
}

func GetTeamCount(ctx context.Context) (count uint16, err error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = conn.QueryRow(timeoutCtx, `
		SELECT count(*)
		FROM teams;
	`).Scan(&count)
	if err != nil {
		return
	}

	return
}

func GenerateTeams(ctx context.Context) (err error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return
	}

	{
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var count uint16
		err = conn.QueryRow(timeoutCtx, `
			SELECT count(*)
			FROM teams;
		`).Scan(&count)
		if err != nil {
			return
		}

		if count > 0 {
			err = errors.New("generation_only_allowed_into_empty_table")
			return
		}
	}

	{
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_, err = conn.Exec(timeoutCtx, `
			INSERT INTO teams(name)
			VALUES 
			('Liverpool'),
			('Arsenal'),
			('Aston Villa'),
			('Milan'),
			('Juventus'),
			('Barcelona'),
			('Bayern Munchen'),
			('Borussia Dortmund'),
			('Manchester City'),
			('Chelsea'),
			('Manchester United'),
			('Inter milan'),
			('Atalanta'),
			('Real Madrid'),
			('Atletico Madrid'),
			('Bayer Leverkusen');
		`)
	}

	return
}
