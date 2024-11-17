package ctrl

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
	"tournament/dto"
	"tournament/internal/model"
	"tournament/pkg/helper"
)

func GetDivisions(ctx context.Context) (interface{}, error) {
	teamsA, err := model.GetTeamListByDivision(ctx, dto.DivisionA, false)
	if err != nil {
		return nil, err
	}

	teamsB, err := model.GetTeamListByDivision(ctx, dto.DivisionB, false)
	if err != nil {
		return nil, err
	}

	matchesA, err := model.GetDivisionMatchList(ctx, dto.MatchTypeDivisionA)
	if err != nil {
		return nil, err
	}

	matchesB, err := model.GetDivisionMatchList(ctx, dto.MatchTypeDivisionB)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"division_a": &dto.DivisionDTO{
			Name:    dto.DivisionA,
			Teams:   teamsA,
			Matches: matchesA,
		},
		"division_b": &dto.DivisionDTO{
			Name:    dto.DivisionB,
			Teams:   teamsB,
			Matches: matchesB,
		},
	}

	return data, nil
}

func PrepareDivisions(ctx context.Context) error {

	n, err := model.GetTeamCount(ctx)
	if err != nil {
		return err
	}

	if n != 16 {
		err = errors.New("must_have_16_teams_to_prepare_divisions")
		return err
	}

	prepared, err := model.DivisionIsPrepared(ctx, dto.DivisionA)
	if err != nil {
		return err
	}

	if prepared {
		err = errors.New("division_is_already_prepared")
		return err
	}

	teams, err := model.GetTeamList(ctx)
	if err != nil {
		return err
	}

	teamsA, teamsB := divideTeamsIntoDivisions(teams)

	var wg sync.WaitGroup
	wg.Add(2)

	errCh := make(chan error, 2) // Buffered channel to hold errors from both goroutines

	var matchesA, matchesB []*dto.MatchDTO

	go func(ctx context.Context) {
		defer wg.Done()
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
		default:
			matchesA, err = generateMatches(ctx, teamsA, dto.DivisionA)
			if err != nil {
				errCh <- err
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		defer wg.Done()
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
		default:
			matchesB, err = generateMatches(ctx, teamsB, dto.DivisionB)
			if err != nil {
				errCh <- err
			}
		}
	}(ctx)

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for e := range errCh {
		if e != nil {
			return e
		}
	}

	return model.PrepareDivisions(ctx, teamsA, teamsB, matchesA, matchesB) // all in one transaction
}

func divideTeamsIntoDivisions(teams []*dto.TeamDTO) ([]*dto.TeamDTO, []*dto.TeamDTO) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(teams), func(i, j int) { // shuffle teams(re-order with random indexes)
		teams[i], teams[j] = teams[j], teams[i]
	})

	n := 8
	for i, v := range teams {
		if i < n {
			v.Division = dto.DivisionA
		} else {
			v.Division = dto.DivisionB
		}
	}

	return teams[:n], teams[n:]
}

func generateMatches(ctx context.Context, teams []*dto.TeamDTO, division string) ([]*dto.MatchDTO, error) {
	matches := make([]*dto.MatchDTO, 0) // generate matches each team plays each other once
	m := make(map[string]bool)

	for i, v := range teams {
		for j, v2 := range teams {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:

				if i == j { // team itself
					continue
				}

				key, keyReverse := fmt.Sprintf("%d_%d", i, j), fmt.Sprintf("%d_%d", j, i) // key and reverse key to meet only play once condition
				_, ok := m[key]
				_, ok2 := m[keyReverse]

				if ok || ok2 { // match exists
					continue
				}

				matchType := dto.MatchTypeDivisionA
				if division == dto.DivisionB {
					matchType = dto.MatchTypeDivisionB
				}

				matches = append(matches, &dto.MatchDTO{
					Name:         fmt.Sprintf("%s - %s", v.Name, v2.Name), // match title
					FirstTeamId:  v.Id,                                    // team 1
					SecondTeamId: v2.Id,                                   // team 2
					MatchType:    matchType,                               // division match a or b
				})
				m[key] = true // mark the key(team1_team2 or team2_team1 ids)
			}
		}
	}

	scheduledMatches := make([]*dto.MatchDTO, 0) // randomly schedule matches

	rand.Seed(time.Now().UnixNano())
	loopIteration := 0

	// regardless of infinite for loops the loop iteration is between 28(number of matches) - 50 + 28 * 30% => 28 - 60,65
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			loopIteration++
			if len(matches) == 0 { // all matches are picked and scheduled, there is no matches left
				return scheduledMatches, nil
			}

			randomIndex := rand.Intn(len(matches)) // pick random match
			match := matches[randomIndex]

			if loopIteration > 50 { // sometimes it is hard to schedule to have relaxing teams within short loop iteration
				// re-schedule matches, insert matches at some indexes instead of picking randomly to match against the last,
				// the scheduled order cannot satisfy this logic, so instead re-schedule at some point
				for {
					if len(matches) == 0 { // all left matches that could not fit the schedule by random picking are scheduled by re-scheduling,
						// no matches left, so exit loop, and exit the upper loop, see below
						break
					}
					v := matches[0] // pick first match from left matches

					for j := 0; j < len(scheduledMatches)-1; j++ { // insert at specific index which meets condition of relaxing team

						if j > 0 { // first match has no previous match
							if scheduledMatches[j-1].FirstTeamId == v.FirstTeamId || scheduledMatches[j-1].FirstTeamId == v.SecondTeamId {
								continue
							}

							if scheduledMatches[j-1].SecondTeamId == v.FirstTeamId || scheduledMatches[j-1].SecondTeamId == v.SecondTeamId {
								continue
							}
						}

						// compare with given index to insert at that index
						if scheduledMatches[j].FirstTeamId == v.FirstTeamId || scheduledMatches[j].FirstTeamId == v.SecondTeamId {
							continue
						}

						if scheduledMatches[j].SecondTeamId == v.FirstTeamId || scheduledMatches[j].SecondTeamId == v.SecondTeamId {
							continue
						}

						matches = append(matches[:0], matches[1:]...)                              // remove first element and append to appropriate index which meets the condition of relaxing
						scheduledMatches = append(scheduledMatches[:j+1], scheduledMatches[j:]...) // move elements to the right by 1
						scheduledMatches[j] = v                                                    // to append match at given index
						break                                                                      // condition met, exit the loop, go to the next left match
					}
				}
				break // and exit the upper loop, no matches left, all matches scheduled
			}

			n := len(scheduledMatches)
			if n != 0 { // just append first element, and compare 1+ to previous match teams to have the condition of teams to be relaxing, not playing in a raw.

				if scheduledMatches[n-1].FirstTeamId == match.FirstTeamId || scheduledMatches[n-1].FirstTeamId == match.SecondTeamId {
					continue
				}

				if scheduledMatches[n-1].SecondTeamId == match.FirstTeamId || scheduledMatches[n-1].SecondTeamId == match.SecondTeamId {
					continue
				}
			}

			matches = append(matches[:randomIndex], matches[randomIndex+1:]...) // pick from matches and add to scheduled, that is done to pick matches one by one further from matches
			scheduledMatches = append(scheduledMatches, match)
		}
	}
}

func StartDivisions(ctx context.Context) error {

	n, err := model.GetTeamCount(ctx)
	if err != nil {
		return err
	}

	if n != 16 {
		err = errors.New("must_have_16_teams_to_start_divisions")
		return err
	}

	prepared, err := model.DivisionIsPrepared(ctx, dto.DivisionA)
	if err != nil {
		return err
	}

	if !prepared {
		err = errors.New("division_is_not_prepared")
		return err
	}

	started, err := model.DivisionIsStarted(ctx)
	if err != nil {
		return err
	}

	if started {
		err = errors.New("division_is_already_started")
		return err
	}

	teamsA, err := model.GetTeamListByDivision(ctx, dto.DivisionA, false)
	if err != nil {
		return err
	}

	teamsB, err := model.GetTeamListByDivision(ctx, dto.DivisionB, false)
	if err != nil {
		return err
	}

	matchesA, err := model.GetDivisionMatchList(ctx, dto.MatchTypeDivisionA)
	if err != nil {
		return err
	}

	matchesB, err := model.GetDivisionMatchList(ctx, dto.MatchTypeDivisionB)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	errCh := make(chan error, 2) // Buffered channel to hold errors from both goroutines

	go func(ctx context.Context) {
		defer wg.Done()
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
		default:
			err = generateMatchesScores(ctx, matchesA, teamsA)
			if err != nil {
				errCh <- err
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		defer wg.Done()
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
		default:
			err = generateMatchesScores(ctx, matchesB, teamsB)
			if err != nil {
				errCh <- err
			}
		}
	}(ctx)

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for e := range errCh {
		if e != nil {
			return e
		}
	}

	return model.StartDivisions(ctx, teamsA, teamsB, matchesA, matchesB)
}

func generateMatchesScores(ctx context.Context, matches []*dto.MatchDTO, teams []*dto.TeamDTO) error {
	maxGoals := 5 // assume max goals
	rand.Seed(time.Now().UnixNano())
	teamMap := make(map[uint16]int) // id -> index

	for i, t := range teams {
		if _, ok := teamMap[t.Id]; !ok {
			teamMap[t.Id] = i
		}
	}

	for _, m := range matches {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m.FirstTeamScore = helper.ToPtrUint16(uint16(rand.Intn(maxGoals)))
			m.SecondTeamScore = helper.ToPtrUint16(uint16(rand.Intn(maxGoals)))

			firstTeamIndex := teamMap[m.FirstTeamId]
			SecondTeamIndex := teamMap[m.SecondTeamId]

			if *m.FirstTeamScore > *m.SecondTeamScore { // first team wins
				m.WinnerId = helper.ToPtrInt16(int16(m.FirstTeamId))
				teams[firstTeamIndex].Points += 3
				teams[firstTeamIndex].Wins++
				teams[SecondTeamIndex].Loses++
			} else if *m.FirstTeamScore < *m.SecondTeamScore { // second team wins
				m.WinnerId = helper.ToPtrInt16(int16(m.SecondTeamId))
				teams[SecondTeamIndex].Points += 3
				teams[SecondTeamIndex].Wins++
				teams[firstTeamIndex].Loses++
			} else { // draw
				m.WinnerId = helper.ToPtrInt16(int16(-1))
				teams[firstTeamIndex].Points++
				teams[SecondTeamIndex].Points++
				teams[firstTeamIndex].Draws++
				teams[SecondTeamIndex].Draws++
			}

			teams[firstTeamIndex].GoalsScored += *m.FirstTeamScore
			teams[SecondTeamIndex].GoalsScored += *m.SecondTeamScore

			teams[firstTeamIndex].GoalsConceded += *m.SecondTeamScore
			teams[SecondTeamIndex].GoalsConceded += *m.FirstTeamScore
		}
	}

	for _, t := range teams {
		t.SetDiff()
	}

	return nil
}
