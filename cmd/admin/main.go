package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/bootstrap"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	username := strings.TrimSpace(os.Getenv("ADMIN_USERNAME"))
	password := os.Getenv("ADMIN_PASSWORD")

	if username == "" {
		log.Fatal("ADMIN_USERNAME not set")
	}

	if password == "" {
		log.Fatal("ADMIN_PASSWORD not set")
	}

	app, err := bootstrap.BuildApp()
	if err != nil {
		log.Fatalf("erro ao inicializar aplicação: %v", err)
	}
	defer app.Close()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("erro ao gerar hash da senha: %v", err)
	}

	_, err = app.DB.Exec(context.Background(), `
		INSERT INTO users (username, password_hash, role)
		VALUES ($1, $2, 'admin')
		ON CONFLICT (LOWER(username))
		DO UPDATE SET
			password_hash = EXCLUDED.password_hash,
			role = 'admin',
			updated_at = NOW()
	`, username, string(passwordHash))
	if err != nil {
		log.Fatalf("erro ao criar admin: %v", err)
	}

	fmt.Printf("admin %q criado/atualizado com sucesso\n", username)
}
