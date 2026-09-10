package usecase

import (
	"analytics/internal/model"
	"context"
)

type AnalyticsUsecase interface {
	Analyse(ctx context.Context, event model.Event) error
}