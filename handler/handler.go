package handler

import "github.com/gin-gonic/gin"

type PageAnalyserHandler interface {
	Analyse(c *gin.Context)
}

type pageAnalyserHandler struct{}

func NewPageAnalyserHandler() PageAnalyserHandler {
	return &pageAnalyserHandler{}
}
