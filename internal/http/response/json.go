package response

import "github.com/gin-gonic/gin"

/*
	Pacote para padronizar respostas JSON
	Reduz repetição e padroniza payloads
*/

// 1. Função para retornar respostas de erro
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// 2. Função para retornar respostas de mensagem
func Message(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"message": message})
}

// 3. Função para retornar respostas de dados
func Data(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}
