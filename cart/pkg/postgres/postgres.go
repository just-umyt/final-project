package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	SSLMode  string
}

func NewDB(cnfg *PostgresConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cnfg.User, cnfg.Password, cnfg.Host, cnfg.Port, cnfg.Dbname, cnfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewDBPool(context context.Context, cnfg *PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cnfg.User, cnfg.Password, cnfg.Host, cnfg.Port, cnfg.Dbname, cnfg.SSLMode)

	dbPool, err := pgxpool.New(context, dsn)
	if err != nil {
		return nil, err
	}

	if err := dbPool.Ping(context); err != nil {
		return nil, err
	}

	return dbPool, nil
}
