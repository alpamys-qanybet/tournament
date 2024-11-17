package dto

type TeamDTO struct {
	Id            uint16 `json:"id"`
	Name          string `json:"name"`
	Division      string `json:"division,omitempty"`
	Wins          uint16 `json:"wins,omitempty"`
	Draws         uint16 `json:"draws,omitempty"`
	Loses         uint16 `json:"loses,omitempty"`
	GoalsScored   uint16 `json:"goals_scored,omitempty"`
	GoalsConceded uint16 `json:"goals_conceded,omitempty"`
	GoalDiff      int16  `json:"goal_diff,omitempty"`
	Points        uint16 `json:"points,omitempty"`
}

func (t *TeamDTO) SetDiff() {
	t.GoalDiff = int16(t.GoalsScored) - int16(t.GoalsConceded)
}
