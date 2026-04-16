package handlers

import (
	"context"
	"errors"   // Pacote padrão do GO -> errors.Is compara erros com segurança
	"net/http" // Pacote padrão para criar APIs REST

	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/response" // Pacote para criar respostas padronizadas (JSON)
	"github.com/joaosimonis/ultimate-fodinha-api/internal/players"       // Pacote de lógica de negócio para jogadores

	"github.com/gin-gonic/gin" // Framework para criar APIs REST
)

type PlayerService interface {
	ListPlayers(ctx context.Context) ([]players.Player, error)
	GetPlayerByID(ctx context.Context, id string) (players.Player, error)
	CreatePlayer(ctx context.Context, req players.CreatePlayerRequest) (players.Player, error)
	DeletePlayer(ctx context.Context, id string) error
}

// Struct que vai guardar a dependência do handler. Essa dependência é o service de jogadores
// Isso é bom pro handler não conhecer SQL nem detalhes de persistência.
type PlayerHandler struct {
	service PlayerService
}

/*
	É o meu construtor
	Recebe o service e devolve o handler pronto
	Agora o handler depende de um contrato, não do repository concreto
*/
func NewPlayerHandler(service PlayerService) *PlayerHandler {
	return &PlayerHandler{
		service: service,
	}
}

/*
	A partir daqui estão meus handlers
	Podemos ver um padrão: recebe um receiver (instância do handler) e também o contexto da requisição no Gin
*/

// 1. Handler que irá listar todos os jogadores
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	// ListPlayers é um método do service, que encapsula a regra da listagem
	// Armazenamos essa busca na variável playerList
	playersList, err := h.service.ListPlayers(c.Request.Context())
	// Se houver erro, retornamos um erro 500 (internet server error)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "erro ao buscar jogadores")
		return
	}

	// Se tudo deu certo, retornamos com status 200 e a lista de jogadores
	response.Data(c, http.StatusOK, playersList)
}

// 2. Handler que irá buscar um jogador pelo ID
func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	// Param retorna o valor do parâmetro da rota, nesse caso estamos extraindo o "id" e atribuindo à uma variável id
	id := c.Param("id")

	// GetPlayerByID é um método do service, que centraliza a busca do jogador
	player, err := h.service.GetPlayerByID(c.Request.Context(), id)
	if err != nil {
		// Verifica se o erro é de jogador não encontrado. Esse ErrNotFound vem do arquivo de errors que criamos!
		if errors.Is(err, players.ErrNotFound) {
			// Se for, retornamos um erro 404 (not found)
			response.Error(c, http.StatusNotFound, players.ErrNotFound.Error())
			return
		}

		// Se não for o erro esperado, retornamos um erro 500 (internal server error)
		response.Error(c, http.StatusInternalServerError, "erro ao buscar jogador")
		return
	}

	// Se tudo deu certo, retornamos com status 200 e o jogador encontrado
	response.Data(c, http.StatusOK, player)
}

// 3. Handler que irá criar um novo jogador
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	// Criamos uma variável req do tipo players.CreatePlayerRequest (meu model que representa a requisição de criação de jogador)
	// Poderíamos fazer assim também: req := players.CreatePlayerRequest
	var req players.CreatePlayerRequest
	// Usamos c.ShouldBindJSON para ler o body do JSON e preencher a Struct
	if err := c.ShouldBindJSON(&req); err != nil {
		// Se houver erro, retornamos um erro 400 (bad request) e a mensagem de dados inválidos
		response.Error(c, http.StatusBadRequest, "dados inválidos")
		return
	}

	// O service centraliza a validação de nome duplicado e a criação do jogador
	player, err := h.service.CreatePlayer(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, players.ErrNameTaken) {
			response.Error(c, http.StatusConflict, players.ErrNameTaken.Error())
			return
		}

		// Se houver erro, retornamos um erro 500 (internal server error)
		response.Error(c, http.StatusInternalServerError, "erro ao criar jogador")
		return
	}

	// Se tudo deu certo, retornamos com status 201 (created) e o jogador criado
	response.Data(c, http.StatusCreated, player)
}

// 4. Handler que irá remover um jogador
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	// Achamos o "id" na URL da rota
	id := c.Param("id")

	// Chamamos o método DeletePlayer do service para remover o jogador
	if err := h.service.DeletePlayer(c.Request.Context(), id); err != nil {
		// Verifica se o erro é de jogador não encontrado. Esse ErrNotFound vem do arquivo de errors que criamos!
		if errors.Is(err, players.ErrNotFound) {
			response.Error(c, http.StatusNotFound, players.ErrNotFound.Error())
			return
		}

		// Se houver outro erro, retornamos um erro 500 (internal server error)
		response.Error(c, http.StatusInternalServerError, "erro ao remover jogador")
		return
	}

	// Se tudo deu certo, retornamos com status 200 (ok) e uma mensagem de sucesso
	response.Message(c, http.StatusOK, "jogador removido com sucesso")
}
