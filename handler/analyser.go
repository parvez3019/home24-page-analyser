package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) Analyse(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "message"})
}
