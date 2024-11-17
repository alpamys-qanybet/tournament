package ctrl

import (
	"context"
	"errors"
	"strings"
	"tournament/internal/model"
)

func GetTeamList(ctx context.Context) (interface{}, error) {
	return model.GetTeamList(ctx)
}

func CreateTeam(ctx context.Context, name string) (uint16, error) {
	if strings.Trim(name, " ") == "" {
		return uint16(0), errors.New("create_team_failure_name_is_required")
	}

	return model.CreateTeam(ctx, name)
}

func GenerateTeams(ctx context.Context) error {
	return model.GenerateTeams(ctx)
}
