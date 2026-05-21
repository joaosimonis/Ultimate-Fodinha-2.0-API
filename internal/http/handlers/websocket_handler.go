package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/response"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/matches"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/ws"
)

type RoomWebSocketService interface {
	GetByRoomCode(ctx context.Context, roomCode string) (matches.Match, error)
	UpdateLiveState(ctx context.Context, roomCode string, hostToken string, liveState json.RawMessage) error
}

type WebSocketHandler struct {
	hub     *ws.Hub
	service RoomWebSocketService
}

func NewWebSocketHandler(hub *ws.Hub, service RoomWebSocketService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:     hub,
		service: service,
	}
}

func (h *WebSocketHandler) JoinRoom(c *gin.Context) {
	roomCode := c.Param("room_code")

	match, err := h.service.GetByRoomCode(c.Request.Context(), roomCode)
	if err != nil {
		if errors.Is(err, matches.ErrNotFound) {
			response.Error(c, http.StatusNotFound, matches.ErrNotFound.Error())
			return
		}

		if errors.Is(err, matches.ErrMatchCancelled) {
			response.Message(c, http.StatusConflict, matches.ErrMatchCancelled.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "erro ao buscar sala")
		return
	}

	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}

	client := ws.NewClient(h.hub, h.service, conn, roomCode, match)
	client.Run(context.Background())
}
