# Execução e configuração

[README do projeto](../../README.md#documentação)

## Requisitos

- Go 1.26.1 ou compatível com o `go.mod`.
- PostgreSQL (o Compose usa `postgres:16-alpine`).
- Docker e Docker Compose para o ambiente em contêineres.
- CLI `migrate` do golang-migrate, com suporte a PostgreSQL, disponível no `PATH`.

Execute os comandos a partir da raiz deste repositório da API.

## Ambiente

```bash
cp .env.example .env
```

Ajuste os valores do `.env`, principalmente senha do banco e segredo JWT. A API carrega o arquivo automaticamente; variáveis já exportadas têm precedência.

| Variável | Obrigatória para a API | Padrão/função |
| --- | --- | --- |
| `APP_PORT` | Não | `PORT`, se definida; senão `8080` |
| `ALLOWED_ORIGINS` | Não | Vazia; origens HTTP permitidas, separadas por vírgula |
| `JWT_SECRET` | Sim | Segredo de assinatura dos tokens administrativos |
| `JWT_EXPIRES_IN_SECONDS` | Não | `3600`; inteiro positivo |
| `DB_HOST` | Sim | Host do PostgreSQL |
| `DB_PORT` | Não | `5432` |
| `DB_USER` | Sim | Usuário do banco |
| `DB_PASSWORD` | Sim | Senha do banco |
| `DB_NAME` | Sim | Nome do banco |
| `DB_SSLMODE` | Não | `disable` |
| `PGADMIN_DEFAULT_EMAIL` | Não | Configuração exclusiva do Compose/pgAdmin |
| `PGADMIN_DEFAULT_PASSWORD` | Não | Configuração exclusiva do Compose/pgAdmin |
| `PGADMIN_PORT` | Não | `5050` no Compose |

O `.env.example` contém sugestões diferentes de alguns defaults internos; copiá-lo define esses valores explicitamente.

## Banco e migrations

```bash
docker compose up -d db
set -a
source .env
set +a
migrate -path internal/db/migrations \
  -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" \
  up
```

O comando pressupõe valores definidos no `.env.example`; caracteres reservados em usuário/senha precisam ser codificados para uso na URL de conexão. O carregamento com `source` funciona em Bash/Zsh e pressupõe que o `.env` tenha sintaxe válida de shell.

As sete migrations criam jogadores, partidas, participantes, estatísticas, usuários administrativos, data de cancelamento e equipes. A migration `007` também associa equipes aos registros existentes. A API e o Compose não aplicam migrations automaticamente.

## Seed e administrador

```bash
go run ./cmd/seed
```

O seed inclui o catálogo de jogadores e ignora conflitos no nome exato (`ON CONFLICT (name) DO NOTHING`); ele não atualiza registros já existentes e não cria usuários administrativos.

Para criar ou atualizar um administrador, exporte as duas variáveis antes de executar o comando:

```bash
export ADMIN_USERNAME=admin
read -rs 'ADMIN_PASSWORD?Senha do administrador: '
echo
export ADMIN_PASSWORD
go run ./cmd/admin
unset ADMIN_PASSWORD
```

O exemplo de leitura de senha acima usa Zsh. Em Bash, use `read -rs -p 'Senha do administrador: ' ADMIN_PASSWORD`. O comando procura `ADMIN_USERNAME` e `ADMIN_PASSWORD` antes de carregar `.env`, portanto essas duas variáveis devem estar no ambiente do processo. As demais configurações da API também são necessárias, incluindo `JWT_SECRET`. Se o username já existir, ignorando maiúsculas/minúsculas, a senha será substituída e o papel definido como `admin`.

## Iniciar a API

```bash
go run ./cmd/api
```

Verifique com `curl http://localhost:8080/api/v1/health`.

## Docker

Após configurar `.env` e aplicar as migrations ao banco:

```bash
docker compose up --build
```

Para incluir o pgAdmin:

```bash
docker compose --profile admin up --build
```

Dentro da rede do Compose, o banco é `db:5432`; a API recebe esses valores automaticamente. No pgAdmin, use esse host para cadastrar a conexão. Do host, use `localhost` e a porta publicada. Dados persistem nos volumes `fodinha-db-data` e `pgadmin-data`.

O Compose publica `${APP_PORT:-8080}:8080`, mas também passa `APP_PORT` ao processo. Para usar a configuração atual sem alterações, mantenha `APP_PORT=8080`; trocar a porta exige alinhar a porta interna do mapeamento. A imagem final contém apenas o executável da API, sem as CLIs de seed, admin ou migrations.

## Verificação e problemas comuns

```bash
go test ./...
go vet ./...
```

- Falha de configuração na inicialização: confira as variáveis obrigatórias e se `JWT_EXPIRES_IN_SECONDS` é positivo.
- Tabela ou coluna inexistente: aplique todas as migrations no banco usado pela instância.
- Login retorna `401`: crie o administrador e confira suas credenciais.
- Navegador bloqueia HTTP por CORS: inclua a origem completa, com protocolo e porta, em `ALLOWED_ORIGINS`.
- Conexão recusada ao banco: confirme se o acesso parte do host (`localhost`) ou de um contêiner (`db`).
