package ws

import "sync"

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]map[*Client]struct{}),
	}
}

func (h *Hub) Register(client *Client) int {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[client.roomCode] == nil {
		h.rooms[client.roomCode] = make(map[*Client]struct{})
	}

	h.rooms[client.roomCode][client] = struct{}{}
	return len(h.rooms[client.roomCode])
}

func (h *Hub) Unregister(client *Client) int {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients := h.rooms[client.roomCode]
	if clients == nil {
		return 0
	}

	delete(clients, client)
	if len(clients) == 0 {
		delete(h.rooms, client.roomCode)
		return 0
	}

	return len(clients)
}

func (h *Hub) Broadcast(roomCode string, message Message) {
	clients := h.clients(roomCode)
	for _, client := range clients {
		client.Send(message)
	}
}

func (h *Hub) clients(roomCode string) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	roomClients := h.rooms[roomCode]
	if roomClients == nil {
		return nil
	}

	clients := make([]*Client, 0, len(roomClients))
	for client := range roomClients {
		clients = append(clients, client)
	}

	return clients
}
