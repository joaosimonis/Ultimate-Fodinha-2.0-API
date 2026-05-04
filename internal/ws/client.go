package ws

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/matches"
)

const (
	MessageTypeRoomSnapshot     = "room_snapshot"
	MessageTypeRoomStateUpdated = "room_state_updated"
	MessageTypeViewerCount      = "viewer_count_updated"
	MessageTypePing             = "ping"
	MessageTypePong             = "pong"
	MessageTypeError            = "error"

	validationTimeout = 5 * time.Second
	sendTimeout       = 5 * time.Second
)

type MatchService interface {
	UpdateLiveState(ctx context.Context, roomCode string, hostToken string, liveState json.RawMessage) error
}

type Message struct {
	Type      string          `json:"type"`
	RoomCode  string          `json:"room_code,omitempty"`
	HostToken string          `json:"host_token,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Error     string          `json:"error,omitempty"`
}

type Client struct {
	hub      *Hub
	service  MatchService
	conn     *websocket.Conn
	roomCode string
	match    matches.Match
	sendMu   sync.Mutex
}

func NewClient(hub *Hub, service MatchService, conn *websocket.Conn, roomCode string, match matches.Match) *Client {
	return &Client{
		hub:      hub,
		service:  service,
		conn:     conn,
		roomCode: roomCode,
		match:    match,
	}
}

func (c *Client) Run(ctx context.Context) {
	viewerCount := c.hub.Register(c)
	defer func() {
		viewerCount := c.hub.Unregister(c)
		c.hub.Broadcast(c.roomCode, viewerCountMessage(c.roomCode, viewerCount))
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	c.Send(snapshotMessage(c.roomCode, c.match))
	c.hub.Broadcast(c.roomCode, viewerCountMessage(c.roomCode, viewerCount))

	for {
		var message Message
		if err := wsjson.Read(ctx, c.conn, &message); err != nil {
			return
		}

		c.handleMessage(message)
	}
}

func (c *Client) Send(message Message) {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), sendTimeout)
	defer cancel()

	if err := wsjson.Write(ctx, c.conn, message); err != nil {
		c.hub.Unregister(c)
	}
}

func (c *Client) handleMessage(message Message) {
	message.Type = strings.TrimSpace(message.Type)
	switch message.Type {
	case MessageTypePing:
		c.Send(Message{Type: MessageTypePong, RoomCode: c.roomCode})
	case MessageTypeRoomStateUpdated:
		c.handleRoomStateUpdated(message)
	case "":
		c.Send(errorMessage(c.roomCode, "type é obrigatório"))
	default:
		c.Send(errorMessage(c.roomCode, "tipo de mensagem não suportado"))
	}
}

func (c *Client) handleRoomStateUpdated(message Message) {
	ctx, cancel := context.WithTimeout(context.Background(), validationTimeout)
	defer cancel()

	if err := c.service.UpdateLiveState(ctx, c.roomCode, message.HostToken, message.Payload); err != nil {
		c.Send(errorMessage(c.roomCode, websocketErrorMessage(err)))
		return
	}

	c.hub.Broadcast(c.roomCode, Message{
		Type:     MessageTypeRoomStateUpdated,
		RoomCode: c.roomCode,
		Payload:  message.Payload,
	})
}

func snapshotMessage(roomCode string, match matches.Match) Message {
	payload, err := json.Marshal(match)
	if err != nil {
		return errorMessage(roomCode, "erro ao montar snapshot da sala")
	}

	return Message{
		Type:     MessageTypeRoomSnapshot,
		RoomCode: roomCode,
		Payload:  payload,
	}
}

func viewerCountMessage(roomCode string, viewerCount int) Message {
	payload, err := json.Marshal(map[string]int{"viewer_count": viewerCount})
	if err != nil {
		return errorMessage(roomCode, "erro ao montar contador da sala")
	}

	return Message{
		Type:     MessageTypeViewerCount,
		RoomCode: roomCode,
		Payload:  payload,
	}
}

func errorMessage(roomCode string, message string) Message {
	return Message{
		Type:     MessageTypeError,
		RoomCode: roomCode,
		Error:    message,
	}
}

func websocketErrorMessage(err error) string {
	switch {
	case errors.Is(err, matches.ErrInvalidHostToken):
		return matches.ErrInvalidHostToken.Error()
	case errors.Is(err, matches.ErrInvalidLiveState):
		return matches.ErrInvalidLiveState.Error()
	case errors.Is(err, matches.ErrNotFound):
		return matches.ErrNotFound.Error()
	default:
		return "erro ao atualizar estado da sala"
	}
}
