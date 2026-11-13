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

## Documentação

| Documento | Conteúdo |
| --- | --- |
| [Execução e configuração](docs/projeto/CONFIGURACAO.md) | Requisitos, ambiente, banco, migrations, seed, administrador e Docker |
| [Arquitetura e persistência](docs/projeto/ARQUITETURA.md) | Organização do código, responsabilidades e tabelas |
| [Convenções e rotas](docs/api/README.md) | URL base, respostas, erros, autenticação por rota e health check |
| [Autenticação](docs/api/AUTENTICACAO.md) | Login administrativo e token JWT |
| [Jogadores](docs/api/JOGADORES.md) | Consulta, cadastro e exclusão |
| [Partidas e salas](docs/api/PARTIDAS.md) | Modelos, criação, participantes, resultados e cancelamento |
| [Ranking](docs/api/RANKING.md) | Filtros, critérios de ordenação e formatos de retorno |
| [WebSocket](docs/websocket/WEBSOCKET.md) | Conexão, eventos, estado ao vivo e exemplo de cliente |

## Por onde começar

1. Prepare o ambiente seguindo o [guia de configuração](docs/projeto/CONFIGURACAO.md).
2. Consulte as [convenções e a tabela de rotas](docs/api/README.md).
3. Liste os jogadores, crie uma partida e associe seus participantes pelo [guia de partidas](docs/api/PARTIDAS.md).
4. Guarde o `host_token` retornado na criação para publicar o estado ao vivo via [WebSocket](docs/websocket/WEBSOCKET.md).
5. Finalize a partida via HTTP para que seus resultados entrem no ranking.

## Manutenção

Ao alterar uma rota, revise também os modelos, as validações do serviço, as consultas do repositório e o guia correspondente. Os exemplos usam dados fictícios; substitua IDs, códigos de sala e tokens pelos retornados pela sua instância. Não há Swagger/OpenAPI servido pela aplicação atualmente.
