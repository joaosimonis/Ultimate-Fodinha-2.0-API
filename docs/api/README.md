# Convenções e referência de rotas

[README do projeto](../../README.md#documentação)

## Endereço e formatos

URL base local: `http://localhost:8080/api/v1`. Todos os caminhos abaixo são relativos a `/api/v1`. Envie corpos JSON com `Content-Type: application/json`. IDs de entidades são UUIDs gerados pelo PostgreSQL; datas são strings RFC 3339, com fuso horário.

O retorno usual de sucesso é `{"data": ...}`. Erros tratados pelos handlers usam `{"error":"mensagem"}`. Respostas informativas usam `{"message":"mensagem"}`. Há duas exceções relevantes: o health check retorna `{"status":"ok"}` e a consulta/conexão de uma sala cancelada retorna `409` com `message`.

Campos com `omitempty` podem ser omitidos: um campo opcional ausente não deve ser presumido como zero. `placement` de participantes/equipes aparece como `null` antes do resultado. Listas vazias de partidas e rankings retornam `data: []`; a lista vazia de jogadores atualmente retorna `data: null`.

## Rotas

| Método | Caminho | Autenticação | Sucesso | Referência |
| --- | --- | --- | --- | --- |
| GET | `/health` | Pública | 200 | Abaixo |
| POST | `/auth/login` | Pública | 200 | [Autenticação](AUTENTICACAO.md) |
| GET | `/players` | Pública | 200 | [Jogadores](JOGADORES.md) |
| GET | `/players/:id` | Pública | 200 | [Jogadores](JOGADORES.md) |
| POST | `/players` | JWT admin | 201 | [Jogadores](JOGADORES.md) |
| DELETE | `/players/:id` | JWT admin | 200 | [Jogadores](JOGADORES.md) |
| GET | `/matches` | Pública | 200 | [Partidas](PARTIDAS.md) |
| GET | `/matches/:id` | Pública | 200 | [Partidas](PARTIDAS.md) |
| POST | `/matches` | Pública | 201 | [Partidas](PARTIDAS.md) |
| POST | `/matches/:id/players` | Pública | 200 | [Partidas](PARTIDAS.md) |
| DELETE | `/matches/:id/players/:player_id` | Pública | 200 | [Partidas](PARTIDAS.md) |
| PATCH | `/matches/:id/finish` | Pública | 200 | [Partidas](PARTIDAS.md) |
| PATCH | `/admin/matches/room/:room_code/cancel` | JWT admin | 200 | [Partidas](PARTIDAS.md) |
| GET | `/rooms/:room_code` | Pública | 200 | [Partidas](PARTIDAS.md) |
| GET | `/ws/rooms/:room_code` | Pública; escrita exige host token | 101 (upgrade) | [WebSocket](../websocket/WEBSOCKET.md) |
| GET | `/ranking` | Pública | 200 | [Ranking](RANKING.md) |
| GET | `/ranking/:mode` | Pública | 200 | [Ranking](RANKING.md) |
| GET | `/ranking/truco/teams` | Pública | 200 | [Ranking](RANKING.md) |

O JWT administrativo e o `host_token` de uma sala têm usos distintos. As mutações HTTP públicas de partidas não exigem nenhum desses tokens na implementação atual.

## Códigos de resposta

| Status | Uso |
| --- | --- |
| 200 | Consulta ou alteração concluída |
| 201 | Jogador ou partida criada |
| 204 | Preflight `OPTIONS` do middleware CORS |
| 400 | JSON inválido, campos obrigatórios ausentes ou regra de domínio rejeitada |
| 401 | Credenciais inválidas ou JWT ausente, inválido ou expirado |
| 404 | Entidade não encontrada nas operações que tratam esse caso |
| 409 | Nome de jogador em uso, conflito de código ou sala cancelada |
| 500 | Falha interna ou erro de banco sem mapeamento específico |

UUIDs malformados e violações de chave estrangeira não têm tratamento dedicado e podem gerar `500`. Rotas inexistentes e falhas no upgrade WebSocket podem usar respostas do framework/biblioteca, fora do envelope acima. Não há paginação nas listagens de jogadores e partidas; rankings aceitam `limit`.

## Health check

```bash
curl http://localhost:8080/api/v1/health
```

```json
{"status":"ok"}
```

A rota responde se o servidor estiver atendendo; não executa uma consulta ao banco a cada chamada.

## CORS

`ALLOWED_ORIGINS` é uma lista de origens separadas por vírgula, com correspondência exata. O middleware libera `Content-Type`, `Authorization`, `Accept`, `Origin`, `Cache-Control` e `X-Requested-With`; métodos anunciados: `GET, POST, PUT, PATCH, DELETE, OPTIONS`. Também define `Access-Control-Allow-Credentials: true` para origens permitidas. Não há expansão de `*`. O upgrade WebSocket usa uma política própria descrita no [guia de WebSocket](../websocket/WEBSOCKET.md).
