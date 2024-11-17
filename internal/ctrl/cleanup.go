package ctrl

import (
	"context"
	"tournament/internal/model"
)

func Cleanup(ctx context.Context) error {
	return model.Cleanup(ctx)
}
