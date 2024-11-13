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

	// schedule matches
}
