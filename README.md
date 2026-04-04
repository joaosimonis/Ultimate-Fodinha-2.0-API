# Ultimate Fodinha API

API em Go para gerenciar dados do projeto Ultimate Fodinha 2.0, com PostgreSQL como banco de dados.

## Requisitos

- Go 1.26+
- Docker e Docker Compose
- PostgreSQL, caso rode sem Docker

## Configuração

Copie o arquivo de exemplo e ajuste as variáveis:

```bash
cp .env.example .env
```

## Rodando localmente

Se quiser subir apenas o banco com Docker:

```bash
docker compose up -d db
```

Depois rode a API localmente:

```bash
go run ./cmd/api
```

A aplicação ficará disponível em `http://localhost:8080`.

## Rodando com Docker

Para subir a API e o banco:

```bash
docker compose up --build
```

Para subir também o pgAdmin:

```bash
docker compose --profile admin up --build
```

## Estrutura

- `cmd/api`: ponto de entrada da aplicação
- `internal/db`: conexão com PostgreSQL e migrations
- `internal/routes`: registro de rotas
- `internal/players`, `internal/matches`, `internal/ranking`: domínios da aplicação

