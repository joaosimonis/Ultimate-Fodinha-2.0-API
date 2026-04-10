package main

import (
	"log"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/bootstrap"
)

func main() {
	app, err := bootstrap.BuildApp()
	if err != nil {
		log.Fatalf("erro ao inicializar aplicação: %v", err)
	}
	defer app.Close()

	log.Printf("servidor rodando na porta %s", app.Config.AppPort)
	log.Fatal(app.Run())
}
