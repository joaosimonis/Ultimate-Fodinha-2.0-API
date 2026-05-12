package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/auth"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/response"
)

func AdminAuth(tokens *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := auth.BearerToken(c.GetHeader("Authorization"))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, auth.ErrInvalidToken.Error())
			c.Abort()
			return
		}

		claims, err := tokens.Validate(token)
		if err != nil {
			status := http.StatusUnauthorized
			if errors.Is(err, auth.ErrExpiredToken) {
				status = http.StatusUnauthorized
			}

			response.Error(c, status, auth.TokenErrorMessage(err))
			c.Abort()
			return
		}

		c.Set("admin_user_id", claims.Subject)
		c.Set("admin_username", claims.Username)
		c.Next()
	}
}
