package model

import (
	"context"
	"database/sql"
	"time"
	"tournament/dto"
	"tournament/pkg/db"
	"tournament/pkg/helper"
)

// table matches:
// id
// first_team_score
// second_team_score
// winner_id
// match_type
// played
// first_team_id
// second_team_id

func GetDivisionMatchList(ctx context.Context, matchType string) ([]*dto.MatchDTO, error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return nil, err
	}

	result := []*dto.MatchDTO{}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := conn.Query(timeoutCtx, `
		SELECT 
			m.id, 
			CONCAT(t1.name, ' - ', t2.name) AS name,
			m.first_team_score,
			m.second_team_score,
			m.first_team_id, 
			m.second_team_id
		FROM matches m
		LEFT JOIN teams t1 ON t1.id = m.first_team_id
		LEFT JOIN teams t2 ON t2.id = m.second_team_id
		WHERE m.match_type = $1
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
			var item dto.MatchDTO

			var firstTeamScore, secondTeamScore sql.NullInt16

			err = rows.Scan(&item.Id, &item.Name, &firstTeamScore, &secondTeamScore, &item.FirstTeamId, &item.SecondTeamId)
			if err != nil {
				return nil, err
			}

			if firstTeamScore.Valid {
				item.FirstTeamScore = helper.ToPtrUint16(uint16(firstTeamScore.Int16))
			}

			if secondTeamScore.Valid {
				item.SecondTeamScore = helper.ToPtrUint16(uint16(secondTeamScore.Int16))
			}

			result = append(result, &item)
		}
	}
	return result, rows.Err()
}

func GetPlayoffMatchList(ctx context.Context, matchType string) ([]*dto.MatchDTO, error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return nil, err
	}

	result := []*dto.MatchDTO{}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := conn.Query(timeoutCtx, `
		SELECT 
			m.id, 
			CONCAT(t1.name, ' - ', t2.name) AS name,
			m.first_team_score,
			m.second_team_score,
			m.first_team_id, 
			m.second_team_id,
			m.winner_id,
			m.match_type
		FROM matches m
		LEFT JOIN teams t1 ON t1.id = m.first_team_id
		LEFT JOIN teams t2 ON t2.id = m.second_team_id
		WHERE m.match_type = $1
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
			var item dto.MatchDTO

			var firstTeamScore, secondTeamScore sql.NullInt16

			err = rows.Scan(&item.Id, &item.Name, &firstTeamScore, &secondTeamScore, &item.FirstTeamId, &item.SecondTeamId, &item.WinnerId, &item.MatchType)
			if err != nil {
				return nil, err
			}

			if firstTeamScore.Valid {
				item.FirstTeamScore = helper.ToPtrUint16(uint16(firstTeamScore.Int16))
			}

			if secondTeamScore.Valid {
				item.SecondTeamScore = helper.ToPtrUint16(uint16(secondTeamScore.Int16))
			}

			result = append(result, &item)
		}
	}
	return result, rows.Err()
}
