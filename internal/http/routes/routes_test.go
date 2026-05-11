package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupRoutesRegistersRankingTeamRoutes(t *testing.T) {
	t.Setenv("GIN_MODE", gin.TestMode)

	router := gin.New()

	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("SetupRoutes panicked: %v", recovered)
		}
	}()

	SetupRoutes(router, nil)
}
