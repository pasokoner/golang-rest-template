package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
	schema   = os.Getenv("DB_SCHEMA")
)

func Connect(connStr string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		return nil, err
	}

	return dbpool, nil
}
