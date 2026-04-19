package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctxx context.Context, databaseUrl string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(ctxx, 10*time.Second)
	defer cancel()
	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		log.Fatalf("Unable to parse database url: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unbale to create connection pool: %v", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL database")
	return pool
}
