package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	SSLMode  string
}

func NewDB(context context.Context, cnfg *PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cnfg.User, cnfg.Password, cnfg.Host, cnfg.Port, cnfg.Dbname, cnfg.SSLMode)

	//sql db for migration
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	//migration
	migration, err := migration(db)
	if err != nil {
		return nil, err
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// pgx Pool
	dbPool, err := pgxpool.New(context, dsn)
	if err != nil {
		return nil, err
	}

	if err := dbPool.Ping(context); err != nil {
		return nil, err
	}

	return dbPool, nil
}

func migration(db *sql.DB) (*migrate.Migrate, error) {
	//driver for migration
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	//migration
	migrate, err := migrate.NewWithDatabaseInstance(
		"file://internal/migrations/postgres",
		"postgres",
		driver,
	)
	if err != nil {
		return nil, err
	}

	return migrate, nil
}
