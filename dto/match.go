package dto

// match types
const (
	MatchTypeDivisionA      = "DA" // division A
	MatchTypeDivisionB      = "DB" // division B
	MatchTypePlayoffQuarter = "PQ" // play-off quarter-final
	MatchTypePlayoffSemi    = "PS" // play-off semi-final
	MatchTypePlayoffFinal   = "PF" // play-off final
)

type MatchDTO struct {
	Id              uint16  `json:"id"`
	Name            string  `json:"name"`
	FirstTeamId     uint16  `json:"first_team_id"`
	SecondTeamId    uint16  `json:"second_team_id"`
	FirstTeamScore  *uint16 `json:"first_team_score"`
	SecondTeamScore *uint16 `json:"second_team_score"`
	WinnerId        *int16  `json:"winner_id,omitempty"`
	MatchType       string  `json:"match_type,omitempty"`
	Played          bool    `json:"played,omitempty"`
}
