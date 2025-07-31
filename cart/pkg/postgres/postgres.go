package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

const (
	driverName = "postgres"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
	SSLMode  string
}

func NewDB(cnfg *PostgresConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s", driverName, cnfg.User, cnfg.Password, cnfg.Host, cnfg.Port, cnfg.DBname, cnfg.SSLMode)

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewDBPool(context context.Context, cnfg *PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s", driverName, cnfg.User, cnfg.Password, cnfg.Host, cnfg.Port, cnfg.DBname, cnfg.SSLMode)

	dbPool, err := pgxpool.New(context, dsn)
	if err != nil {
		return nil, err
	}

	if err := dbPool.Ping(context); err != nil {
		return nil, err
	}

	return dbPool, nil
}
