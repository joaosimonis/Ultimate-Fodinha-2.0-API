package players

import "time"

type Ultimate struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type Player struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    AvatarKey string    `json:"avatar_key"`
    Title     string    `json:"title"`
    Lore      string    `json:"lore"`
    Ultimate  *Ultimate `json:"ultimate"`
    CreatedAt time.Time `json:"created_at"`
}

type CreatePlayerRequest struct {
    Name      string    `json:"name"      binding:"required"`
    AvatarKey string    `json:"avatar_key" binding:"required"`
    Title     string    `json:"title"`
    Lore      string    `json:"lore"`
    Ultimate  *Ultimate `json:"ultimate"  binding:"required"`
}