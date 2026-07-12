package migration

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"subscriptions/internal/config"
)

const defaultMigrationsDir = "migrations"

type Options struct {
	DatabaseURL   string
	MigrationsDir string
}

func DefaultOptions() (*Options, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return &Options{DatabaseURL: cfg.Database.DSN(), MigrationsDir: defaultMigrationsDir}, nil
}

func Up(options *Options) error {
	return withDB(options, func(db *sql.DB, dir string) error { return goose.Up(db, dir) })
}
func Down(options *Options) error {
	return withDB(options, func(db *sql.DB, dir string) error { return goose.Down(db, dir) })
}
func Status(options *Options) error {
	return withDB(options, func(db *sql.DB, dir string) error { return goose.Status(db, dir) })
}

func withDB(options *Options, operation func(*sql.DB, string) error) error {
	if options == nil {
		var err error
		options, err = DefaultOptions()
		if err != nil {
			return err
		}
	}
	if options.MigrationsDir == "" {
		options.MigrationsDir = defaultMigrationsDir
	}
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	db, err := sql.Open("pgx", options.DatabaseURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	return operation(db, options.MigrationsDir)
}
