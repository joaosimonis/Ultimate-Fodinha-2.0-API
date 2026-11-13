# Arquitetura e persistência

[README do projeto](../../README.md#documentação)

## Organização

| Caminho | Responsabilidade |
| --- | --- |
| `cmd/api` | Inicialização do servidor |
| `cmd/admin` | Criação/atualização de administrador com bcrypt |
| `cmd/seed` | Importação do catálogo inicial de jogadores |
| `internal/bootstrap` | Composição da configuração, conexão, middleware e rotas |
| `internal/config` | Leitura e validação do ambiente |
| `internal/db` | Pool PostgreSQL e migrations SQL |
| `internal/http/routes` | Registro de endpoints e proteção administrativa |
| `internal/http/handlers` | Leitura HTTP, chamada aos serviços e mapeamento de erros |
| `internal/http/middleware` | JWT administrativo e CORS |
| `internal/http/response` | Envelopes JSON |
| `internal/auth` | Login, consulta de usuários, geração e validação de JWT |
| `internal/players` | Cadastro e consulta de jogadores |
| `internal/matches` | Regras de partidas, equipes, resultados e estado ao vivo |
| `internal/ranking` | Validação de filtros e agregações dos rankings |
| `internal/ws` | Clientes WebSocket e distribuição de mensagens por sala |

O fluxo HTTP é rota → middleware, quando aplicável → handler → serviço → repositório → PostgreSQL. Handlers e serviços usam interfaces para suas dependências. O bootstrap conecta implementações concretas com Gin e pgx. O WebSocket usa o serviço de partidas para validar o host e persistir o estado antes de distribuí-lo.

## Modelo de dados

| Tabela | Uso e relacionamentos |
| --- | --- |
| `players` | Perfil do jogador e `ultimate` em JSONB; nome único |
| `matches` | Modo, status, código único da sala, hash do host token, configuração, JSONB de estado e datas |
| `match_players` | Participação única por `(match_id, player_id)`, equipe e resultado individual |
| `match_teams` | Equipes únicas por `(match_id, name)`, com ordem de exibição |
| `users` | Administradores, hash bcrypt e índice único por `LOWER(username)` |
| `player_stats` | Estrutura de estatísticas por jogador/modo; não é usada pelas consultas atuais de ranking |

Excluir uma partida no banco remove suas equipes e participações por cascata. Excluir uma equipe deixa `team_id` dos participantes nulo. A referência de `match_players` a `players` impede excluir um jogador que tenha participações. A API não oferece exclusão de partidas ou equipes.

Os resultados de equipes são calculados a partir de seus participantes: maior `final_score`, menor `placement` e existência de vencedor. A finalização grava resultados dos participantes e status/data da partida em uma transação. A inclusão de participantes e equipes também usa transação. Rankings consultam as partidas finalizadas e seus participantes diretamente.

## Estado ao vivo

`matches.live_state` persiste o último JSON publicado pelo host. Seu conteúdo é definido pelo cliente; o servidor não executa a lógica das rodadas nem deriva o resultado final desse JSON. Finalizar exige uma chamada HTTP explícita.

O hub de conexões e a contagem de espectadores vivem na memória de cada processo. Não há barramento compartilhado entre réplicas: conexões atendidas por processos diferentes não recebem automaticamente os mesmos broadcasts. Reiniciar a API encerra conexões, mas o último estado persistido permanece no banco.

## Pontos de manutenção

O código atual tem testes de validação da finalização de Truco em `internal/matches/match_service_test.go`. Não há suíte HTTP/banco completa no repositório. Ao evoluir o contrato, confira a [referência de rotas](../api/README.md) e os comportamentos específicos documentados em cada guia.
