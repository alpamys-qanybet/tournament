package ctrl

import (
	"context"
	"log"
	"testing"
	"tournament/dto"
	"tournament/internal/model"
	"tournament/pkg/db"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

var conn *pgxpool.Pool
var ctx context.Context = context.Background()

func TestPreFn(t *testing.T) {
	databaseUrl := "postgresql://postgres:postgres@localhost:5433/tournament" // lets assume you are using docker database
	var err error

	conn, err = db.Connect(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("Error on postgres database: %v\n", err)
	}
}

func TestDivisionRandomization(t *testing.T) {
	err := model.Cleanup(ctx)
	if err != nil {
		return
	}

	err = model.GenerateTeams(ctx)
	if err != nil {
		return
	}

	teams, err := model.GetTeamList(ctx)
	if err != nil {
		return
	}

	originalTeams := make([]*dto.TeamDTO, len(teams))
	_ = copy(originalTeams, teams)

	n := 8
	oTeamsA, oTeamsB := originalTeams[:n], originalTeams[n:]

	teamsA, teamsB := divideTeamsIntoDivisions(teams)

	// check randomization for identicality
	identicalCount := 0
	for i, v := range oTeamsA {
		if teamsA[i].Id == v.Id {
			identicalCount++
		}
	}

	for i, v := range oTeamsB {
		if teamsB[i].Id == v.Id {
			identicalCount++
		}
	}
	assert.NotEqual(t, identicalCount, 16) // not 100%

	// check division boxes are not 10 collectively similar just with random indexes
	teamMap := make(map[uint16]int) // id -> similary

	for _, t := range oTeamsA { // checking only one box(A or B) is enough to check 10 teams similarity regardles indexes
		teamMap[t.Id] = 1
	}

	for _, t := range teamsA {
		if _, ok := teamMap[t.Id]; ok {
			teamMap[t.Id] = 2
		} else {
			teamMap[t.Id] = 3
		}
	}

	// 1 - only in original, 2 - both original and randomized, 3 - only in randomized
	bothSimilarCount := 0
	for _, s := range teamMap {
		// fmt.Println(s)
		if s == 2 {
			bothSimilarCount++
		}
	}

	assert.NotEqual(t, identicalCount, 8) // not 100%

	matchesA, err := generateMatches(ctx, teamsA, dto.DivisionA)
	if err != nil {
		return
	}
	assert.Equal(t, len(matchesA), 28) // number of matches for each division must be 28

	err = generateMatchesScores(ctx, matchesA, teamsA)
	if err != nil {
		return
	}

	// every team must play exactly 7 matches by wins, draws, loses
	for _, team := range teamsA {
		assert.Equal(t, int(team.Wins+team.Draws+team.Loses), 7)
	}
}

func TestPostFn(t *testing.T) {
	defer conn.Close()
}
