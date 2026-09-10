package usecase

import (
	"analytics/internal/model"
	"analytics/internal/store"
	"context"
	"time"
)

type analyticsUsecase struct {
	store store.AnalyticsStore
}

func NewAnalyticsUsecase(store store.AnalyticsStore) (*analyticsUsecase) {
	return &analyticsUsecase{store: store}
}

func(au *analyticsUsecase) Analyse(ctx context.Context, event model.Event) ( error) {
	windowStart := event.Timestamp.Truncate(30 * time.Second)
	windowEnd := windowStart.Add(30 * time.Second)
	
	analytics := model.EventAnalytics{
		EventType: event.EventType,
		WindowStart: windowStart,
		WindowEnd: windowEnd,
	}
	return au.store.SaveAnalytics(ctx, analytics)
}