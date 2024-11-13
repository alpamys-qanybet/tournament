package main

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

func main() {
	// fmt.Println("Division")

	// divideTeamsIntoDivisions()
	generateMatches()
}

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

func divideTeamsIntoDivisions() {
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

	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	r.Shuffle(len(teams), func(i, j int) { // shuffle teams(re-order with random indexes)
		teams[i], teams[j] = teams[j], teams[i]
	})

	fmt.Println("\nre-ordered teams:")
	for i, v := range teams {
		fmt.Println(i+1, fmt.Sprintf("id: %d, name: %s", v.id, v.name))
	}

	n := 8
	divisionATeams := make([]*Team, 0)
	divisionATeams = append(divisionATeams, teams[:n]...)
	divisionBTeams := teams[n:]

	fmt.Println("\ndivision A teams:")
	for i, v := range divisionATeams {
		fmt.Println(i+1, fmt.Sprintf("id: %d, name: %s", v.id, v.name))
	}

	fmt.Println("\ndivision B teams:")
	for i, v := range divisionBTeams {
		fmt.Println(i+1, fmt.Sprintf("id: %d, name: %s", v.id, v.name))
	}
}

func generateMatches() { // assume taking 8 teams
	var teams = []*Team{
		{1, "Liverpool", 0, 0, 0, 0, 0, 0},
		{2, "Arsenal", 0, 0, 0, 0, 0, 0},
		{3, "Aston Villa", 0, 0, 0, 0, 0, 0},
		{4, "Milan", 0, 0, 0, 0, 0, 0},
		{5, "Juventus", 0, 0, 0, 0, 0, 0},
		{6, "Barcelona", 0, 0, 0, 0, 0, 0},
		{7, "Bayern Munchen", 0, 0, 0, 0, 0, 0},
		{8, "Borussia Dortmund", 0, 0, 0, 0, 0, 0},
	}

	fmt.Println("teams:")
	for i, v := range teams {
		fmt.Println(i+1, fmt.Sprintf("id: %d, name: %s", v.id, v.name))
	}

	matches := make([]*Match, 0)
	m := make(map[string]bool)

	id := 1
	for i, v := range teams {
		for j, v2 := range teams {
			if i == j {
				continue
			}

			key, keyReverse := fmt.Sprintf("%d_%d", i, j), fmt.Sprintf("%d_%d", j, i)
			_, ok := m[key]
			_, ok2 := m[keyReverse]

			if !(ok || ok2) { // none exists
				matches = append(matches, &Match{
					id,
					v,
					v2,
					0,
					0,
					-1,
				})
				m[key] = true
				id++
			}
		}
	}

	fmt.Println("\nmatches:")
	for i, v := range matches {
		fmt.Println(i+1, v.id, v.name())
	}
	// 28 matches // correct

	scheduledMatches := make([]*Match, 0)

	rand.Seed(uint64(time.Now().Unix()))

	loopIteration := 0
	fmt.Println()
	for {
		loopIteration++
		if len(matches) == 0 {
			break
		}

		randomIndex := rand.Intn(len(matches))
		match := matches[randomIndex]

		if loopIteration > 50 { // sometimes it is hard to schedule to have relaxing teams
			fmt.Println("re-schedule", len(matches))
			// re-schedule matches, insert matches at some indexes
			for {
				if len(matches) == 0 {
					break
				}
				v := matches[0]

				for j := 0; j < len(scheduledMatches)-1; j++ {

					if j == 0 { // first
						if scheduledMatches[j].firstTeam.id == v.firstTeam.id || scheduledMatches[j].firstTeam.id == v.secondTeam.id {
							continue
						}

						if scheduledMatches[j].secondTeam.id == v.firstTeam.id || scheduledMatches[j].secondTeam.id == v.secondTeam.id {
							continue
						}
					} else { // 1 ... n

						if scheduledMatches[j-1].firstTeam.id == v.firstTeam.id || scheduledMatches[j-1].firstTeam.id == v.secondTeam.id {
							continue
						}

						if scheduledMatches[j-1].secondTeam.id == v.firstTeam.id || scheduledMatches[j-1].secondTeam.id == v.secondTeam.id {
							continue
						}

						if scheduledMatches[j].firstTeam.id == v.firstTeam.id || scheduledMatches[j].firstTeam.id == v.secondTeam.id {
							continue
						}

						if scheduledMatches[j].secondTeam.id == v.firstTeam.id || scheduledMatches[j].secondTeam.id == v.secondTeam.id {
							continue
						}
					}

					matches = append(matches[:0], matches[1:]...)
					scheduledMatches = append(scheduledMatches[:j+1], scheduledMatches[j:]...)
					scheduledMatches[j] = v
					break
				}
			}

			break
		}

		n := len(scheduledMatches)
		if n != 0 {

			if scheduledMatches[n-1].firstTeam.id == match.firstTeam.id || scheduledMatches[n-1].firstTeam.id == match.secondTeam.id {
				continue
			}

			if scheduledMatches[n-1].secondTeam.id == match.firstTeam.id || scheduledMatches[n-1].secondTeam.id == match.secondTeam.id {
				continue
			}
		}
		matches = append(matches[:randomIndex], matches[randomIndex+1:]...)
		scheduledMatches = append(scheduledMatches, match)
	}

	fmt.Println("\nscheduled matches:")
	for i, v := range scheduledMatches {
		fmt.Println(i+1, v.id, v.name()) // , v.score()
	}

	// generate scores of matches
	maxGoals := 5 // assume max goals
	for _, v := range scheduledMatches {
		v.firstTeamScore = rand.Intn(maxGoals)
		v.secondTeamScore = rand.Intn(maxGoals)

		if v.firstTeamScore > v.secondTeamScore {
			v.winnerId = v.firstTeam.id
			v.firstTeam.points += 3
			v.firstTeam.wins++
			v.secondTeam.loses++
		} else if v.firstTeamScore < v.secondTeamScore {
			v.winnerId = v.secondTeam.id
			v.secondTeam.points += 3
			v.secondTeam.wins++
			v.firstTeam.loses++
		} else {
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

	fmt.Println("\nscores of matches:")
	for i, v := range scheduledMatches {
		fmt.Println(i+1, v.id, v.name(), v.score())
	}

	fmt.Println("\nteam points:")
	for i, v := range teams {
		fmt.Println(i+1, fmt.Sprintf("id %d, name %s, W %d, D %d, L %d, F %d, A %d, GD %d, P %d", v.id, v.name, v.wins, v.draws, v.loses, v.goalsScored, v.goalsConceded, v.diff(), v.points))
	}

	// leave sorting algorithm to db, points, then diff, then if you have same teams look at the matches between???
}
