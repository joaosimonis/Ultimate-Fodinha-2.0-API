package config

import (
	"fmt" // Erros formatados (fmt.Errorf)
	"os" // Interagir com o Sistema Operacional, principalmente ler variáveis de ambiente (os.LookupEnv) e verificar erro de arquivos inexistentes (os.IsNotExist)

	"github.com/joho/godotenv"
)

// Centraliza as configurações que a API precisa
type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// Carrega a configuração
func LoadConfig() (Config, error) {
	// 1. Carregar .env
	// Não existe -> tudo bem, segue o fluxo igual
	// Deu outro tipo de erro -> retorna o erro (erro ao abrir arquivo ou ao interpretar o conteúdo de .env)
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		// Retorno de erro
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	/* 
		Exemplo desse "outro erro" que pode acontecer
		No meu .env, mandei um DB-HOST=localhost (não está com o parser válido no nome da variável)
		Retornaria -> load .env: unexpected character "-" in variable name near "DB-HOST=localhost\n
	*/

	// 2. Lendo variáveis obrigatórias
	dbHost, err := requiredEnv("DB_HOST")
	if err != nil {
		return Config{}, err
	}

	dbUser, err := requiredEnv("DB_USER")
	if err != nil {
		return Config{}, err
	}

	dbPassword, err := requiredEnv("DB_PASSWORD")
	if err != nil {
		return Config{}, err
	}

	dbName, err := requiredEnv("DB_NAME")
	if err != nil {
		return Config{}, err
	}

	// 3. Monta o Struct final
	return Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     dbHost,
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}, nil // Retornamos nil como 2º valor retornado em caso de sucesso. Se desse erro, retornaríamos err 
}

// Lê uma variável de ambiente opcional
func getEnv(key, fallback string) string {
	// value: valor da variável de ambiente (é nativo do os.LookupEnv)
	// ok: um bool dizendo se ela existe ou não (é nativo do os.LookupEnv)
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	return value
}

// Lê uma variável de ambiente obrigatória
func requiredEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return "", fmt.Errorf("environment variable %s not set", key)
	}

	return value, nil
}
