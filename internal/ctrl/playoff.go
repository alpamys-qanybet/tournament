package ctrl

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"
	"tournament/dto"
	"tournament/internal/model"
	"tournament/pkg/helper"
)

func GetPlayoffs(ctx context.Context) (interface{}, error) {

	teamsA, err := model.GetTeamListByDivision(ctx, dto.DivisionA, true)
	if err != nil {
		return nil, err
	}

	teamsB, err := model.GetTeamListByDivision(ctx, dto.DivisionB, true)
	if err != nil {
		return nil, err
	}

	teamsA = append(teamsA, teamsB...)

	matchesQ, err := model.GetPlayoffMatchList(ctx, dto.MatchTypePlayoffQuarter)
	if err != nil {
		return nil, err
	}

	teamsS, err := model.GetTeamListByPlayoffWinners(ctx, dto.MatchTypePlayoffQuarter)
	if err != nil {
		return nil, err
	}

	matchesS, err := model.GetPlayoffMatchList(ctx, dto.MatchTypePlayoffSemi)
	if err != nil {
		return nil, err
	}

	teamsF, err := model.GetTeamListByPlayoffWinners(ctx, dto.MatchTypePlayoffSemi)
	if err != nil {
		return nil, err
	}

	matchesF, err := model.GetPlayoffMatchList(ctx, dto.MatchTypePlayoffFinal)
	if err != nil {
		return nil, err
	}

	stage, err := getPlayoffStage(ctx)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"quarter": map[string]interface{}{
			"name":    "Quarter-Final",
			"teams":   teamsA,
			"matches": matchesQ,
		},
		"semi": map[string]interface{}{
			"name":    "Semi-Final",
			"teams":   teamsS,
			"matches": matchesS,
		},
		"final": map[string]interface{}{
			"name":    "Final",
			"teams":   teamsF,
			"matches": matchesF,
		},
		"stage": stage,
	}

	return data, nil
}

func getPlayoffStage(ctx context.Context) (stage string, err error) {
	stage = "prepare-quarter"
	// prepare-quarter
	// start-quarter
	// prepare - semi
	// start-semi
	// prepare-final
	// start-final
	// ended

	finalStarted, err := model.PlayoffIsStarted(ctx, dto.MatchTypePlayoffFinal)
	if err != nil {
		return
	}
	if finalStarted { // we have a final winner ===> ended
		stage = "ended"
		return
	}

	finalPrepared, err := model.PlayoffIsPrepared(ctx, dto.MatchTypePlayoffFinal, false)
	if err != nil {
		return
	}
	if finalPrepared { // final match is prepared, so start a final match ===> start-final
		stage = "start-final"
		return
	}

	semiStarted, err := model.PlayoffIsStarted(ctx, dto.MatchTypePlayoffSemi)
	if err != nil {
		return
	}
	if semiStarted { // we have semi winners ===> prepare-final
		stage = "prepare-final"
		return
	}

	semiPrepared, err := model.PlayoffIsPrepared(ctx, dto.MatchTypePlayoffSemi, false)
	if err != nil {
		return
	}
	if semiPrepared { // semi matches are prepared, so start semi matches ===> start-semi
		stage = "start-semi"
		return
	}

	quarterStarted, err := model.PlayoffIsStarted(ctx, dto.MatchTypePlayoffQuarter)
	if err != nil {
		return
	}
	if quarterStarted { // we have quarter winners ===> prepare-semi
		stage = "prepare-semi"
		return
	}

	quarterPrepared, err := model.PlayoffIsPrepared(ctx, dto.MatchTypePlayoffQuarter, false)
	if err != nil {
		return
	}
	if quarterPrepared { // quarter matches are prepared, so start quarter matches ===> start-quarter
		stage = "start-quarter"
		return
	}

	return
}

func PreparePlayoffQuarter(ctx context.Context) (err error) {

	started, err := model.DivisionIsStarted(ctx)
	if err != nil {
		return
	}

	if !started {
		err = errors.New("division_is_not_started")
		return
	}

	prepared, err := model.PlayoffIsPrepared(ctx, dto.MatchTypePlayoffQuarter, true)
	if err != nil {
		return
	}

	if prepared {
		err = errors.New("playoff_quarter_is_already_prepared")
		return
	}

	teamsA, err := model.GetTeamListByDivision(ctx, dto.DivisionA, true)
	if err != nil {
		return
	}

	teamsB, err := model.GetTeamListByDivision(ctx, dto.DivisionB, true)
	if err != nil {
		return
	}

	matches := generateQuarterMatches(teamsA, teamsB) // quarter matches(team1 better than team2)

	err = model.PreparePlayoffMatches(ctx, matches) // all in one transaction

	return
}

// Play-off initial schedule is made by principle - best team plays against worst team.
// 4 vs 4 => 1 vs 4, 2 vs 3, 3 vs 2, 4 vs 1 <=== maybe it will be static logic or with some random
func generateQuarterMatches(teamsA []*dto.TeamDTO, teamsB []*dto.TeamDTO) []*dto.MatchDTO {
	matches := make([]*dto.MatchDTO, 0)

	// main logic:
	// best worst
	// better worse
	// worse better
	// worst best

	// but re-order teams to have nice look:
	// best worst
	// better worse
	// better worse
	// best worst

	// then re-order again to look even better
	// best worst
	// best worst
	// better worse
	// better worse

	// final step: boolean randomly swap 1,2 to 2,1 indexes and 3,4 to 4,3
	// bestA worstB
	// bestB worstA
	// betterA worseB
	// betterB worseA
	// to
	// randomly re-order 1,2
	// bestA worstB
	// bestB worstA
	// or
	// bestB worstA
	// bestA worstB
	// randomly re-order 3,4
	// betterA worseB
	// betterB worseA
	// or
	// betterB worseA
	// betterA worseB

	n := 4
	for i := 0; i < n; i++ {
		// 0 0 division A best team = team 1 <- better
		// 3 4-0-1 division B worst team = team 2

		// 1 1 division A 2nd = team 1 <- better
		// 2 4-1-1 division B 3rd = team 2

		// 2 2 division A 3rd = team 2
		// 1 4-2-1 division B 2nd = team 1 <- better

		// 3 3 division A worst team = team 2
		// 0 4-3-1 division B best team = team 1 <- better

		// so team1 is always better than team2

		team1Index := i
		team2Index := n - i - 1

		var team1, team2 *dto.TeamDTO

		if i < 2 {
			team1 = teamsA[team1Index]
			team2 = teamsB[team2Index]
		} else {
			team1Index, team2Index = team2Index, team1Index // swap, left side is always better(condition best team plays worst team and nice look that shows clearly all left sides are best)

			team1 = teamsB[team1Index]
			team2 = teamsA[team2Index]
		}

		matches = append(matches, &dto.MatchDTO{
			FirstTeamId:  team1.Id,
			SecondTeamId: team2.Id,
			MatchType:    dto.MatchTypePlayoffQuarter, // playoff quarter final 1/4 match
		})
	}

	// now we have
	// best worst
	// better worse
	// better worse
	// best worst

	elementToMove := matches[3]      // Extract the 4th element
	copy(matches[2:4], matches[1:3]) // Shift elements at indices 1 and 2 to indices 2 and 3
	matches[1] = elementToMove       // Place the 4th element at index 1

	// now we have
	// bestA worstB
	// bestB worstA
	// betterA worseB
	// betterB worseA

	// Randomly swap 1st and 2nd matches
	if helper.RandomBool() {
		matches[0], matches[1] = matches[1], matches[0]
	}

	// Randomly swap 3rd and 4th matches
	if helper.RandomBool() {
		matches[2], matches[3] = matches[3], matches[2]
	}

	// now we have regardless of division
	// best worst
	// best worst
	// better worse
	// better worse

	return matches
}

func generatePlayoffMatches(ctx context.Context, teams []*dto.TeamDTO, matchType string) ([]*dto.MatchDTO, error) {
	matches := make([]*dto.MatchDTO, 0)

	if len(teams)%2 == 1 {
		return nil, errors.New("incorrect_number_of_matches_must_be_even")
	}

	rand.Seed(time.Now().UnixNano())
	for { // it is everytime even, so do not worry
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if len(teams) == 0 { // all teams are picked, there is no teams left
				return matches, nil
			}

			randomIndex := rand.Intn(len(teams)) // pick random team
			team1 := teams[randomIndex]
			teams = append(teams[:randomIndex], teams[randomIndex+1:]...) // remove from teams

			randomIndex = rand.Intn(len(teams)) // pick random team
			team2 := teams[randomIndex]
			teams = append(teams[:randomIndex], teams[randomIndex+1:]...) // remove from teams

			matches = append(matches, &dto.MatchDTO{
				FirstTeamId:  team1.Id,
				SecondTeamId: team2.Id,
				MatchType:    matchType,
			})
		}
	}
}

func StartPlayoffQuarter(ctx context.Context) (err error) {

	prepared, err := model.PlayoffIsPrepared(ctx, dto.MatchTypePlayoffQuarter, true)
	if err != nil {
		return
	}

	if !prepared {
		err = errors.New("playoff_quarter_is_not_prepared")
		return
	}

	started, err := model.PlayoffIsStarted(ctx, dto.MatchTypePlayoffQuarter)
	if err != nil {
		return
	}

	if started {
		err = errors.New("playoff_quarter_is_already_started")
		return
	}

	matches, err := model.GetPlayoffMatchList(ctx, dto.MatchTypePlayoffQuarter)
	if err != nil {
		return
	}

	err = generatePlayoffMatchesScores(ctx, matches)
	if err != nil {
		return
	}
	err = model.GeneratePlayoffMatchesScores(ctx, matches)
	return
}

func generatePlayoffMatchesScores(ctx context.Context, matches []*dto.MatchDTO) error {
	var wg sync.WaitGroup
	maxGoals := 5 // assume max goals
	rand.Seed(time.Now().UnixNano())

	errCh := make(chan error, len(matches))

	for _, match := range matches {
		wg.Add(1)
		go func(ctx context.Context, m *dto.MatchDTO) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
			default:
				m.FirstTeamScore = helper.ToPtrUint16(uint16(rand.Intn(maxGoals)))
				m.SecondTeamScore = helper.ToPtrUint16(uint16(rand.Intn(maxGoals)))

				if *m.FirstTeamScore == *m.SecondTeamScore { // draw is not allowed in play-off
					if helper.RandomBool() { // first team wins by random goals ahead
						*m.FirstTeamScore += uint16(rand.Intn(maxGoals) + 1) // must be bigger minimum by 1
					} else { // or second team wins
						*m.SecondTeamScore += uint16(rand.Intn(maxGoals) + 1) // must be bigger minimum by 1
					}
				}

				if *m.FirstTeamScore > *m.SecondTeamScore { // first team wins
					m.WinnerId = helper.ToPtrInt16(int16(m.FirstTeamId))
				} else if *m.SecondTeamScore > *m.FirstTeamScore { // second team wins
					m.WinnerId = helper.ToPtrInt16(int16(m.SecondTeamId))
				}
			}
		}(ctx, match)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for e := range errCh {
		if e != nil {
			return e
		}
	}

	return nil
}

func PreparePlayoffSemi(ctx context.Context) (err error) {

	started, err := model.DivisionIsStarted(ctx)
	if err != nil {
		return
	}

	if !started {
		err = errors.New("division_is_not_started")
		return
	}

	prepared, err := model.PlayoffIsPrepared(ctx, dto.MatchTypePlayoffSemi, true)
	if err != nil {
		return
	}

	if prepared {
		err = errors.New("playoff_semi_is_already_prepared")
		return
	}

	teams, err := model.GetTeamListByPlayoffWinners(ctx, dto.MatchTypePlayoffQuarter)
	if err != nil {
		return
	}

	matches, err := generatePlayoffMatches(ctx, teams, dto.MatchTypePlayoffSemi)
	if err != nil {
		return
	}

	err = model.PreparePlayoffMatches(ctx, matches)
	return
}

func StartPlayoffSemi(ctx context.Context) (err error) {

	prepared, err := model.PlayoffIsPrepared(ctx, dto.MatchTypePlayoffSemi, true)
	if err != nil {
		return
	}

	if !prepared {
		err = errors.New("playoff_semi_is_not_prepared")
		return
	}

	started, err := model.PlayoffIsStarted(ctx, dto.MatchTypePlayoffSemi)
	if err != nil {
		return
	}

	if started {
		err = errors.New("playoff_semi_is_already_started")
		return
	}

	matches, err := model.GetPlayoffMatchList(ctx, dto.MatchTypePlayoffSemi)
	if err != nil {
		return
	}

	err = generatePlayoffMatchesScores(ctx, matches)
	if err != nil {
		return
	}

	err = model.GeneratePlayoffMatchesScores(ctx, matches)
	return
}

func PreparePlayoffFinal(ctx context.Context) (err error) {

	started, err := model.DivisionIsStarted(ctx)
	if err != nil {
		return
	}

	if !started {
		err = errors.New("division_is_not_started")
		return
	}

	prepared, err := model.PlayoffIsPrepared(ctx, dto.MatchTypePlayoffFinal, true)
	if err != nil {
		return
	}

	if prepared {
		err = errors.New("playoff_final_is_already_prepared")
		return
	}

	teams, err := model.GetTeamListByPlayoffWinners(ctx, dto.MatchTypePlayoffSemi)
	if err != nil {
		return
	}

	matches, err := generatePlayoffMatches(ctx, teams, dto.MatchTypePlayoffFinal)
	if err != nil {
		return
	}

	err = model.PreparePlayoffMatches(ctx, matches)
	return
}

func StartPlayoffFinal(ctx context.Context) (err error) {

	prepared, err := model.PlayoffIsPrepared(ctx, dto.MatchTypePlayoffFinal, true)
	if err != nil {
		return
	}

	if !prepared {
		err = errors.New("playoff_final_is_not_prepared")
		return
	}

	started, err := model.PlayoffIsStarted(ctx, dto.MatchTypePlayoffFinal)
	if err != nil {
		return
	}

	if started {
		err = errors.New("playoff_final_is_already_started")
		return
	}

	matches, err := model.GetPlayoffMatchList(ctx, dto.MatchTypePlayoffFinal)
	if err != nil {
		return
	}

	err = generatePlayoffMatchesScores(ctx, matches)
	if err != nil {
		return
	}

	err = model.GeneratePlayoffMatchesScores(ctx, matches)
	return
}
