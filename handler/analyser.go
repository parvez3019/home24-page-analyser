package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *pageAnalyserHandler) Analyse(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "message"})
}
