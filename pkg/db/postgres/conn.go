package postgres

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/0sokrat0/BookAPI/internal/config"
	"github.com/0sokrat0/BookAPI/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, cfg *config.Config) (*postgres, error) {

	lg := logger.FromContext(ctx)
	if lg == nil {
		return nil, errors.New("logger not found in context")
	}

	var err error

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?pool_max_conns=%d&pool_min_conns=%d",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.MaxConn,
		cfg.Database.MinConn,
	)

	pgOnce.Do(func() {
		var db *pgxpool.Pool
		db, err = pgxpool.New(ctx, connString)
		if err != nil {
			lg.Errorf("error creating pg pool", zap.Error(err))
			return
		}

		pgInstance = &postgres{db: db}
	})

	if err != nil {
		return nil, err
	}

	migrateString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	m, err := migrate.New(
		"file://migrations",
		migrateString)
	if err != nil {
		lg.Errorf("unable to create migration", zap.Error(err))
		return nil, err
	}

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		if err.Error() == "Dirty database version 1. Fix and force version." {
			lg.Warn("Database is dirty, forcing version 1...")
			if forceErr := m.Force(1); forceErr != nil {
				lg.Errorf("Failed to force version: %v", forceErr)
				return nil, forceErr
			}

			err = m.Up()
			if err != nil && !errors.Is(err, migrate.ErrNoChange) {
				lg.Errorf("Unable to migrate database after forcing version: %v", err)
				return nil, err
			}
		} else {
			lg.Errorf("unable to migrate database", zap.Error(err))
			return nil, err
		}
	}

	lg.Info("Pools created successfully")

	return pgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}
