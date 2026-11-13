# Salas em tempo real via WebSocket

[README do projeto](../../README.md#documentação) · [Partidas](../api/PARTIDAS.md)

## Conexão

Conecte em `ws://localhost:8080/api/v1/ws/rooms/:room_code`. Em uma implantação com HTTPS, use `wss://` no mesmo caminho.

A sala é consultada antes do upgrade. Código inexistente retorna HTTP `404`; sala cancelada retorna HTTP `409` com `{"message":"Essa partida foi cancelada pelo admin."}`. O código precisa corresponder exatamente ao retornado na criação. Salas finalizadas permitem conexão. Não é exigido JWT nem host token para assistir.

A aceitação de conexões usa `InsecureSkipVerify: true`, desabilitando a verificação de origem WebSocket; `ALLOWED_ORIGINS` configura o CORS HTTP, mas não restringe esse upgrade.

## Envelope

Mensagens são JSON de texto, sem o envelope HTTP `data`:

```json
{
  "type": "room_state_updated",
  "room_code": "ABC234",
  "payload": {"round": 2}
}
```

`type` define o evento. `payload` é JSON, não uma string contendo JSON. `room_code`, `payload`, `host_token` e `error` são omitidos quando vazios. A conexão determina a sala; enviar outro `room_code` no corpo não muda seu destino.

## Eventos do servidor

| `type` | Quando | Conteúdo |
| --- | --- | --- |
| `room_snapshot` | Após conectar | `payload` é o modelo completo da partida, incluindo o último `live_state`, se existir |
| `viewer_count_updated` | Entrada/saída de conexão | `payload: {"viewer_count": N}` |
| `room_state_updated` | Atualização do host persistida | `payload` contém o JSON recebido do host |
| `pong` | Resposta a `ping` | Sem payload |
| `error` | Mensagem rejeitada | Campo `error` com a descrição; sem envelope HTTP |

Snapshot inicial:

```json
{
  "type":"room_snapshot",
  "room_code":"ABC234",
  "payload":{
    "id":"aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
    "mode":"fodinha",
    "status":"ongoing",
    "room_code":"ABC234",
    "lives":3,
    "created_at":"2026-11-13T12:00:00Z"
  }
}
```

O snapshot não inclui `host_token`. O contador mede conexões, incluindo a do host e múltiplas abas da mesma pessoa. Seu valor e os broadcasts são locais ao processo da API.

## Mensagens do cliente

### Ping da aplicação

```json
{"type":"ping"}
```

Resposta: `{"type":"pong","room_code":"ABC234"}`. Este é um evento JSON da aplicação, diferente dos frames de controle do protocolo WebSocket. Não há heartbeat periódico implementado pelo servidor.

### Publicar estado

Use o `host_token` retornado por `POST /matches`:

```json
{
  "type":"room_state_updated",
  "host_token":"TOKEN_RETORNADO_NA_CRIACAO",
  "payload":{
    "round":2,
    "current_player_id":"11111111-1111-4111-8111-111111111111"
  }
}
```

Os campos do exemplo em `payload` são ilustrativos: o backend aceita qualquer valor JSON válido, inclusive `null`, sem validar regras das rodadas ou estrutura interna. O payload é obrigatório no envelope. O servidor valida o token, substitui integralmente `matches.live_state` e transmite o mesmo payload para todos os clientes da sala, incluindo o remetente. Não inclui o host token no broadcast.

Atualizações não são patches nem são mescladas ao estado anterior. Não finalizam a partida, não alteram o ranking e não exigem status `ongoing`. A validação/persistência tem timeout de cinco segundos; cada envio também tem timeout de cinco segundos.

### Erros

```json
{"type":"error","room_code":"ABC234","error":"host_token inválido"}
```

Mensagens possíveis: `type é obrigatório`, `tipo de mensagem não suportado`, `estado da sala inválido`, `host_token inválido`, `partida não encontrada` e `erro ao atualizar estado da sala`. Uma rejeição de mensagem não encerra automaticamente a conexão; JSON malformado ou falha de leitura encerra o loop do cliente.

## Exemplo no navegador

Execute no contexto do frontend, substituindo código e token:

```javascript
const roomCode = "ABC234";
const socket = new WebSocket(
  `ws://localhost:8080/api/v1/ws/rooms/${encodeURIComponent(roomCode)}`
);

socket.addEventListener("message", ({ data }) => {
  const event = JSON.parse(data);
  switch (event.type) {
    case "room_snapshot":
      console.log("Partida inicial:", event.payload);
      break;
    case "room_state_updated":
      console.log("Novo estado:", event.payload);
      break;
    case "viewer_count_updated":
      console.log("Conexões:", event.payload.viewer_count);
      break;
    case "error":
      console.error(event.error);
      break;
  }
});

socket.addEventListener("open", () => {
  socket.send(JSON.stringify({ type: "ping" }));
});

function publishState(hostToken, state) {
  if (socket.readyState !== WebSocket.OPEN) return;
  socket.send(JSON.stringify({
    type: "room_state_updated",
    host_token: hostToken,
    payload: state
  }));
}
```

Um espectador apenas lê eventos; o host chama `publishState` quando precisar publicar. O cliente precisa implementar reconexão, se desejado; uma nova conexão recebe outro snapshot. Mutações HTTP como inclusão de participantes, finalização e cancelamento não geram broadcasts automaticamente. Consulte a API HTTP para atualizar esses dados.
