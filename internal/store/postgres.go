package store

import (
	"analytics/internal/model"
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresStore struct {
	db *sqlx.DB
}

func NewPostgres(databaseURL string) (*PostgresStore, error) {
	db, err := sqlx.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать пул соед БД: %w", err)
	}
	
	if err := db.PingContext(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("не удалось создать пул соед БД: %w", err)
	}

	m, err := migrate.New("file:///app/migrations", databaseURL)
    if err != nil {
		db.Close()
        return nil, fmt.Errorf("не удалось создать экземпляр migrate: %w", err)
    }
    defer m.Close()

    err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
		db.Close()
        return nil, err
    }

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

func (s *PostgresStore) SaveAnalytics(ctx context.Context, analytics model.EventAnalytics) error {
	sqlquery := `
	INSERT INTO event_analytics (
    event_type,
    window_start,
    window_end,
    count
	)
	VALUES ($1, $2, $3, 1)
	ON CONFLICT (event_type, window_start)
	DO UPDATE SET
    	count = event_analytics.count + 1;
	`
	if _, err := s.db.ExecContext(ctx, sqlquery, analytics.EventType, analytics.WindowStart, analytics.WindowEnd);
	err != nil {
		return fmt.Errorf("не получилось сохранить аналитику: %w", err)
	}
	
	return nil
}