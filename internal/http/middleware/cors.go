package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigins string) gin.HandlerFunc {
	origins := parseAllowedOrigins(allowedOrigins)

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" && origins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func parseAllowedOrigins(rawOrigins string) map[string]bool {
	origins := make(map[string]bool)

	for _, origin := range strings.Split(rawOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins[origin] = true
		}
	}

	return origins
}
