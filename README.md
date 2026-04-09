# Ultimate Fodinha API

API em Go para gerenciar os dados do projeto Ultimate Fodinha 2.0, usando PostgreSQL como banco de dados.

## Requisitos

- Go 1.26+
- Docker e Docker Compose
- PostgreSQL, caso prefira rodar o banco sem Docker

## Configuração

Copie o arquivo de exemplo e ajuste as variáveis de ambiente:

```bash
cp .env.example .env
```

## Fluxo recomendado para rodar localmente

### 1. Suba o banco

Se quiser subir apenas o PostgreSQL com Docker:

```bash
docker compose up -d db
```

### 2. Carregue as variáveis do `.env` no shell

Isso ajuda a reaproveitar os mesmos valores nos comandos de migration e seed:

```bash
set -a
source .env
set +a
```

### 3. Rode as migrations

⚠️ Certifique-se de que o banco de dados esteja rodando antes de aplicar as migrations.

Se ainda não tiver a CLI `migrate` instalada, no Linux você pode fazer:

```bash
curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

Depois disso, aplique todas as migrations com:

```bash
migrate -path internal/db/migrations \
  -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" \
  up
```

Se preferir, você também pode usar valores fixos:

```bash
migrate -path internal/db/migrations \
  -database "postgres://postgres:postgres@localhost:5432/fodinha?sslmode=disable" \
  up
```

### 4. Rode o seed

Para popular a tabela de jogadores com os dados iniciais:

```bash
go run ./cmd/seed
```

### 5. Inicie a API

```bash
go run ./cmd/api
```

A aplicação ficará disponível em `http://localhost:8080`.

## Testando a API

Após iniciar a aplicação, você pode testar se está funcionando:

```bash
curl http://localhost:8080/api/v1/health
```

Ou testar o endpoint no Postman.

## Rodando com Docker

Para subir a API e o banco:

```bash
docker compose up --build
```

Para subir também o pgAdmin:

```bash
docker compose --profile admin up --build
```

## Estrutura do projeto

- `cmd/api`: ponto de entrada da aplicação
- `cmd/seed`: carga inicial de dados
- `internal/db`: conexão com PostgreSQL e arquivos de migration
- `internal/http/routes`: registro de rotas
- `internal/players`, `internal/matches`, `internal/ranking`: domínios da aplicação
