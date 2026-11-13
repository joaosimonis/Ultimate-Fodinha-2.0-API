# Jogadores

[README do projeto](../../README.md#documentação) · [Rotas](README.md)

## Modelo

```json
{
  "id": "11111111-1111-4111-8111-111111111111",
  "name": "Jogador Exemplo",
  "avatar_key": "exemplo",
  "title": "Guardião",
  "lore": "História do personagem.",
  "ultimate": {"name": "Proteção", "description": "Protege os aliados."},
  "created_at": "2026-11-13T12:00:00Z"
}
```

`ultimate` pode ser `null` em registros consultados; no cadastro é obrigatório enviar um objeto. `title` e `lore` são opcionais no cadastro. `avatar_key` é uma chave de identificação do avatar, não uma URL; a API não fornece upload de imagens.

## Listar e consultar

`GET /api/v1/players` — público. Retorna apenas jogadores com `is_active = TRUE`, ordenados por `avatar_key` crescente, com `200` e `{"data":[<jogador>, ...]}`. Sem jogadores, retorna `{"data":[]}`. Não há busca ou paginação.

Para ocultar um jogador do catálogo e da seleção sem apagar seu histórico, defina `players.is_active = FALSE` no banco. Para reexibir, use `TRUE`. Aplique a migration `008` antes de publicar esta versão da API. A consulta por ID, as partidas e os rankings continuam incluindo os jogadores ocultos. O filtro controla a seleção no frontend; não impede adicionar um ID conhecido diretamente pela API.

`GET /api/v1/players/:id` — público. Retorna `200` com `{"data":<jogador>}` ou `404` com `{"error":"jogador não encontrado"}`.

```bash
curl http://localhost:8080/api/v1/players
curl http://localhost:8080/api/v1/players/11111111-1111-4111-8111-111111111111
```

## Criar

`POST /api/v1/players` — exige [JWT administrativo](AUTENTICACAO.md).

```bash
curl -X POST http://localhost:8080/api/v1/players \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "name":"Jogador Exemplo",
    "avatar_key":"exemplo",
    "title":"Guardião",
    "lore":"História do personagem.",
    "ultimate":{"name":"Proteção","description":"Protege os aliados."}
  }'
```

Obrigatórios: `name`, `avatar_key` e `ultimate`. Os campos internos de `ultimate` não têm validação de obrigatoriedade. O serviço verifica duplicidade de nome ignorando maiúsculas/minúsculas, sem remover espaços.

Sucesso: `201` com `{"data":<jogador criado>}`. JSON/campos obrigatórios inválidos: `400`, `dados inválidos`. Nome em uso: `409`, `já existe um jogador com esse nome`. Falhas internas: `500`, `erro ao criar jogador`.

## Excluir

`DELETE /api/v1/players/:id` — exige JWT administrativo; sem corpo.

```bash
curl -X DELETE http://localhost:8080/api/v1/players/11111111-1111-4111-8111-111111111111 \
  -H "Authorization: Bearer $TOKEN"
```

Sucesso: `200`, `{"message":"jogador removido com sucesso"}`. Jogador inexistente: `404`. Jogadores referenciados em `match_players` não podem ser excluídos pela restrição do banco; atualmente essa falha retorna `500`, `erro ao remover jogador`.

Não há endpoint de edição de jogador. Nas consultas, falhas internas retornam `500` com `erro ao buscar jogadores` ou `erro ao buscar jogador`.
