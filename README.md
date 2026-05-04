# Ultimate Fodinha API

API em Go para gerenciar jogadores, partidas, salas ao vivo, ranking e eventos em tempo real do projeto Ultimate Fodinha 2.0.

O projeto usa:

- Go
- Gin
- PostgreSQL
- pgx
- golang-migrate
- Docker Compose
- WebSocket com `github.com/coder/websocket`

## Requisitos

- Go 1.26+
- Docker e Docker Compose
- PostgreSQL, caso prefira rodar o banco sem Docker
- CLI `migrate` para aplicar migrations

## Configuração

Copie o arquivo de exemplo e ajuste as variáveis:

```bash
cp .env.example .env
```

## Rodando Localmente

### 1. Suba o banco

```bash
docker compose up -d db
```

### 2. Carregue as variáveis do `.env`

```bash
set -a
source .env
set +a
```

### 3. Rode as migrations

Se ainda não tiver a CLI `migrate`, instale:

```bash
curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

Aplique as migrations:

```bash
migrate -path internal/db/migrations \
  -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" \
  up
```

### 4. Rode o seed

```bash
go run ./cmd/seed
```

### 5. Inicie a API

```bash
go run ./cmd/api
```

A aplicação ficará disponível em:

```txt
http://localhost:8080
```

## Docker

Para subir API e banco:

```bash
docker compose up --build
```

Para subir também o pgAdmin:

```bash
docker compose --profile admin up --build
```
