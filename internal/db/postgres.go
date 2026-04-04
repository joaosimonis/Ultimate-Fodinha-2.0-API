package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jsimonis/ultimate-fodinha-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg config.Config) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("erro ao criar pool de conexões: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("banco inacessível: %v", err)
	}

	log.Println("✅ conexão com PostgreSQL estabelecida")
	return pool
}