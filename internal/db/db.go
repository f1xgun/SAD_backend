package db

import (
	"context"
	"fmt"
	"sad/internal/config"

	"github.com/jackc/pgx/v5"
)

func NewDBConnection(config config.Config) (*pgx.Conn, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DBUserName, config.DBUserPassword, config.DBHost, config.DBPort, config.DBName)

	conn, err := pgx.Connect(context.Background(), url)

	return conn, err
}
