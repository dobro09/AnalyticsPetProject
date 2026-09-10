package store

import (
	"analytics/internal/model"
	"context"
)

type AnalyticsStore interface {
	SaveAnalytics(ctx context.Context, analytics model.EventAnalytics) error
}