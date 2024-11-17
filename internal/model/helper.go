package model

import (
	"strconv"
	"tournament/dto"
)

func prepareTeamInIdsStr(teams []*dto.TeamDTO) (s string) {

	for _, v := range teams {
		s += strconv.Itoa(int(v.Id)) + ","
	}

	if len(s) > 0 {
		s = s[:len(s)-1] // Remove the last comma
	}

	return
}
