package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jsimonis/ultimate-fodinha-api/internal/config"

	"github.com/jackc/pgx/v5"
)

func Connect(cfg config.Config) *pgx.Conn {
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

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}

	log.Println("conexão com PostgreSQL estabelecida com sucesso")
	return conn
}