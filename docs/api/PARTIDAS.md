# Partidas e salas

[README do projeto](../../README.md#documentação) · [Rotas](README.md)

As rotas HTTP de partidas são públicas, exceto o cancelamento administrativo. `:id` é o UUID da partida; `:room_code` é o código da sala; `:player_id` é o UUID do jogador cadastrado.

## Modelo de resposta

Exemplo de detalhe de uma partida de Fodinha com um participante representado:

```json
{
  "data": {
    "id": "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
    "mode": "fodinha",
    "status": "ongoing",
    "room_code": "ABC234",
    "lives": 3,
    "created_at": "2026-11-13T12:00:00Z",
    "players": [{
      "id": "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb",
      "match_id": "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
      "player_id": "11111111-1111-4111-8111-111111111111",
      "placement": null,
      "is_winner": false
    }]
  }
}
```

| Campo de partida | Significado |
| --- | --- |
| `id`, `mode`, `status`, `room_code`, `created_at` | Identificação e estado; modos `fodinha`/`truco`, status `ongoing`/`finished`/`cancelled` |
| `host_token` | String retornada somente na criação; guarde-a para publicar estado ao vivo |
| `lives` | Vidas iniciais da Fodinha |
| `max_score` | Pontuação máxima do Truco |
| `live_state` | Último JSON enviado pelo host, quando existente |
| `finished_at`, `cancelled_at` | Datas presentes após a respectiva operação |
| `players` | Participações, quando carregadas e não vazias |
| `teams` | Equipes, quando carregadas e não vazias |

Cada participação tem `id`, `match_id`, `player_id`, `placement` (inteiro ou `null`) e `is_winner`. Pode incluir `team_id`, `team_name`, `final_score` e `remaining_lives`. O `id` da participação é diferente do `player_id` do jogador.

Cada equipe tem `id`, `match_id`, `team_name`, `placement` (inteiro ou `null`), `is_winner` e, após resultado, `final_score`. Os resultados das equipes são agregados a partir dos participantes. Arrays vazios e outros campos opcionais podem ser omitidos.

## Criar partida

`POST /api/v1/matches` — retorna `201` com a partida criada, ainda sem participantes, com `status: ongoing`, `room_code` e `host_token`.

Fodinha:

```bash
curl -X POST http://localhost:8080/api/v1/matches \
  -H 'Content-Type: application/json' \
  -d '{"mode":"fodinha","lives":3}'
```

Truco:

```bash
curl -X POST http://localhost:8080/api/v1/matches \
  -H 'Content-Type: application/json' \
  -d '{"mode":"truco","max_score":15}'
```

| Modo | Obrigatório | Não permitido |
| --- | --- | --- |
| `fodinha` | `lives`: inteiro de 1 a 5 | `max_score` com valor não nulo |
| `truco` | `max_score`: 15, 24 ou 30 | `lives` com valor não nulo |

`mode` deve corresponder exatamente a um desses valores. Código de sala e token são gerados pelo servidor; não são parâmetros de criação. O código possui seis caracteres de `ABCDEFGHJKLMNPQRSTUVWXYZ23456789`. O token é aleatório, e somente seu hash SHA-256 é persistido. Consultas posteriores não devolvem o token e não há rota para recuperá-lo.

## Listar e consultar

- `GET /api/v1/matches`: retorna `200` com uma lista, da mais recente para a mais antiga. Cada item contém os campos da partida, sem carregar `players` ou `teams`.
- `GET /api/v1/matches?status=ongoing`: filtra pelo status. `live` é um alias de `ongoing`; `finished` e `cancelled` filtram os respectivos estados. Um valor desconhecido é repassado ao filtro e normalmente resulta em lista vazia, sem `400`.
- `GET /api/v1/matches/:id`: retorna o detalhe com participantes e equipes, inclusive se a partida estiver cancelada.
- `GET /api/v1/rooms/:room_code`: retorna o mesmo detalhe, mas rejeita partidas canceladas com `409`, `{"message":"Essa partida foi cancelada pelo admin."}`.

As consultas por código de sala são sensíveis a maiúsculas/minúsculas; use o código exatamente como retornado na criação. Partidas inexistentes retornam `404`. Não há paginação.

## Adicionar participantes

`POST /api/v1/matches/:id/players` — retorna `200` com o detalhe atualizado.

Fodinha:

```json
{
  "players": [
    {"player_id":"11111111-1111-4111-8111-111111111111"},
    {"player_id":"22222222-2222-4222-8222-222222222222"}
  ]
}
```

Truco em duplas:

```json
{
  "players": [
    {"player_id":"11111111-1111-4111-8111-111111111111","team_name":"Equipe A"},
    {"player_id":"22222222-2222-4222-8222-222222222222","team_name":"Equipe A"},
    {"player_id":"33333333-3333-4333-8333-333333333333","team_name":"Equipe B"},
    {"player_id":"44444444-4444-4444-8444-444444444444","team_name":"Equipe B"}
  ]
}
```

Cada chamada deve conter pelo menos dois jogadores cadastrados, IDs não vazios e sem duplicação. Truco exige exatamente quatro ou seis jogadores no corpo, divididos em duas equipes de dois ou três integrantes. `team_name` é obrigatório no Truco; espaços nas extremidades são removidos e os nomes diferenciam maiúsculas/minúsculas. O servidor cria/reutiliza as equipes e devolve seus IDs.

A operação adiciona participantes; não substitui a lista anterior. Um jogador já associado gera `400`. As regras de quantidade/equilíbrio são verificadas no corpo da chamada, não no conjunto acumulado. Para a configuração usual, envie todos os participantes de uma vez. A implementação não bloqueia essa operação pelo status da partida.

## Remover participante

`DELETE /api/v1/matches/:id/players/:player_id` — sem corpo, retorna `200` com a partida atualizada ou `404` se a associação não existir. Remove a participação, preservando o cadastro do jogador. Não revalida quantidade/equilíbrio, não remove equipes vazias e não verifica o status da partida.

## Finalizar Fodinha

`PATCH /api/v1/matches/:id/finish`:

```json
{
  "players": [
    {"player_id":"11111111-1111-4111-8111-111111111111","remaining_lives":2,"placement":1},
    {"player_id":"22222222-2222-4222-8222-222222222222","remaining_lives":0,"placement":2}
  ]
}
```

Envie exatamente os jogadores associados, uma vez cada. `remaining_lives` é obrigatório, entre zero e as vidas iniciais. `placement` é obrigatório, único e entre 1 e a quantidade de participantes; o primeiro colocado é o vencedor. A validação não deriva colocação das vidas nem exige que os demais tenham zero vidas. Não envie equipes. `final_score` não faz parte do corpo de resultado de Fodinha.

## Finalizar Truco

Use o mesmo endpoint com as duas equipes:

```json
{
  "teams": [
    {"team_id":"aaaaaaaa-1111-4111-8111-111111111111","final_score":15},
    {"team_id":"bbbbbbbb-2222-4222-8222-222222222222","final_score":9}
  ]
}
```

Substitua os IDs pelos de `teams` retornados na associação/consulta. Como alternativa, envie `team_name` com o nome exato cadastrado. Se ambos forem informados, `team_id` tem precedência, inclusive quando inválido. Prefira o ID para evitar depender da apresentação do nome.

Envie exatamente duas equipes distintas e conhecidas. `final_score` é obrigatório para ambas, entre zero e `max_score`. Uma equipe precisa atingir exatamente `max_score`, com pontuações diferentes. Não envie uma lista não vazia em `players`. O servidor atribui aos integrantes a pontuação da equipe, colocação 1 ou 2 e `is_winner`.

Nos dois modos, sucesso retorna `200` com o detalhe, `status: finished` e `finished_at`. Partida já finalizada retorna `400`; cancelada retorna `409` com `error`. Resultado final não é derivado de `live_state`. A finalização não transmite eventos WebSocket automaticamente.

## Cancelar sala

`PATCH /api/v1/admin/matches/room/:room_code/cancel` — exige JWT administrativo; sem corpo.

```bash
curl -X PATCH http://localhost:8080/api/v1/admin/matches/room/ABC234/cancel \
  -H "Authorization: Bearer $TOKEN"
```

Retorna `200` com a partida e `status: cancelled`. Define `cancelled_at` apenas se ainda não existir. O código é comparado sem diferenciar maiúsculas/minúsculas nesta operação. Partida inexistente retorna `404`. O cancelamento também aceita partidas finalizadas; como o ranking considera apenas `finished`, a partida deixa de compô-lo. Não apaga os resultados nem `finished_at`.

Novas consultas por sala e novas conexões WebSocket recebem `409`. Conexões já abertas não são encerradas ou notificadas automaticamente; atualizações de estado não verificam status.

## Erros e limites atuais

JSON inválido retorna `400`, `dados inválidos`. Validações de configuração, participantes e resultados retornam `400` com a mensagem do domínio. Exemplos: `modo de partida inválido`, `pontuação máxima não permitida para truco`, `jogador duplicado na partida`.

Conflito ao gerar código retorna `409`. IDs de jogadores inexistentes, UUIDs malformados e outras falhas não mapeadas podem retornar `500`. Erros internos usam mensagens específicas da operação, como `erro ao criar partida` e `erro ao finalizar partida`. Consulte as mensagens completas em `internal/matches/match_errors.go`.
