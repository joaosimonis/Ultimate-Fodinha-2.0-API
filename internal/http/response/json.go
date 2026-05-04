package response

import "github.com/gin-gonic/gin"

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func Message(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"message": message})
}

func Data(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}
