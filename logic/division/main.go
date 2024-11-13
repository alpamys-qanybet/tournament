package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

func main() {
	// fmt.Println("Division")

	// divideTeamsIntoDivision()
	generateMatches()
}

func divideTeamsIntoDivision() {
	var teams = [16]string{ // static 16 teams
		"Liverpool",
		"Arsenal",
		"Aston Villa",
		"Milan",
		"Juventus",
		"Barcelona",
		"Bayern Munchen",
		"Borussia Dortmund",
		"Manchester City",
		"Chelsea",
		"Manchester United",
		"Inter milan",
		"Atalanta",
		"Real Madrid",
		"Atletico Madrid",
		"Bayer Leverkusen",
	}
	// fmt.Println(teams)

	/*
		divisionATeams := make([]string, 0)                  // division A
		divisionBTeams := make([]string, 0)                  // division B
		divisionATeams = append(divisionATeams, teams[:]...) // add all teams to division A and them pick 8 one by one random

		// fmt.Println(divisionATeams)

		rand.Seed(uint64(time.Now().Unix()))
		for i := 0; i < len(teams)/2; i++ { // 8
			randomIndex := rand.Intn(len(divisionATeams))

			divisionBTeams = append(divisionBTeams, divisionATeams[randomIndex]) // add from A to B
			divisionATeams = append(divisionATeams[:randomIndex], divisionATeams[randomIndex+1:]...)
			// fmt.Println(divisionATeams)
			// fmt.Println(divisionBTeams)
			// fmt.Println()
			// fmt.Println()
		}

		fmt.Println(teams)          // 16 teams
		fmt.Println(divisionATeams) // 8 teams
		fmt.Println(divisionBTeams) // 8 teams
	*/

	n := 8

	teamsSlice := make([]string, 0)
	teamsSlice = append(teamsSlice, teams[:]...)

	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	r.Shuffle(len(teamsSlice), func(i, j int) { // shuffle teams(re-order with random indexes)
		teamsSlice[i], teamsSlice[j] = teamsSlice[j], teamsSlice[i]
	})

	divisionATeams := make([]string, 0)
	divisionATeams = append(divisionATeams, teamsSlice[:n]...) // len 8, cap 8
	// divisionATeams := teamsSlice[:n] // len 8, cap 16
	divisionBTeams := teamsSlice[n:] // len 8, cap 8

	fmt.Println(teams, len(teams), cap(teams)) // 16 teams
	fmt.Println()
	fmt.Println(teamsSlice, len(teamsSlice), cap(teamsSlice)) // 16 randomly ordered teams
	fmt.Println()
	fmt.Println(divisionATeams, len(divisionATeams), cap(divisionATeams)) // 8 teams
	fmt.Println()
	fmt.Println(divisionBTeams, len(divisionBTeams), cap(divisionBTeams)) // 8 teams

}

func generateMatches() { // assume taking 8 teams
	var teams = [8]string{ // static 8 teams
		"Liverpool",
		"Arsenal",
		"Aston Villa",
		"Milan",
		"Juventus",
		"Barcelona",
		"Bayern Munchen",
		"Borussia Dortmund",
	}

	fmt.Println(teams)

	teamsSlice := make([]string, 0)
	teamsSlice = append(teamsSlice, teams[:]...)

	a := GenerateRoundRobinTournamentMatchesByTeams(teamsSlice)
	fmt.Println(a)
}

// GenerateRoundRobinTournamentMatches generates a 2d slice of matches of a single round robin tournament.
// Each team will play one time against all other teams.
func GenerateRoundRobinTournamentMatchesByTeams(teams []string) [][]string {

	matches := make([][]string, 0)

	dummy := "even"
	if len(teams)%2 != 0 {
		dummy = uuid.New().String()
		teams = append(teams, dummy)
	}

	rotation := make([]string, len(teams))
	copy(rotation, teams)

	for i := 0; i < (len(teams) - 1); i++ {
		rotationLen := len(rotation)
		for i := 0; i < len(rotation); i = i + 2 {
			matches = append(matches, []string{rotation[i], rotation[i+1]})
		}

		// rotate teams for next round
		rotationHelper := append([]string{}, rotation[0])                     // append first team
		rotationHelper = append(rotationHelper, rotation[rotationLen-1])      // append last team
		rotationHelper = append(rotationHelper, rotation[1:rotationLen-1]...) // append remaining teams

		rotation = rotationHelper
	}

	// remove dummy matches
	if dummy != "even" {
		i := 0
		for _, x := range matches {
			if !stringSlicecontains(x, dummy) {
				matches[i] = x
				i++
			}
		}
		matches = matches[:i]
	}

	return matches
}

// Iterate over slice of string  to check whether it an element or not
func stringSlicecontains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
