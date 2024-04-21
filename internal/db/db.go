package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"sad/internal/config"
)

func NewDBConnection(config config.Config) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DBUserName, config.DBUserPassword, config.DBHost, config.DBPort, config.DBName)

	conn, err := pgxpool.New(context.Background(), url)

	return conn, err
}
