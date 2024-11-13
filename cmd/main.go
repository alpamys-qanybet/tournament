package main

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

type Team struct {
	id            int
	name          string
	points        int
	goalsScored   int
	goalsConceded int
	wins          int
	draws         int
	loses         int
}

func (t *Team) diff() int {
	return t.goalsScored - t.goalsConceded
}

type Match struct {
	id              int
	firstTeam       *Team
	secondTeam      *Team
	firstTeamScore  int
	secondTeamScore int
	winnerId        int
}

func (m *Match) name() string {
	return fmt.Sprintf("%s - %s", m.firstTeam.name, m.secondTeam.name)
}

func (m *Match) score() string {
	return fmt.Sprintf("%d - %d", m.firstTeamScore, m.secondTeamScore)
}

func main() {
	// fmt.Println("Division")
	var teams = []*Team{
		{1, "Liverpool", 0, 0, 0, 0, 0, 0},
		{2, "Arsenal", 0, 0, 0, 0, 0, 0},
		{3, "Aston Villa", 0, 0, 0, 0, 0, 0},
		{4, "Milan", 0, 0, 0, 0, 0, 0},
		{5, "Juventus", 0, 0, 0, 0, 0, 0},
		{6, "Barcelona", 0, 0, 0, 0, 0, 0},
		{7, "Bayern Munchen", 0, 0, 0, 0, 0, 0},
		{8, "Borussia Dortmund", 0, 0, 0, 0, 0, 0},
		{9, "Manchester City", 0, 0, 0, 0, 0, 0},
		{10, "Chelsea", 0, 0, 0, 0, 0, 0},
		{11, "Manchester United", 0, 0, 0, 0, 0, 0},
		{12, "Inter milan", 0, 0, 0, 0, 0, 0},
		{13, "Atalanta", 0, 0, 0, 0, 0, 0},
		{14, "Real Madrid", 0, 0, 0, 0, 0, 0},
		{15, "Atletico Madrid", 0, 0, 0, 0, 0, 0},
		{16, "Bayer Leverkusen", 0, 0, 0, 0, 0, 0},
	}

	fmt.Println("teams:")
	for i, v := range teams {
		fmt.Println(i+1, fmt.Sprintf("id: %d, name: %s", v.id, v.name))
	}

	divisionATeams, divisionBTeams := divideTeamsIntoDivisions(teams)

	fmt.Println("\ndivision A teams:")
	for i, v := range divisionATeams {
		fmt.Println(i+1, fmt.Sprintf("id: %d, name: %s", v.id, v.name))
	}

	fmt.Println("\ndivision B teams:")
	for i, v := range divisionBTeams {
		fmt.Println(i+1, fmt.Sprintf("id: %d, name: %s", v.id, v.name))
	}

	matchesA := generateMatches(divisionATeams, 1)
	fmt.Println("\ndivision A matches:")
	for i, v := range matchesA {
		fmt.Println(i+1, v.id, v.name(), v.score())
	}

	matchesB := generateMatches(divisionBTeams, len(matchesA)+1)
	fmt.Println("\ndivision B matches:")
	for i, v := range matchesB {
		fmt.Println(i+1, v.id, v.name(), v.score())
	}

	generateMatchesScores(matchesA)
	fmt.Println("\ndivision A scores of matches:")
	for i, v := range matchesA {
		fmt.Println(i+1, v.id, v.name(), v.score())
	}

	generateMatchesScores(matchesB)
	fmt.Println("\ndivision B scores of matches:")
	for i, v := range matchesB {
		fmt.Println(i+1, v.id, v.name(), v.score())
	}

	fmt.Println("\ndivision A team points:")
	for i, v := range divisionATeams {
		fmt.Println(i+1, fmt.Sprintf("id %d, name %s, W %d, D %d, L %d, F %d, A %d, GD %d, P %d", v.id, v.name, v.wins, v.draws, v.loses, v.goalsScored, v.goalsConceded, v.diff(), v.points))
	}

	fmt.Println("\ndivision B team points:")
	for i, v := range divisionBTeams {
		fmt.Println(i+1, fmt.Sprintf("id %d, name %s, W %d, D %d, L %d, F %d, A %d, GD %d, P %d", v.id, v.name, v.wins, v.draws, v.loses, v.goalsScored, v.goalsConceded, v.diff(), v.points))
	}

	// assuming having ranking tables by db sorting

	// 5 id 5, name Juventus, W 5, D 1, L 1, F 20, A 14, GD 6, P 16
	// 6 id 4, name Milan, W 4, D 2, L 1, F 20, A 16, GD 4, P 14
	// 4 id 1, name Liverpool, W 4, D 1, L 2, F 16, A 13, GD 3, P 13
	// 3 id 2, name Arsenal, W 3, D 1, L 3, F 18, A 15, GD 3, P 10
	// 2 id 11, name Manchester United, W 2, D 2, L 3, F 17, A 14, GD 3, P 8
	// 7 id 6, name Barcelona, W 2, D 1, L 4, F 13, A 17, GD -4, P 7
	// 8 id 9, name Manchester City, W 1, D 2, L 4, F 17, A 24, GD -7, P 5
	// 1 id 16, name Bayer Leverkusen, W 1, D 2, L 4, F 14, A 22, GD -8, P 5

	// 8 id 7, name Bayern Munchen, W 6, D 1, L 0, F 22, A 12, GD 10, P 19
	// 5 id 10, name Chelsea, W 4, D 2, L 1, F 22, A 14, GD 8, P 14
	// 7 id 14, name Real Madrid, W 4, D 1, L 2, F 13, A 12, GD 1, P 13
	// 2 id 15, name Atletico Madrid, W 3, D 2, L 2, F 19, A 16, GD 3, P 11
	// 4 id 12, name Inter milan, W 3, D 0, L 4, F 16, A 16, GD 0, P 9
	// 1 id 8, name Borussia Dortmund, W 2, D 0, L 5, F 16, A 18, GD -2, P 6
	// 6 id 3, name Aston Villa, W 1, D 2, L 4, F 12, A 18, GD -6, P 5
	// 3 id 13, name Atalanta, W 0, D 2, L 5, F 9, A 23, GD -14, P 2
}

func divideTeamsIntoDivisions(teams []*Team) ([]*Team, []*Team) {
	// rethink of getting randomly ordered from db???
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	r.Shuffle(len(teams), func(i, j int) { // shuffle teams(re-order with random indexes)
		teams[i], teams[j] = teams[j], teams[i]
	})

	fmt.Println("\nrandomly re-ordered teams:")
	for i, v := range teams {
		fmt.Println(i+1, fmt.Sprintf("id: %d, name: %s", v.id, v.name))
	}

	n := 8
	return teams[:n], teams[n:] // capacity has no impact, the len of first division is 8 <=== len, cap : 8,16; 8,8
}

func generateMatches(teams []*Team, id int) []*Match {
	matches := make([]*Match, 0) // generate matches each team plays each other once
	m := make(map[string]bool)

	for i, v := range teams {
		for j, v2 := range teams {
			if i == j { // team itself
				continue
			}

			key, keyReverse := fmt.Sprintf("%d_%d", i, j), fmt.Sprintf("%d_%d", j, i) // key and reverse key to meet only play once condition
			_, ok := m[key]
			_, ok2 := m[keyReverse]

			if ok || ok2 { // match exists
				continue
			}

			matches = append(matches, &Match{
				id, // match id
				v,  // team 1
				v2, // team 2
				0,  // team 1 score
				0,  // team 2 score
				-1, // winner id
			})
			m[key] = true // mark the key(team1_team2 or team2_team1 ids)
			id++

		}
	}

	scheduledMatches := make([]*Match, 0) // randomly schedule matches

	rand.Seed(uint64(time.Now().Unix()))
	loopIteration := 0

	// regardless of infinite for loops the loop iteration is between 28(number of matches) - 50 + 28 * 30% => 28 - 60,65
	for {
		loopIteration++
		if len(matches) == 0 { // all matches are picked and scheduled, there is no matches left
			break // so exit loop
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
						if scheduledMatches[j-1].firstTeam.id == v.firstTeam.id || scheduledMatches[j-1].firstTeam.id == v.secondTeam.id {
							continue
						}

						if scheduledMatches[j-1].secondTeam.id == v.firstTeam.id || scheduledMatches[j-1].secondTeam.id == v.secondTeam.id {
							continue
						}
					}

					// compare with given index to insert at that index
					if scheduledMatches[j].firstTeam.id == v.firstTeam.id || scheduledMatches[j].firstTeam.id == v.secondTeam.id {
						continue
					}

					if scheduledMatches[j].secondTeam.id == v.firstTeam.id || scheduledMatches[j].secondTeam.id == v.secondTeam.id {
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

			if scheduledMatches[n-1].firstTeam.id == match.firstTeam.id || scheduledMatches[n-1].firstTeam.id == match.secondTeam.id {
				continue
			}

			if scheduledMatches[n-1].secondTeam.id == match.firstTeam.id || scheduledMatches[n-1].secondTeam.id == match.secondTeam.id {
				continue
			}
		}

		matches = append(matches[:randomIndex], matches[randomIndex+1:]...) // pick from matches and add to scheduled, that is done to pick matches one by one further from matches
		scheduledMatches = append(scheduledMatches, match)
	}

	return scheduledMatches
}

func generateMatchesScores(matches []*Match) {

	// generate scores of matches
	maxGoals := 5 // assume max goals
	rand.Seed((uint64(time.Now().UnixNano())))
	for _, v := range matches {
		v.firstTeamScore = rand.Intn(maxGoals)
		v.secondTeamScore = rand.Intn(maxGoals)

		if v.firstTeamScore > v.secondTeamScore { // first team wins
			v.winnerId = v.firstTeam.id
			v.firstTeam.points += 3
			v.firstTeam.wins++
			v.secondTeam.loses++
		} else if v.firstTeamScore < v.secondTeamScore { // second team wins
			v.winnerId = v.secondTeam.id
			v.secondTeam.points += 3
			v.secondTeam.wins++
			v.firstTeam.loses++
		} else { // draw
			v.firstTeam.points++
			v.secondTeam.points++
			v.secondTeam.draws++
			v.firstTeam.draws++
		}

		v.firstTeam.goalsScored += v.firstTeamScore
		v.secondTeam.goalsScored += v.secondTeamScore

		v.firstTeam.goalsConceded += v.secondTeamScore
		v.secondTeam.goalsConceded += v.firstTeamScore
	}

	// leave sorting algorithm to db, points, then diff, then if you have same teams look at the matches between???
}
