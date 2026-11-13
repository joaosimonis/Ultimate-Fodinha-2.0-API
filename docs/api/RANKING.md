# Rankings

[README do projeto](../../README.md#documentação) · [Rotas](README.md)

Todas as rotas de ranking são públicas e retornam `200` com `{"data":[...]}`. Sem resultados, `data` é `[]`. As consultas consideram apenas partidas com `status: finished`; estado ao vivo e partidas canceladas não contam.

## Ranking individual

Use `GET /api/v1/ranking?mode=fodinha` ou `GET /api/v1/ranking/fodinha`. Para Truco, use `GET /api/v1/ranking/truco?team_size=2`.

| Parâmetro | Regra |
| --- | --- |
| `mode` | Obrigatório na rota `/ranking`; `fodinha` ou `truco`. Também pode vir de `/:mode`; se a query `mode` não estiver vazia, tem precedência sobre o caminho |
| `team_size` | Obrigatório no Truco: `2` (duplas) ou `3` (trios). Não permitido na Fodinha |
| `limit` | Inteiro positivo, padrão `10`; valores acima de `100` são limitados a `100` |

O serviço remove espaços nas extremidades dos filtros e converte `mode` para minúsculas. Não há offset ou paginação por cursor. Jogadores sem partidas finalizadas no recorte não aparecem.

### Fodinha

```bash
curl 'http://localhost:8080/api/v1/ranking/fodinha?limit=10'
```

```json
{
  "data": [{
    "position": 1,
    "player_id": "11111111-1111-4111-8111-111111111111",
    "name": "Jogador Exemplo",
    "avatar_key": "exemplo",
    "games_played": 5,
    "gold_medals": 3,
    "silver_medals": 1,
    "bronze_medals": 1,
    "podiums": 5
  }]
}
```

Ouro, prata e bronze contam colocações 1, 2 e 3; `podiums` conta colocações de 1 a 3. Ordenação: ouro, prata, bronze, pódios e partidas jogadas, todos decrescentes; depois `avatar_key` crescente. Campos de vitórias/derrotas/aproveitamento são omitidos neste modo.

### Truco individual

```bash
curl 'http://localhost:8080/api/v1/ranking/truco?team_size=3&limit=10'
```

```json
{
  "data": [{
    "position": 1,
    "player_id": "11111111-1111-4111-8111-111111111111",
    "name": "Jogador Exemplo",
    "avatar_key": "exemplo",
    "games_played": 10,
    "wins": 7,
    "losses": 3,
    "win_rate": 70
  }]
}
```

O recorte exige partidas com exatamente duas equipes e quantidade de integrantes igual ao `team_size` em cada uma, agrupadas por `team_name`. `win_rate` é percentual de 0 a 100, arredondado a duas casas. Ordenação: vitórias, aproveitamento e partidas jogadas, decrescentes; depois `avatar_key` crescente. Campos de medalhas são omitidos.

## Ranking de equipes de Truco

`GET /api/v1/ranking/truco/teams?team_size=2&limit=10` — `team_size` é obrigatório (2 ou 3); `limit` segue as regras anteriores. O modo é fixo em Truco nesta rota.

```json
{
  "data": [{
    "position": 1,
    "team_key": "11111111-1111-4111-8111-111111111111:22222222-2222-4222-8222-222222222222",
    "players": [
      {"id":"11111111-1111-4111-8111-111111111111","name":"Jogador A","avatar_key":"a"},
      {"id":"22222222-2222-4222-8222-222222222222","name":"Jogador B","avatar_key":"b"}
    ],
    "games_played": 10,
    "wins": 7,
    "losses": 3,
    "win_rate": 70,
    "points_for": 132,
    "points_against": 94,
    "point_diff": 38,
    "last_played_at": "2026-11-13T12:00:00Z"
  }]
}
```

`team_key` identifica a composição da equipe por UUIDs de jogadores ordenados e unidos por `:`. A mesma composição é acumulada entre partidas mesmo que o nome da equipe mude. Não é o `team_id` de uma partida. A lista `players` segue a ordem dos IDs.

`points_for` soma a pontuação da equipe, `points_against` soma a do adversário e `point_diff` é a diferença. `last_played_at` é a data mais recente de finalização, usando criação como fallback quando necessário.

Ordenação: vitórias, aproveitamento, saldo de pontos, pontos marcados e partidas jogadas, decrescentes; depois `team_key` crescente. `position` usa numeração sequencial, sem compartilhar posições em empates.

## Erros

Validações retornam `400` com `error`:

- `modo do ranking é obrigatório` ou `modo do ranking inválido`.
- `team_size é obrigatório para ranking de truco`.
- `team_size não é permitido para ranking de fodinha`.
- `team_size deve ser 2 para dupla ou 3 para trio`.
- `limit deve ser maior que zero` (também para valores não inteiros).

Falhas internas retornam `500` com `erro ao buscar ranking` ou `erro ao buscar ranking de duplas` — esta última mensagem também é usada para trios.
