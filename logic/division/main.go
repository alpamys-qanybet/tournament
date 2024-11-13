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
	id   int
	name string
}

type Match struct {
	id         int
	firstTeam  *Team
	secondTeam *Team
}

func (m *Match) name() string {
	return fmt.Sprintf("%s - %s", m.firstTeam.name, m.secondTeam.name)
}

func divideTeamsIntoDivisions() {
	var teams = []*Team{
		{1, "Liverpool"},
		{2, "Arsenal"},
		{3, "Aston Villa"},
		{4, "Milan"},
		{5, "Juventus"},
		{6, "Barcelona"},
		{7, "Bayern Munchen"},
		{8, "Borussia Dortmund"},
		{9, "Manchester City"},
		{10, "Chelsea"},
		{11, "Manchester United"},
		{12, "Inter milan"},
		{13, "Atalanta"},
		{14, "Real Madrid"},
		{15, "Atletico Madrid"},
		{16, "Bayer Leverkusen"},
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
		{1, "Liverpool"},
		{2, "Arsenal"},
		{3, "Aston Villa"},
		{4, "Milan"},
		{5, "Juventus"},
		{6, "Barcelona"},
		{7, "Bayern Munchen"},
		{8, "Borussia Dortmund"},
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
		fmt.Println(i+1, v.id, v.name())
	}
}
