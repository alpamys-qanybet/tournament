package dto

const (
	DivisionA = "A"
	DivisionB = "B"
)

type DivisionDTO struct {
	Name    string      `json:"name"`
	Teams   []*TeamDTO  `json:"teams"`
	Matches []*MatchDTO `json:"matches"`
}
