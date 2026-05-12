package auth

type User struct {
	ID           string
	Username     string
	PasswordHash string
	Role         string
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Authenticated bool      `json:"authenticated"`
	AccessToken   string    `json:"access_token"`
	TokenType     string    `json:"token_type"`
	ExpiresIn     int       `json:"expires_in"`
	User          AdminUser `json:"user"`
}

type AdminUser struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
