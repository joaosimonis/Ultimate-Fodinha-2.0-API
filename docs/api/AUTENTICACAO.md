# Autenticação administrativa

[README do projeto](../../README.md#documentação) · [Rotas](README.md)

## Login

`POST /api/v1/auth/login` — público.

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"sua-senha"}'
```

`username` e `password` são obrigatórios. O serviço remove espaços nas extremidades do username e a consulta ao banco ignora maiúsculas/minúsculas. A senha é comparada ao hash bcrypt. Apenas usuários com papel `admin` podem autenticar.

Resposta `200`:

```json
{
  "data": {
    "authenticated": true,
    "access_token": "JWT_RETORNADO_PELA_API",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {"username": "admin", "role": "admin"}
  }
}
```

`expires_in` acompanha `JWT_EXPIRES_IN_SECONDS`. O token gerado usa HS256 e contém `sub` (UUID do usuário), `username`, `role`, `iat` e `exp`.

## Usar o token

Envie `Authorization: Bearer <access_token>` nas rotas administrativas: criar jogador, excluir jogador e cancelar partida por código de sala. Por exemplo, com `TOKEN` definido no seu shell:

```bash
curl -X PATCH http://localhost:8080/api/v1/admin/matches/room/ABC234/cancel \
  -H "Authorization: Bearer $TOKEN"
```

Não há endpoint de cadastro de usuário, renovação de token ou logout. Para obter outro token, autentique novamente. O [comando administrativo](../projeto/CONFIGURACAO.md) cria/atualiza usuários; alterar a senha não revoga automaticamente tokens já emitidos.

## Erros

| Status | Situação | Corpo |
| --- | --- | --- |
| 400 | JSON inválido ou campo obrigatório ausente | `{"error":"dados inválidos"}` |
| 401 | Usuário/senha inválidos ou papel diferente de admin | `{"error":"credenciais inválidas"}` |
| 401 | Bearer ausente ou token inválido | `{"error":"token inválido"}` |
| 401 | JWT expirado | `{"error":"token expirado"}` |
| 500 | Falha interna no login | `{"error":"erro ao autenticar"}` |

O `host_token` de partidas não é um JWT administrativo. Ele é retornado apenas na criação da partida e autoriza as atualizações do [estado via WebSocket](../websocket/WEBSOCKET.md).
